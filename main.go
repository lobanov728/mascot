package main

import (
	"context"
	"net/http"

	"github.com/Lobanov728/mascot/config"
	"github.com/Lobanov728/mascot/internal/billing/ports/jsonrpcport"
	"github.com/Lobanov728/mascot/internal/billing/service"
)

func main() {
	ctx := context.Background()

	cfg := config.Config{}
	config.Init("./config.yaml", &cfg)
	billingApp := service.NewApplication(ctx, cfg)

	handler := jsonrpcport.NewJSONRpcServer(billingApp)
	http.Handle("/", handler)

	srv := http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	err := srv.ListenAndServe()

	if err != nil {
		panic("Unable to start HTTP server")
	}
}
