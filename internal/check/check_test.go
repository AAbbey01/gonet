package check

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckHTTP(t *testing.T) {
	ok := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ok.Close()
	broken := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer broken.Close()

	if err := checkHTTP(context.Background(), ok.URL); err != nil {
		t.Fatalf("200 server: want nil, got %v", err)
	}
	if err := checkHTTP(context.Background(), broken.URL); err == nil {
		t.Fatal("500 server: want an error, got nil")
	}
}
