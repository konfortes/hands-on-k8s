package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	opentracing "github.com/opentracing/opentracing-go"
	traceLog "github.com/opentracing/opentracing-go/log"
)

// UserInput ...
type UserInput struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func usersHandler(w http.ResponseWriter, req *http.Request) {
	tracer, closer := initJaeger("hands-on-k8s-web")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	span := tracer.StartSpan("userHandler")
	defer span.Finish()
	ctx := context.Background()
	ctx = opentracing.ContextWithSpan(ctx, span)

	userInput, err := decodeInput(ctx, req)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid input"))
		return
	}

	processUser(ctx, &userInput)

	if err := persistUser(ctx, userInput); err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func decodeInput(ctx context.Context, req *http.Request) (UserInput, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "decodeInput")
	defer span.Finish()

	span.LogKV("event", "decodeInput")
	userInput := UserInput{}
	if err := json.NewDecoder(req.Body).Decode(&userInput); err != nil {
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
	span, _ := opentracing.StartSpanFromContext(ctx, "persistUser")
	defer span.Finish()

	err := userService.CreateUser(user)
	if err != nil {
		span.LogFields(
			traceLog.String("event", "error"),
			traceLog.String("value", err.Error()),
		)
	}
	return err
}
