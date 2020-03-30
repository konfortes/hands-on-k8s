package main

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
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

	if err := saveUser(c.Request.Context(), input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "created"})
}

func saveUser(ctx context.Context, user UserInput) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "saveUser")
	defer span.Finish()

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(1000)
	time.Sleep(time.Duration(n) * time.Millisecond)

	return nil
}
