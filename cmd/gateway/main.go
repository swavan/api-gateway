package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/swavan.io/gateway/config"
	srvConfig "github.com/swavan.io/gateway/internal/config"
	"github.com/swavan.io/gateway/internal/db"
	"github.com/swavan.io/gateway/internal/handler"
	"github.com/swavan.io/gateway/pkg/authentication"
)

func main() {
	err := config.New(config.Configuration(), &srvConfig.Config)
	if err != nil {
		panic(fmt.Errorf("could not load configuration: %v", err))
	}

	connection, err := db.Start(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	authConfig := authentication.NewConfig()
	err = config.New(
		config.Configuration().
			SetFileExtension(authConfig.Extension()).
			SetFileName(authConfig.Name()).
			SetFilePath(authConfig.Path()),
		&authConfig)

	if err != nil {
		panic(fmt.Errorf("could not load configuration: %v", err))
	}
	authentication, err := authentication.New(connection.GetDB(), authConfig)
	if err != nil {
		panic(err)
	}
	mux := http.NewServeMux()
	handler.Run(context.Background(), mux, authentication)
	http.ListenAndServe(":8000", mux)
}
