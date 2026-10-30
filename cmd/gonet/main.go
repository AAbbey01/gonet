package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/AAbbey01/gonet/internal/check"
	"github.com/AAbbey01/gonet/internal/metrics"
)

func main() {
	configPath := flag.String("config", "targets.yaml", "path to the targets file")
	addr := flag.String("addr", ":8080", "listen address for /metrics")
	flag.Parse()

	cfg, err := check.Load(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	reg := prometheus.NewRegistry()
	m := metrics.New(reg)

	// ctx is cancelled on Ctrl+C or docker stop.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	for _, t := range cfg.Targets {
		go watch(ctx, cfg, t, m)
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	srv := &http.Server{Addr: *addr, Handler: mux}
	go func() {
		<-ctx.Done()
		srv.Shutdown(context.Background())
	}()

	log.Printf("checking %d targets every %s; metrics on %s/metrics", len(cfg.Targets), cfg.Interval, *addr)
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

// watch checks one target every cfg.Interval until ctx is cancelled.
func watch(ctx context.Context, cfg check.Config, t check.Target, m *metrics.Metrics) {
	// TODO 1: ticker := time.NewTicker(cfg.Interval); defer ticker.Stop()
	// TODO 2: loop forever:
	//   cctx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	//   r := check.Run(cctx, t); cancel()
	//   m.Record(r); log.Printf one line with name, kind, r.Up, r.Latency, r.Err
	//   select { case <-ctx.Done(): return; case <-ticker.C: }
	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	for {
		cctx, cancel := context.WithTimeout(ctx, cfg.Timeout)
		r := check.Run(cctx, t)
		cancel()
		m.Record(r)
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:

		}
	}
}
