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

func main() {
	initialize()

	router := gin.Default()

	serverutils.SetMiddlewares(router, tracer)
	serverutils.SetMonitoringHandler(router)
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
		tracer = serverutils.InitJaeger("hands-on-k8s-web")
	}
}

func setRoutes(router *gin.Engine) {
	// http localhost:8080/health
	router.GET("/health", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", []byte("OK"))
	})

	router.POST("/users", usersHandler)
}
