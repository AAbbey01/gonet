package check

import (
	"context"
	"net"
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

func TestCheckTCP(t *testing.T){
	ln, _ := net.Listen("tcp","127.0.0.1:0")
	if len(ln.Addr().String()) != 0{
		ln.Close()
	}else{
		t.Fatal("Failed to connect TCP first time")
	}
	if len(ln.Addr().String()) >0{
		t.Fatal("Connection was not closed")
	}
}

func TestRunTimeout(t *testing.T){
	
}