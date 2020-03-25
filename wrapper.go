package apigateway

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

type RequestContext struct {
	AccountID    string `json:"accountId"`
	ApiId        string `json:"apiId"`
	DomainName   string `json:"domainName"`
	DomainPrefix string `json:"domainPrefix"`
	Http         struct {
		Method    string `json:"method"`
		Path      string `json:"path"`
		Protocol  string `json:"protocol"`
		SourceIp  string `json:"sourceIp"`
		UserAgent string `json:"userAgent"`
	} `json:"http"`
	QueryStringParameters map[string]string `json:"queryStringParameters"`
	RequestId             string            `json:"requestId"`
	RouteKey              string            `json:"routeKey"`
	Stage                 string            `json:"stage"`
	Time                  string            `json:"time"`
	TimeEpoch             int64             `json:"timeEpoch"`
}

type Request struct {
	Version         string            `json:"version"`
	RouteKey        string            `json:"routeKey"`
	RawPath         string            `json:"rawPath"`
	RawQueryString  string            `json:"rawQueryString"`
	Headers         map[string]string `json:"headers"`
	RequestContext  RequestContext    `json:"requestContext"`
	StageVariables  map[string]string `json:"stageVariables"`
	Body            string            `json:"body"`
	IsBase64Encoded bool              `json:"isBase64Encoded"`
}

type Response struct {
	StatusCode      int               `json:"statusCode"`
	Headers         map[string]string `json:"headers"`
	Body            string            `json:"body"`
	IsBase64Encoded bool              `json:"isBase64Encoded,omitempty"`
}

func Wrap(handler http.Handler) func(ctx context.Context, event Request) (Response, error) {
	return func(ctx context.Context, event Request) (Response, error) {
		req, err := makeRequest(event)
		if err != nil {
			return Response{}, err
		}

		for k, v := range event.Headers {
			if n := strings.Index(v, ","); n > 0 {
				req.Header[k] = strings.Split(v, ",")
				continue
			}
			req.Header[k] = []string{v}
		}

		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		resp, err := makeResponse(w)
		if err != nil {
			return Response{}, err
		}

		return resp, nil
	}
}

func makeRequest(event Request) (*http.Request, error) {
	var uri string
	switch event.RawQueryString {
	case "":
		uri = "http://" + event.RequestContext.DomainName + event.RawPath
	default:
		uri = "http://" + event.RequestContext.DomainName + event.RawPath + "?" + event.RawQueryString
	}

	var body io.Reader
	switch {
	case event.Body != "" && event.IsBase64Encoded:
		data, err := base64.StdEncoding.DecodeString(event.Body)
		if err != nil {
			return nil, fmt.Errorf("unable to base64 decode request body: %w", err)
		}
		body = bytes.NewReader(data)

	case event.Body != "":
		body = strings.NewReader(event.Body)

	default:
		body = bytes.NewReader(nil)
	}

	req, err := http.NewRequest(event.RequestContext.Http.Method, uri, body)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %w", err)
	}

	return req, nil
}

func makeResponse(w *httptest.ResponseRecorder) (Response, error) {
	var body string
	if raw := w.Body.Bytes(); len(raw) > 0 {
		body = base64.StdEncoding.EncodeToString(raw)
	}

	headers := map[string]string{}
	for k, v := range w.Header() {
		switch len(v) {
		case 0:
			// ok
		case 1:
			headers[k] = v[0]
		default:
			headers[k] = strings.Join(v, ",")
		}
	}

	return Response{
		StatusCode:      w.Code,
		Headers:         headers,
		Body:            body,
		IsBase64Encoded: len(body) > 0,
	}, nil
}
