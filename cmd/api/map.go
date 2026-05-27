package main

import (
	"context"
	"fmt"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielskowronski/geoguessrwatchdog/internal/apischema"
	db "github.com/danielskowronski/geoguessrwatchdog/internal/db/generated"
	"github.com/jackc/pgx/v5"
)

type MapHistoryInput struct {
	ID string `path:"id" doc:"Map ID"`
}

type MapHistoryEntry struct {
	Timestamp time.Time         `json:"timestamp"`
	Info      apischema.MapInfo `json:"info"`
}

type MapHistoryOutput struct {
	Body struct {
		Name        string            `json:"name"`
		Description string            `json:"description"`
		URL         string            `json:"url"`
		Entries     []MapHistoryEntry `json:"entries"`
	}
}

func (a *App) GetMapHistory(ctx context.Context, input *MapHistoryInput) (*MapHistoryOutput, error) {
	a.logger.Debug("handling GetMapHistory", "map_id", input.ID)
	resp := &MapHistoryOutput{}

	q := db.New(a.db)
	info, err := q.GetMapInfo(ctx, input.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, huma.Error404NotFound("map not found")
		}
		a.logger.Warn("failed to get map history", "err", err)
		return nil, huma.Error500InternalServerError("failed to get map history")
	}
	resp.Body.Name = info.Name
	resp.Body.Description = info.Description.String
	resp.Body.URL = fmt.Sprintf(apischema.MAP_FRONTEND_URL, input.ID)

	history, err := q.GetMapHistory(ctx, input.ID)
	if err != nil {
		a.logger.Warn("failed to get map history", "err", err)
		return nil, huma.Error500InternalServerError("failed to get map history")
	}

	resp.Body.Entries = make([]MapHistoryEntry, len(history))
	for i := range history {
		resp.Body.Entries[i] = MapHistoryEntry{
			Timestamp: unwrapTimestamptzDefault(history[i].Timestamp),
			Info: apischema.MapInfo{
				ID:               info.ApiID,
				BoundsMinLat:     unwrapFloat8Default(history[i].BoundsMinLat),
				BoundsMinLng:     unwrapFloat8Default(history[i].BoundsMinLon),
				BoundsMaxLat:     unwrapFloat8Default(history[i].BoundsMaxLat),
				BoundsMaxLng:     unwrapFloat8Default(history[i].BoundsMaxLon),
				MaxErrorDistance: unwrapInt8Default(history[i].MaxErrDistance),
				UpdatedAt:        unwrapTimestamptzDefault(history[i].ApiUpdated),
				CoordinateCount:  history[i].LocationCount,
			},
		}
	}
	return resp, nil
}
