package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/swavan.io/gateway/internal/config"
	"github.com/swavan.io/gateway/pkg/authentication"
)

func Run(ctx context.Context, mux *http.ServeMux, auth authentication.AuthenticationAPI) error {
	mux.HandleFunc("/health", health)

	authMiddleware, err := NewAuthMiddleware(ctx, auth)
	if err != nil {
		return err
	}
	for _, resource := range config.Config.Resources {

		if !resource.Active {
			continue
		}

		url, err := url.Parse(resource.Destination)
		if err != nil {
			return err
		}

		if resource.Authenticated {
			mux.HandleFunc(
				resource.Endpoint,
				authMiddleware.Guard(NewReverseProxy(url, resource.Endpoint).ServeHTTP),
			)
			continue
		}

		fmt.Println(resource.Endpoint)
		mux.HandleFunc(
			resource.Endpoint,
			NewReverseProxy(url, resource.Endpoint).ServeHTTP,
		)

	}

	NewLogger(mux)
	return nil
}
