package main

import (
	"context"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/konfortes/go-server-utils/tracing"
	opentracing "github.com/opentracing/opentracing-go"
	traceLog "github.com/opentracing/opentracing-go/log"
)

// UserInput ...
type UserInput struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" binding:"required"`
}

func usersHandler(c *gin.Context) {
	var input UserInput
	if err := c.BindJSON(&input); err != nil {
		span := opentracing.SpanFromContext(c.Request.Context())
		tracing.Error(span, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	processUser(c.Request.Context(), &input)

	// some async op
	go func(ctx context.Context) {
		span, _ := opentracing.StartSpanFromContext(ctx, "AsyncOp")
		defer span.Finish()

		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(1000)
		time.Sleep(time.Duration(n) * time.Millisecond)
	}(c.Request.Context())

	if err := persistUser(c.Request.Context(), input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "created"})
}

func processUser(ctx context.Context, user *UserInput) {
	span, _ := opentracing.StartSpanFromContext(ctx, "processUser")
	defer span.Finish()

	span.LogFields(
		traceLog.String("processing", user.Email),
	)

	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)
	user.Email = strings.TrimSpace(user.Email)
}

func persistUser(ctx context.Context, user UserInput) error {

	if err := userService.CreateUser(ctx, user); err != nil {
		span := opentracing.SpanFromContext(ctx)
		tracing.Error(span, err)
		return err
	}
	return nil
}
