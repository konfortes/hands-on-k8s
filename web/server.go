package main

import (
	"log"
	"net/http"

	"github.com/konfortes/go-server-utils/server"
	"github.com/konfortes/go-server-utils/utils"
	opentracing "github.com/opentracing/opentracing-go"
)

var (
	tracer      *opentracing.Tracer
	userService UserService
)

const (
	appName = "hands-on-k8s-web"
)

func main() {
	initialize()

	serverConfig := server.Config{
		AppName:     appName,
		Port:        utils.GetEnvOr("PORT", "4431"),
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

func initialize() {
	userService = UserService{
		Host: utils.GetEnvOr("HANDS_ON_USER_SERVICE_SERVICE_HOST", "hands-on-user-service"),
		Port: utils.GetEnvOr("HANDS_ON_USER_SERVICE_SERVICE_PORT", "4432"),
	}
}

func handlers() []server.Handler {
	return []server.Handler{
		{
			// http POST localhost:4431/users first_name='Ronen' last_name='Konfortes' email='konfortes@gmail.com'
			Method:  http.MethodPost,
			Pattern: "/users",
			H:       usersHandler,
		},
	}
}
