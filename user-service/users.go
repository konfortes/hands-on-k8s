package main

import (
	"context"
	"errors"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
		span := opentracing.SpanFromContext(c)
		span.SetTag("error", true)
		span.LogFields(
			traceLog.String("event", "error"),
			traceLog.String("message", err.Error()),
		)
		return
	}

	if err := saveUser(c.Request.Context(), input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "created"})
}

func saveUser(ctx context.Context, user UserInput) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "saveUser")
	defer span.Finish()

	// sleep random time
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(1000)
	time.Sleep(time.Duration(n) * time.Millisecond)

	// randomally return error
	if n > 800 {
		span.SetTag("error", true)
		span.LogFields(
			traceLog.String("event", "error"),
			traceLog.String("message", "just a random error while saving user"),
		)
		return errors.New("error saving user")
	}
	return nil
}
