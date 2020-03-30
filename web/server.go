package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/konfortes/go-server-utils/serverutils"
	opentracing "github.com/opentracing/opentracing-go"
)

var (
	tracer      *opentracing.Tracer
	userService UserService
)

const (
	serviceName = "hands-on-k8s-web"
)

func main() {
	initialize()

	router := gin.Default()

	serverutils.SetMiddlewares(router, tracer)
	serverutils.SetRoutes(router, serviceName)
	setRoutes(router)

	srv := &http.Server{
		Addr:    ":" + serverutils.GetEnvOr("PORT", "4431"),
		Handler: router,
	}

	go func() {
		log.Println("listening on " + srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	serverutils.GracefulShutdown(srv)
}

func initialize() {
	userService = UserService{
		Host: serverutils.GetEnvOr("HANDS_ON_USER_SERVICE_SERVICE_HOST", "hands-on-user-service"),
		Port: serverutils.GetEnvOr("HANDS_ON_USER_SERVICE_SERVICE_PORT", "4432"),
	}

	if serverutils.GetEnvOr("TRACING_ENABLED", "false") == "true" {
		tracer = serverutils.InitJaeger(serviceName)
	}
}

func setRoutes(router *gin.Engine) {
	router.POST("/users", usersHandler)
}
