package handler

import (
	"context"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/swavan.io/gateway/pkg/authentication"
	"github.com/swavan.io/gateway/pkg/authentication/key"
	"github.com/swavan.io/gateway/pkg/identity"
)

type Logger struct {
	handler http.Handler
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	l.handler.ServeHTTP(w, r)
	log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
}

func NewLogger(handlerToWrap http.Handler) *Logger {
	return &Logger{handlerToWrap}
}

type Auth struct {
	api authentication.AuthenticationAPI
	key *key.Key
}

func NewAuthMiddleware(ctx context.Context, api authentication.AuthenticationAPI) (*Auth, error) {
	key, err := api.
		Key().
		FetchKey(ctx, os.Getenv("APP_NAME"))
	if err != nil {
		return nil, err
	}
	return &Auth{api, key}, nil
}

func (a *Auth) Guard(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get("Authorization")
		if accessToken != "" {
			accessToken = strings.TrimPrefix(accessToken, "Bearer ")
		} else {
			cookie, err := r.Cookie(string(identity.AccessToken))
			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			accessToken = cookie.Value
		}

		if accessToken == "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		claims, _, err := authentication.ParseAsymmetricToken(
			accessToken,
			a.key.PublicKey,
		)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		w.Header().Add("X-AUTH-USER", claims.Username)

		h.ServeHTTP(
			w,
			r.WithContext(
				context.WithValue(
					r.Context(),
					identity.AuthenticatedUser, claims)))

	})
}

func (a *Auth) Access(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		path := r.URL.Path
		if slices.Contains(a.api.Config().IgnoreAccess, path) {
			h.ServeHTTP(w, r)
			return
		}

		claims := r.Context().
			Value(identity.AuthenticatedUser).(*authentication.Claims)

		action := r.Method
		s, err := a.api.Access().Enforcer().Enforce(
			claims.Username,
			claims.Domain.ID,
			path,
			action,
		)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !s {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	})
}
