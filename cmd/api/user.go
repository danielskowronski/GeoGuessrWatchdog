package main

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	db "github.com/danielskowronski/geoguessrwatchdog/internal/db/generated"
	"github.com/jackc/pgx/v5"
)

type UsersOutput struct {
	Body struct {
		Users map[string]string `json:"users"`
	}
}

func (a *App) GetUsers(ctx context.Context, input *struct{}) (*UsersOutput, error) {
	resp := &UsersOutput{}
	resp.Body.Users = a.userAliases
	return resp, nil
}

type UserStatsInput struct {
	ID string `query:"id" doc:"User ID"`
}
type UserStatsOutput struct {
	Body struct {
		Entries []db.GetUserFetchCombinedHistoryRow `json:"entries"`
	}
}

func (a *App) GetUserStats(ctx context.Context, input *UserStatsInput) (*UserStatsOutput, error) {
	resp := &UserStatsOutput{}

	q := db.New(a.db)
	history, err := q.GetUserFetchCombinedHistory(ctx, input.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, huma.Error404NotFound("user not found")
		}
		return nil, huma.Error500InternalServerError("failed to get user stats")
	}

	resp.Body.Entries = history
	return resp, nil
}
