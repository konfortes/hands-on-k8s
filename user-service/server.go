package main

import (
	"log"
	"net/http"

	"github.com/konfortes/go-server-utils/server"
	"github.com/konfortes/go-server-utils/utils"
	opentracing "github.com/opentracing/opentracing-go"
)

var (
	tracer *opentracing.Tracer
)

const (
	serviceName = "hands-on-k8s-user-service"
)

func main() {
	serverConfig := server.Config{
		AppName:     "my-app-name",
		Port:        utils.GetEnvOr("PORT", "4432"),
		Env:         utils.GetEnvOr("ENV", "development"),
		Handlers:    handlers(),
		WithTracing: utils.GetEnvOr("TRACING_ENABLED", "false") == "true",
	}

	srv := server.Initialize(serverConfig)

	go func() {
		log.Println("listening on " + srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	server.GracefulShutdown(srv)
}

func handlers() []server.Handler {
	return []server.Handler{
		{
			Method:  http.MethodPost,
			Pattern: "/users",
			H:       usersHandler,
		},
	}
}
