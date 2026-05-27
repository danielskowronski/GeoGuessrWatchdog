package main

import (
	"context"
	"time"

	"github.com/danielgtaylor/huma/v2"
	db "github.com/danielskowronski/geoguessrwatchdog/internal/db/generated"
)

type DivisionsOutput struct {
	Body struct {
		Divisions []DivisionModeMapInfo `json:"divisions"`
	}
}

type DivisionModeMapInfo struct {
	DivisionName string    `json:"divisionName"`
	GameMode     string    `json:"gameMode"`
	MapID        string    `json:"mapId"`
	MapName      string    `json:"mapName"`
	LastChanged  time.Time `json:"lastChanged"`
}

func (a *App) GetDivisions(ctx context.Context, input *struct{}) (*DivisionsOutput, error) {
	a.logger.Debug("handling GetDivisions")
	q := db.New(a.db)
	divisions, err := q.GetAllDivisionsInfo(ctx)
	if err != nil {
		a.logger.Warn("failed to get divisions", "err", err)
		return nil, huma.Error500InternalServerError("failed to get divisions")
	}

	resp := &DivisionsOutput{}
	resp.Body.Divisions = make([]DivisionModeMapInfo, len(divisions))
	for i := range divisions {
		resp.Body.Divisions[i] = DivisionModeMapInfo{
			DivisionName: divisions[i].DivisionName,
			GameMode:     divisions[i].GameMode,
			MapID:        divisions[i].MapID,
			MapName:      divisions[i].MapName,
			LastChanged:  divisions[i].LastChanged.Time,
		}
	}
	return resp, nil
}
