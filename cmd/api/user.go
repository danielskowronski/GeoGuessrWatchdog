package main

import (
	"context"
	"time"

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
	a.logger.Debug("handling GetUsers")
	resp := &UsersOutput{}
	resp.Body.Users = a.userAliases
	// TODO: probably need to add unique user from DB without aliases as well
	return resp, nil
}

type UserStatsInput struct {
	ID string `path:"id" doc:"User ID"`
}

type UserStatsDailyEntry struct {
	FetchID                 int64     `json:"fetchID"`
	FetchTimestamp          time.Time `json:"fetchTimestamp"`
	DivisionName            string    `json:"divisionName"`
	RatingOverall           int32     `json:"ratingOverall"`
	RatingMoving            int32     `json:"ratingMoving"`
	RatingNomove            int32     `json:"ratingNomove"`
	RatingNmpz              int32     `json:"ratingNmpz"`
	GuessedFirst            float64   `json:"guessedFirst"`
	BestCountries           string    `json:"bestCountries"`
	WorstCountries          string    `json:"worstCountries"`
	RankedTeamMovingGames   int64     `json:"rankedTeamMovingGames"`
	RankedTeamMovingWins    int64     `json:"rankedTeamMovingWins"`
	RankedTeamNomoveGames   int64     `json:"rankedTeamNomoveGames"`
	RankedTeamNomoveWins    int64     `json:"rankedTeamNomoveWins"`
	RankedTeamNmpzGames     int64     `json:"rankedTeamNmpzGames"`
	RankedTeamNmpzWins      int64     `json:"rankedTeamNmpzWins"`
	RankedSoloMovingGames   int64     `json:"rankedSoloMovingGames"`
	RankedSoloMovingWins    int64     `json:"rankedSoloMovingWins"`
	RankedSoloNomoveGames   int64     `json:"rankedSoloNomoveGames"`
	RankedSoloNomoveWins    int64     `json:"rankedSoloNomoveWins"`
	RankedSoloNmpzGames     int64     `json:"rankedSoloNmpzGames"`
	RankedSoloNmpzWins      int64     `json:"rankedSoloNmpzWins"`
	UnrankedSoloMovingGames int64     `json:"unrankedSoloMovingGames"`
	UnrankedSoloMovingWins  int64     `json:"unrankedSoloMovingWins"`
	UnrankedSoloNomoveGames int64     `json:"unrankedSoloNomoveGames"`
	UnrankedSoloNomoveWins  int64     `json:"unrankedSoloNomoveWins"`
	UnrankedSoloNmpzGames   int64     `json:"unrankedSoloNmpzGames"`
	UnrankedSoloNmpzWins    int64     `json:"unrankedSoloNmpzWins"`
}

func GetUserFetchCombinedHistoryDailyRowToUserStatsDailyEntry(row db.GetUserFetchCombinedHistoryDailyRow) UserStatsDailyEntry {
	return UserStatsDailyEntry{
		FetchID:                 row.FetchID,
		FetchTimestamp:          row.FetchTimestamp.Time,
		DivisionName:            row.DivisionName.String,
		RatingOverall:           row.RatingOverall.Int32,
		RatingMoving:            row.RatingMoving.Int32,
		RatingNomove:            row.RatingNomove.Int32,
		RatingNmpz:              row.RatingNmpz.Int32,
		GuessedFirst:            row.GuessedFirst.Float64,
		BestCountries:           row.BestCountries.String,
		WorstCountries:          row.WorstCountries.String,
		RankedTeamMovingGames:   row.RankedTeamMovingGames.Int64,
		RankedTeamMovingWins:    row.RankedTeamMovingWins.Int64,
		RankedTeamNomoveGames:   row.RankedTeamNomoveGames.Int64,
		RankedTeamNomoveWins:    row.RankedTeamNomoveWins.Int64,
		RankedTeamNmpzGames:     row.RankedTeamNmpzGames.Int64,
		RankedTeamNmpzWins:      row.RankedTeamNmpzWins.Int64,
		RankedSoloMovingGames:   row.RankedSoloMovingGames.Int64,
		RankedSoloMovingWins:    row.RankedSoloMovingWins.Int64,
		RankedSoloNomoveGames:   row.RankedSoloNomoveGames.Int64,
		RankedSoloNomoveWins:    row.RankedSoloNomoveWins.Int64,
		RankedSoloNmpzGames:     row.RankedSoloNmpzGames.Int64,
		RankedSoloNmpzWins:      row.RankedSoloNmpzWins.Int64,
		UnrankedSoloMovingGames: row.UnrankedSoloMovingGames.Int64,
		UnrankedSoloMovingWins:  row.UnrankedSoloMovingWins.Int64,
		UnrankedSoloNomoveGames: row.UnrankedSoloNomoveGames.Int64,
		UnrankedSoloNomoveWins:  row.UnrankedSoloNomoveWins.Int64,
		UnrankedSoloNmpzGames:   row.UnrankedSoloNmpzGames.Int64,
		UnrankedSoloNmpzWins:    row.UnrankedSoloNmpzWins.Int64,
	}
}

type UserStatsOutput struct {
	Body struct {
		Entries []UserStatsDailyEntry `json:"entries"`
	}
}

func (a *App) GetUserStats(ctx context.Context, input *UserStatsInput) (*UserStatsOutput, error) {
	a.logger.Debug("handling GetUserStats", "user_id", input.ID)
	resp := &UserStatsOutput{}

	q := db.New(a.db)
	history, err := q.GetUserFetchCombinedHistoryDaily(ctx, input.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, huma.Error404NotFound("user not found")
		}
		a.logger.Warn("failed to get user stats", "err", err)
		return nil, huma.Error500InternalServerError("failed to get user stats")
	}

	resp.Body.Entries = make([]UserStatsDailyEntry, len(history))
	for i, row := range history {
		resp.Body.Entries[i] = GetUserFetchCombinedHistoryDailyRowToUserStatsDailyEntry(row)
	}
	return resp, nil
}
