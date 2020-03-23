package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	traceLog "github.com/opentracing/opentracing-go/log"
)

// UserService ...
type UserService struct {
	Host string
	Port string
}

// CreateUser ...
func (us UserService) CreateUser(ctx context.Context, user UserInput) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "UserService#CreateUser")

	requestBody, err := json.Marshal(user)
	if err != nil {
		span.LogFields(
			traceLog.Error(err),
		)
		span.SetTag("status", "error")
		return err
	}

	url := fmt.Sprintf("http://%s:%s/users", us.Host, us.Port)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		span.LogFields(
			traceLog.Error(err),
		)
		span.SetTag("status", "error")
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: time.Second * 10}

	opentracing.GlobalTracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header))

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
		span.LogFields(
			traceLog.Error(err),
		)
		span.SetTag("status", "error")
		return err
	}

	span.LogFields(
		traceLog.String("responseStatus", resp.Status),
	)

	if resp.StatusCode > 400 {
		span.LogFields(
			traceLog.String("responseError", resp.Status),
		)
	}

	return nil
}
