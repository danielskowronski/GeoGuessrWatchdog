package main

import (
	"context"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	db "github.com/danielskowronski/geoguessrwatchdog/internal/db/generated"
)

type DivisionsOutput struct {
	Body struct {
		Divisions []db.DivisionInfo `json:"divisions"`
	}
}

func (a *App) GetDivisions(ctx context.Context, input *struct{}) (*DivisionsOutput, error) {
	q := db.New(a.db)
	divisions, err := q.GetAllDivisionsInfo(ctx)
	if err != nil {
		fmt.Printf("failed to get divisions: %v\n", err)
		return nil, huma.Error500InternalServerError("failed to get divisions")
	}

	resp := &DivisionsOutput{}
	resp.Body.Divisions = divisions
	return resp, nil
}
