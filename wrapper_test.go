package apigateway

import (
	"encoding/json"
	"io/ioutil"
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

func Test_makeV2Request(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/request-get.json")
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}

	var event Request
	err = json.Unmarshal(data, &event)
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}

	req, err := makeV2Request(event, "/dev")
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}

	if want, got := "v6rlfj8f7b.execute-api.us-west-2.amazonaws.com", req.Host; got != want {
		t.Fatalf("got %v; want %v", got, want)
	}
	if want, got := "73.189.109.118", req.RemoteAddr; got != want {
		t.Fatalf("got %v; want %v", got, want)
	}
	if want, got := "/abc", req.URL.Path; got != want {
		t.Fatalf("got %v; want %v", got, want)
	}
	if want, got := "/abc", req.RequestURI; got != want {
		t.Fatalf("got %v; want %v", got, want)
	}
	if want, got := int64(4), req.ContentLength; got != want {
		t.Fatalf("got %v; want %v", got, want)
	}
	if want, got := http.MethodGet, req.Method; got != want {
		t.Fatalf("got %v; want %v", got, want)
	}
	if want, got := "/abc", req.URL.Path; got != want {
		t.Fatalf("got %v; want %v", got, want)
	}
	got, err := req.Cookie("abc")
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}
	if want := "123"; want != got.Value {
		t.Fatalf("got %v; want %v; err %v", got, want, err)
	}
}

func Test_makeV1Request(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/request-v1-post.json")
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}

	var event Request
	err = json.Unmarshal(data, &event)
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}

	req, err := makeV1Request(event)
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}

	if got, want := "dev-d-loadb-1xcoy0h9gw154-247507225.us-west-2.elb.amazonaws.com", req.Host; got != want {
		t.Fatalf("got %v; want %v", got, want)
	}
	if got, want := "73.189.109.118", req.RemoteAddr; got != want {
		t.Fatalf("got %v; want %v", got, want)
	}
	if got, want := "/graphql", req.URL.Path; got != want {
		t.Fatalf("got %v; want %v", got, want)
	}
	if got, want := "/graphql", req.RequestURI; got != want {
		t.Fatalf("got %v; want %v", got, want)
	}
	if got, want := int64(4), req.ContentLength; got != want {
		t.Fatalf("got %v; want %v", got, want)
	}
	if got, want := http.MethodPost, req.Method; got != want {
		t.Fatalf("got %v; want %v", got, want)
	}
	if got, want := "/graphql", req.URL.Path; got != want {
		t.Fatalf("got %v; want %v", got, want)
	}
}

func Test_makeV1RequestMultipart(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/request-v1-multipart.json")
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}

	var event Request
	err = json.Unmarshal(data, &event)
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}

	req, err := makeV1Request(event)
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}

	err = req.ParseMultipartForm(1e6)
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}

	if req.MultipartForm == nil {
		t.Fatalf("got nil; want not nil")
	}
}
