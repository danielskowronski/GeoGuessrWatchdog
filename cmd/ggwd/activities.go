package main

import (
	"context"
	"time"

	"github.com/danielskowronski/geoguessrwatchdog/internal/config"
)

type ActivityEnum int

const (
	ActivityEnumDivisionFetch ActivityEnum = iota
	ActivityEnumMapFetch
	ActivityEnumNotifyAboutMapChange
	ActivityEnumUserProgressFetch
	ActivityEnumUserStatsFetch
)

type Activities struct {
	Config config.Config
}

// TODO: split fetch and inserts activities to allow DB operations to be retried independently from API calls, and avoid hitting API rate limits on retries of DB operations

func (a *Activities) DivisionFetchActivity(ctx context.Context, input DivisionFetchInput) (DivisionFetchResult, error) {
	result, err := RunDivisionFetch(ctx, a.Config.Database, a.Config.GeoguessrAPI)

	if err != nil {
		return DivisionFetchResult{}, err
	}

	return result, nil
}

func RunDivisionFetch(ctx context.Context, dbConfig config.DatabaseConfig, ggConfig config.GeoguessrAPIConfig) (DivisionFetchResult, error) {
	api, err := NewAPIClient(ggConfig.Proxy, ggConfig.BaseURL, ggConfig.Cookie)
	if err != nil {
		return DivisionFetchResult{}, err
	}

	db := NewDB(dbConfig.URL)

	dmmiList, apiCode, err := api.FetchDivisions(ctx)
	if err != nil {
		return DivisionFetchResult{}, err
	}

	uniqueMaps := NewStringSet()

	changedCount := 0

	for _, dmmi := range dmmiList {
		uniqueMaps.Add(dmmi.MapID)

		changed, err := db.UpsertDivisionModeMapInfo(ctx, dmmi)
		if err != nil {
			return DivisionFetchResult{}, err
		}
		if changed {
			changedCount++
		}
	}

	return DivisionFetchResult{
		ApiResultCode: apiCode,
		UniqueMaps:    uniqueMaps.values,
		Changed:       changedCount,
		FetchedAt:     time.Now().UTC(),
	}, nil
}

func (a *Activities) MapFetchActivity(ctx context.Context, input MapFetchInput) (MapFetchResult, error) {
	result, err := RunMapFetch(ctx, a.Config.Database, a.Config.GeoguessrAPI, input.MapID)

	if err != nil {
		return MapFetchResult{}, err
	}

	return result, nil
}

func RunMapFetch(ctx context.Context, dbConfig config.DatabaseConfig, ggConfig config.GeoguessrAPIConfig, mapID string) (MapFetchResult, error) {
	api, err := NewAPIClient(ggConfig.Proxy, ggConfig.BaseURL, ggConfig.Cookie)
	if err != nil {
		return MapFetchResult{}, err
	}

	db := NewDB(dbConfig.URL)

	mapInfo, apiCode, err := api.FetchMapInfo(ctx, mapID)
	if err != nil {
		return MapFetchResult{}, err
	}

	changed, err := db.UpsertMapInfo(ctx, *mapInfo)
	if err != nil {
		return MapFetchResult{}, err
	}
	changedCount := 0
	if changed {
		changedCount = 1
	}

	return MapFetchResult{
		ApiResultCode: apiCode,
		Changed:       changedCount,
		Locations:     mapInfo.CoordinateCount,
		FetchedAt:     time.Now().UTC(),
	}, nil
}

func (a *Activities) NotifyAboutMapChangeActivity(ctx context.Context, input NotifierInput) (NotifierResult, error) {
	// FIXME: implement this
	return NotifierResult{
		Success: true,
	}, nil
}

func (a *Activities) UserProgressFetchActivity(ctx context.Context, input UserProgressFetchInput) (UserProgressFetchOutput, error) {
	result, err := RunUserProgressFetch(ctx, a.Config.Database, a.Config.GeoguessrAPI, input.UserID)

	if err != nil {
		return UserProgressFetchOutput{}, err
	}

	return result, nil
}
func RunUserProgressFetch(ctx context.Context, dbConfig config.DatabaseConfig, ggConfig config.GeoguessrAPIConfig, userID string) (UserProgressFetchOutput, error) {
	api, err := NewAPIClient(ggConfig.Proxy, ggConfig.BaseURL, ggConfig.Cookie)
	if err != nil {
		return UserProgressFetchOutput{}, err
	}

	db := NewDB(dbConfig.URL)

	userProgress, apiCode, err := api.FetchUserProgress(ctx, userID)
	if err != nil {
		return UserProgressFetchOutput{}, err
	}

	_, err = db.InsertUserProgressHistory(ctx, *userProgress, userID)
	if err != nil {
		return UserProgressFetchOutput{}, err
	}

	return UserProgressFetchOutput{
		ApiResultCode: apiCode,
		RatingOverall: uint64(userProgress.RatingOverall),
		FetchedAt:     time.Now().UTC(),
	}, nil
}

func (a *Activities) UserStatsFetchActivity(ctx context.Context, input UserStatsFetchInput) (UserStatsFetchOutput, error) {
	result, err := RunUserStatsFetch(ctx, a.Config.Database, a.Config.GeoguessrAPI, input.UserID)
	if err != nil {
		return UserStatsFetchOutput{}, err
	}

	return result, nil
}

func RunUserStatsFetch(ctx context.Context, dbConfig config.DatabaseConfig, ggConfig config.GeoguessrAPIConfig, userID string) (UserStatsFetchOutput, error) {
	api, err := NewAPIClient(ggConfig.Proxy, ggConfig.BaseURL, ggConfig.Cookie)
	if err != nil {
		return UserStatsFetchOutput{}, err
	}

	db := NewDB(dbConfig.URL)

	userStats, apiCode, err := api.FetchUserStats(ctx, userID)
	if err != nil {
		return UserStatsFetchOutput{}, err
	}

	_, err = db.InsertUserStatsHistory(ctx, *userStats, userID)
	if err != nil {
		return UserStatsFetchOutput{}, err
	}

	gamesCount := userStats.RankedTeamMovingGames + userStats.RankedTeamNoMoveGames + userStats.RankedTeamNMPZGames +
		userStats.RankedSoloMovingGames + userStats.RankedSoloNoMoveGames + userStats.RankedSoloNMPZGames +
		userStats.UnrankedSoloMovingGames + userStats.UnrankedSoloNoMoveGames + userStats.UnrankedSoloNMPZGames

	return UserStatsFetchOutput{
		ApiResultCode: apiCode,
		GamesCount:    uint64(gamesCount),
		FetchedAt:     time.Now().UTC(),
	}, nil
}
