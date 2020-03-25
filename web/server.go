package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
)

var (
	tracer            opentracing.Tracer
	tracerCloser      io.Closer
	shutdownHooks     []func()
	customMiddlewares []gin.HandlerFunc
	userService       UserService
)

func main() {
	initialize()

	router := gin.Default()

	setMiddlewares(router)
	setRoutes(router)

	srv := &http.Server{
		Addr:    ":" + getEnvOr("PORT", "4431"),
		Handler: router,
	}

	go func() {
		log.Println("listening on " + srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	gracefulShutdown(srv)
}

func initialize() {
	userService = UserService{
		Host: getEnvOr("HANDS_ON_USER_SERVICE_SERVICE_HOST", "hands-on-user-service"),
		Port: getEnvOr("HANDS_ON_USER_SERVICE_SERVICE_PORT", "4432"),
	}

	if isTracingEnabled() {
		initJaeger("hands-on-k8s-web")
		customMiddlewares = append(customMiddlewares, jaegerMiddleware)
	}
}

func setMiddlewares(router *gin.Engine) {
	// TODO: add jaegerMiddleware only to relevant routes
	for _, middleware := range customMiddlewares {
		router.Use(middleware)
	}

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
}

func setRoutes(router *gin.Engine) {
	// http localhost:8080/health
	router.GET("/health", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", []byte("OK"))
	})

	router.POST("/users", usersHandler)
}

func getEnvOr(env, ifNotFound string) string {
	foundEnv, found := os.LookupEnv(env)

	if found {
		return foundEnv
	}

	return ifNotFound
}

func gracefulShutdown(srv *http.Server) {
	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shuting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for _, hook := range shutdownHooks {
		hook()
	}
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
