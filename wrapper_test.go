package apigateway

import (
	"net/http"
	"testing"
)

func TestNewRequest(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}
	if req.Header == nil {
		t.Fatalf("got nil; want not nil")
	}
}
