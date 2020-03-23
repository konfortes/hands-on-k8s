package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	traceLog "github.com/opentracing/opentracing-go/log"
)

// UserInput ...
type UserInput struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func handleUsers(w http.ResponseWriter, req *http.Request) {
	tracer, closer := initJaeger("hands-on-k8s-user-service")
	defer closer.Close()

	upstreamCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	span := tracer.StartSpan("handleUsers", ext.RPCServerOption(upstreamCtx))
	defer span.Finish()

	ctx := context.Background()
	ctx = opentracing.ContextWithSpan(ctx, span)

	userInput := UserInput{}
	if err := decodeInput(ctx, req, &userInput); err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := saveUser(ctx, userInput); err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func decodeInput(ctx context.Context, req *http.Request, user *UserInput) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "decodeInput")
	defer span.Finish()

	if err := json.NewDecoder(req.Body).Decode(user); err != nil {
		span.LogFields(
			traceLog.Error(err),
		)
		span.SetTag("status", "error")
		return err
	}

	return nil
}

func saveUser(ctx context.Context, user UserInput) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "saveUser")
	defer span.Finish()

	time.Sleep(time.Millisecond * 1000)

	return nil
}
