package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/konfortes/go-server-utils/tracing"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
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
	defer span.Finish()

	requestBody, err := json.Marshal(user)
	if err != nil {
		tracing.Error(span, err)
		return err
	}

	url := fmt.Sprintf("http://%s:%s/users", us.Host, us.Port)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: time.Second * 10}
	// recommended opentracing semantic conventions tags
	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, url)
	ext.HTTPMethod.Set(span, "POST")

	opentracing.GlobalTracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header))

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
		tracing.Error(span, err)
		return err
	}

	span.LogFields(
		traceLog.String("responseStatus", resp.Status),
	)

	if resp.StatusCode > 400 {
		span.SetTag("error", true)
		span.LogFields(
			traceLog.String("event", "error"),
			traceLog.Int("http.response_code", resp.StatusCode),
		)
		return fmt.Errorf("got %d from user-service", resp.StatusCode)
	}

	return nil
}
