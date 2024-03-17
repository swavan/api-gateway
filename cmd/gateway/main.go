package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/swavan.io/gateway/config"
	srvConfig "github.com/swavan.io/gateway/internal/config"
	"github.com/swavan.io/gateway/internal/handler"
)

func main() {
	err := config.New(config.Configuration(), &srvConfig.Config)
	if err != nil {
		panic(fmt.Errorf("could not load configuration: %v", err))
	}
	mux := http.NewServeMux()
	handler.Run(context.Background(), mux)
	http.ListenAndServe(":8000", mux)
}
