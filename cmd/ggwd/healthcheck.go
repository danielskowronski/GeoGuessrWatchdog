package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/danielskowronski/geoguessrwatchdog/internal/buildinfo"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.temporal.io/sdk/client"
)

type HealthState struct {
	started      atomic.Bool
	shuttingDown atomic.Bool
}

func startHealthServer(
	ctx context.Context,
	addr string,
	db *pgxpool.Pool,
	tc client.Client,
	state *HealthState,
) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(buildinfo.GetBuildInfo() + "\n"))
	})

	mux.HandleFunc("/livez", func(w http.ResponseWriter, r *http.Request) {
		if state.shuttingDown.Load() {
			http.Error(w, "shutting down", http.StatusServiceUnavailable)
			return
		}
		if !state.started.Load() {
			http.Error(w, "worker not started", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})

	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if !state.started.Load() {
			http.Error(w, "worker not started", http.StatusServiceUnavailable)
			return
		}
		if state.shuttingDown.Load() {
			http.Error(w, "shutting down", http.StatusServiceUnavailable)
			return
		}

		checkCtx, cancel := context.WithTimeout(r.Context(), 500*time.Millisecond)
		defer cancel()

		if err := db.Ping(checkCtx); err != nil {
			http.Error(w, "db not ready", http.StatusServiceUnavailable)
			return
		}

		if _, err := tc.CheckHealth(checkCtx, &client.CheckHealthRequest{}); err != nil {
			http.Error(w, "temporal not ready", http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready\n"))
	})

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		<-ctx.Done()
		state.shuttingDown.Store(true)

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_ = srv.Shutdown(shutdownCtx)
	}()

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("health server failed", "error", err)
		}
	}()

	return srv
}
