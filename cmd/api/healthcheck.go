package main

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/danielgtaylor/huma/v2"
	buildInfo "github.com/danielskowronski/geoguessrwatchdog/internal/buildinfo"
)

type HealthState struct {
	started      atomic.Bool
	shuttingDown atomic.Bool
}

type EmptyInput struct{}

type HealthcheckOutput struct {
	Body struct {
		Status string `json:"status"`
	}
}
type VersionOutput struct {
	Body struct {
		Version   string `json:"version"`
		BuildDate string `json:"buildDate"`
	}
}

func (a *App) Version(ctx context.Context, input *EmptyInput) (*VersionOutput, error) {
	resp := &VersionOutput{}
	resp.Body.Version = buildInfo.Version
	resp.Body.BuildDate = buildInfo.BuildDate
	return resp, nil
}

func (a *App) Livez(ctx context.Context, input *EmptyInput) (*HealthcheckOutput, error) {
	resp := &HealthcheckOutput{}

	if a.state.shuttingDown.Load() {
		resp.Body.Status = "shutting down"
		return resp, huma.Error500InternalServerError("shutting down")
	}
	if !a.state.started.Load() {
		resp.Body.Status = "worker not started"
		return resp, huma.Error500InternalServerError("worker not started")
	}

	resp.Body.Status = "ok"
	return resp, nil
}

func (a *App) Readyz(ctx context.Context, input *struct{}) (*HealthcheckOutput, error) {
	resp := &HealthcheckOutput{}

	if a.state.shuttingDown.Load() {
		return resp, huma.Error500InternalServerError("shutting down")
	}
	if !a.state.started.Load() {
		return resp, huma.Error500InternalServerError("worker not started")
	}

	checkCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	if err := a.db.Ping(checkCtx); err != nil {
		return resp, huma.Error500InternalServerError("db not ready")
	}

	resp.Body.Status = "ok"
	return resp, nil
}
