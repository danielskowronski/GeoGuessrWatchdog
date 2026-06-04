package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	dbConst "github.com/danielskowronski/geoguessrwatchdog/internal/db"
	db "github.com/danielskowronski/geoguessrwatchdog/internal/db/generated"
)

type FetchStatusInput struct {
	MaxAgeDays int `query:"maxAgeDays" default:"7"`
}

type FetchStatus struct {
	FetchType   string    `json:"fetchType"`
	Description string    `json:"description,omitempty"`
	LastSuccess time.Time `json:"lastSuccess"`
}
type FetchStatusOutput struct {
	Body struct {
		Fetches []FetchStatus `json:"fetches"`
	}
}

func getFetchStatusDescription(fetchType string, userAliases map[string]string) string {
	if fetchType == dbConst.FetchTypeDivisionsAndAllRelatedMaps {
		return "All maps for current divisions and game modes"
	}
	if strings.HasPrefix(fetchType, dbConst.FetchTypeAllUserStatsPrefix) {
		userID := strings.TrimPrefix(fetchType, dbConst.FetchTypeAllUserStatsPrefix)
		if alias, ok := userAliases[userID]; ok {
			return fmt.Sprintf("All stats for known user %s", alias)
		}
		return fmt.Sprintf("All stats for user %s", userID)
	}
	return "unknown"
}

func (a *App) GetFetchStatuses(ctx context.Context, input *FetchStatusInput) (*FetchStatusOutput, error) {
	a.logger.Debug("handling GetFetchStatuses")
	resp := &FetchStatusOutput{}
	resp.Body.Fetches = make([]FetchStatus, 0)

	q := db.New(a.db)
	fetchStatuses, err := q.GetAllStatuses(ctx, int32(input.MaxAgeDays))
	if err != nil {
		a.logger.Warn("failed to get fetch statuses", "err", err)
		return nil, huma.Error500InternalServerError("failed to get fetch statuses")
	}

	for _, row := range fetchStatuses {
		resp.Body.Fetches = append(resp.Body.Fetches, FetchStatus{
			FetchType:   row.FetchType,
			Description: getFetchStatusDescription(row.FetchType, a.userAliases),
			LastSuccess: row.LastSuccess.Time,
		})
	}

	return resp, nil
}
