package e2e

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"
)

type TestClient struct {
	BaseURL string
	Client  *http.Client
	T       *testing.T
}

type Response struct {
	Body       string
	StatusCode int
}

func NewTestClient(t *testing.T, baseURL string) *TestClient {
	return &TestClient{
		BaseURL: baseURL,
		Client:  &http.Client{},
		T:       t,
	}
}

func (tc *TestClient) Post(path string, payload []byte) Response {
	url := fmt.Sprintf("%s%s", tc.BaseURL, path)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		tc.T.Fatalf("Failed to create POST request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	return tc.do(req)
}

func (tc *TestClient) Get(path string) Response {
	url := fmt.Sprintf("%s%s", tc.BaseURL, path)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		tc.T.Fatalf("Failed to create GET request: %v", err)
	}

	return tc.do(req)
}

func (tc *TestClient) Delete(path string) Response {
	url := fmt.Sprintf("%s%s", tc.BaseURL, path)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		tc.T.Fatalf("Failed to create DELETE request: %v", err)
	}

	return tc.do(req)
}

func (tc *TestClient) do(req *http.Request) Response {
	res, err := tc.Client.Do(req)
	if err != nil {
		tc.T.Fatalf("Failed to execute request: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		tc.T.Fatalf("Failed to read response body: %v", err)
	}

	return Response{
		Body:       string(body),
		StatusCode: res.StatusCode,
	}
}
