package handler

import (
	"context"
	"net/http"
	"net/url"

	"github.com/swavan.io/gateway/internal/config"
)

func Run(ctx context.Context, mux *http.ServeMux) error {
	mux.HandleFunc("/health", health)
	for _, resource := range config.Config.Services {
		url, err := url.Parse(resource.Destination)
		if err != nil {
			return err
		}
		mux.HandleFunc(
			resource.Endpoint,
			NewReverseProxy(url, resource.Endpoint).Redirect,
		)
	}
	return nil
}
