package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// TODO: write error to trace?
		return
	}

	processUser(c.Request.Context(), &input)

	if err := persistUser(c.Request.Context(), input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "created"})
}

func decodeInput(ctx context.Context, req *http.Request) (UserInput, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "decodeInput")
	defer span.Finish()

	userInput := UserInput{}
	if err := json.NewDecoder(req.Body).Decode(&userInput); err != nil {
		span.LogFields(
			traceLog.Error(err),
		)
		span.SetTag("status", "error")
		return userInput, err
	}

	span.LogFields(
		traceLog.String("event", "userDecoded"),
		traceLog.String("value", fmt.Sprintf("%+v", userInput)),
	)

	return userInput, nil
}

func processUser(ctx context.Context, user *UserInput) {
	span, _ := opentracing.StartSpanFromContext(ctx, "processUser")
	defer span.Finish()

	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)
	user.Email = strings.TrimSpace(user.Email)
}

func persistUser(ctx context.Context, user UserInput) error {
	span, persistCtx := opentracing.StartSpanFromContext(ctx, "persistUser")
	defer span.Finish()

	err := userService.CreateUser(persistCtx, user)
	if err != nil {
		span.LogFields(
			traceLog.String("event", "error"),
			traceLog.String("value", err.Error()),
		)
		span.SetTag("status", "error")
	}
	return err
}
