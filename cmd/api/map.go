package main

import (
	"context"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielskowronski/geoguessrwatchdog/internal/apischema"
	db "github.com/danielskowronski/geoguessrwatchdog/internal/db/generated"
	"github.com/jackc/pgx/v5"
)

type MapHistoryInput struct {
	ID string `query:"id" doc:"Map ID"`
}
type MapHistoryOutput struct {
	Body struct {
		Name        string                `json:"name"`
		Description string                `json:"description"`
		URL         string                `json:"url"`
		Entries     []db.GetMapHistoryRow `json:"entries"`
	}
}

func (a *App) GetMapHistory(ctx context.Context, input *MapHistoryInput) (*MapHistoryOutput, error) {
	resp := &MapHistoryOutput{}

	q := db.New(a.db)
	info, err := q.GetMapInfo(ctx, input.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, huma.Error404NotFound("map not found")
		}
		fmt.Printf("failed to get map history: %v\n", err)
		return nil, huma.Error500InternalServerError("failed to get map history")
	}
	resp.Body.Name = info.Name
	resp.Body.Description = info.Description.String
	resp.Body.URL = fmt.Sprintf(apischema.MAP_FRONTEND_URL, input.ID)

	history, err := q.GetMapHistory(ctx, input.ID)
	if err != nil {
		fmt.Printf("failed to get map history: %v\n", err)
		return nil, huma.Error500InternalServerError("failed to get map history")
	}

	resp.Body.Entries = history
	return resp, nil
}
