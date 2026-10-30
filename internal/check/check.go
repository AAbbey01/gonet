package check

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-playground/locales/or"
)

// Result is the outcome of one check.
type Result struct {
	Target  Target
	Up      bool
	Latency time.Duration
	Err     error
}

// Run checks one target and times it. ctx carries the timeout.
func Run(ctx context.Context, t Target) Result {
	start := time.Now()
	var err error
	switch t.Kind {
	case "tcp":
		err = checkTCP(ctx, t.Address)
	case "http":
		err = checkHTTP(ctx, t.Address)
	case "dns":
		err = checkDNS(ctx, t.Address)
	default:
		err = fmt.Errorf("unknown kind %q", t.Kind)
	}
	return Result{Target: t, Up: err == nil, Latency: time.Since(start), Err: err}
}

func checkTCP(ctx context.Context, addr string) error {
	// TODO: var d net.Dialer; conn, err := d.DialContext(ctx, "tcp", addr)
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", addr)
	// on success close the connection and return nil
	if err != nil{
		//conn.Close()
		return fmt.Errorf("Error connecting %w",err)
	}
	conn.Close()
	return nil
}

func checkHTTP(ctx context.Context, url string) error {
	// TODO: req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	// resp, err := http.DefaultClient.Do(req); always defer resp.Body.Close()
	// return an error when resp.StatusCode >= 400
	req, err := http.NewRequestWithContext(ctx,http.MethodGet, url, nil)
	resp, err := http.DefaultClient.Do(req)
	if resp.StatusCode >= 400{
		return fmt.Errorf("Error code %d: %w",resp.StatusCode,err)
	}
	
	defer resp.Body.Close()

	return nil
}

func checkDNS(ctx context.Context, host string) error {
	// TODO: addrs, err := net.DefaultResolver.LookupHost(ctx, host)
	// return an error when err != nil or len(addrs) == 0
	addrs, err := net.DefaultResolver.LookupHost(ctx,host)
	if err!= nil || len(addrs) == 0{
		return fmt.Errorf("Error with dns lookup: %w", err)
	}
	return nil
}