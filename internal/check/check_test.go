package check

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
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
	ln, err := net.Listen("tcp","127.0.0.1:0")
	if err != nil{
		t.Fatal(err)
	}
	addr := ln.Addr().String()
	if err:= checkTCP(context.Background(), addr); err != nil{
		t.Fatalf("open port: want nil, got %v", err)
	}
	ln.Close()
	if err := checkTCP(context.Background(), addr); err == nil{
		t.Fatal("closed port: want an error, got nil")
	}
}

func TestRunTimeout(t *testing.T){
	slow := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                time.Sleep(200 * time.Millisecond)
        }))
	defer slow.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	r := Run(ctx, Target{Name: "slow", Kind: "http", Address: slow.URL})
	if r.Up{
		t.Fatal("want down after timeout, got up")
	}
	if r.Latency > 150*time.Millisecond{
		t.Fatalf("check took %v; the 50ms deadline was not enforced",r.Latency)
	}
}

func TestLoad(t *testing.T){
	path := filepath.Join(t.TempDir(),"targets.yaml")
	good := "interval: 10s\ntargets:\n  - name: a\n    kind: tcp\n    address: 127.0.0.1:1\n"
	if err:= os.WriteFile(path, []byte(good), 0o644); err != nil{
		t.Fatal(err)
	}
	cfg, err := Load(path)
	if err != nil{
		t.Fatalf("valid file: want nil, got %v",err)
	}
	if cfg.Interval != 10*time.Second{
		t.Errorf("Interval = %v, want 10s", cfg.Interval)
	}
	if cfg.Timeout != 5*time.Second{
		t.Errorf("Timeout = %v, want 5s", cfg.Timeout)
	}
	if len(cfg.Targets) != 1{
		t.Errorf("got %d targets, want 1", len(cfg.Targets))
		
	}
	bad := "targets:\n  - name: a\n    kind: ftp\n    address: x\n"
	if err := os.WriteFile(path, []byte(bad), 0o644); err != nil {
			t.Fatal(err)
	}
	if _, err := Load(path); err == nil {
			t.Fatal("kind ftp: want an error, got nil")
	}

	empty := "targets: []\n"
	if err := os.WriteFile(path, []byte(empty), 0o644); err != nil {
			t.Fatal(err)
	}
	if _, err := Load(path); err == nil {
			t.Fatal("no targets: want an error, got nil")
	}
}