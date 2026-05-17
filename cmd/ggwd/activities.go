package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
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
	for _, dm := range a.Config.Watchdogs.CompetitiveMaps.NotifyAbout {
		err := RunNotifyAboutMapChange(ctx, a.Config.Database, a.Config.NotifierAPI.Discord, dm)
		if err != nil {
			return NotifierResult{
				Success: false,
			}, err
		}
	}

	return NotifierResult{
		Success: true,
	}, nil
}

type DivisionMapDelta struct {
	MapPointerChanged bool
	MapDetailsChanged bool
	Details           string
}

func (dmmd *DivisionModeMapDetails) CompareWithNotification(lastNotification *DivisionModeMapDetails) DivisionMapDelta {
	// TODO: cover it woth feature flags
	delta := DivisionMapDelta{}
	delta.Details = "Maps didn't change"

	if lastNotification.mapInfo.ID != dmmd.mapInfo.ID {
		delta.MapPointerChanged = true
		delta.Details = "New map assigned"
	} else {
		changes := []string{}
		if lastNotification.mapInfo.BoundsMinLat != dmmd.mapInfo.BoundsMinLat ||
			lastNotification.mapInfo.BoundsMinLng != dmmd.mapInfo.BoundsMinLng ||
			lastNotification.mapInfo.BoundsMaxLat != dmmd.mapInfo.BoundsMaxLat ||
			lastNotification.mapInfo.BoundsMaxLng != dmmd.mapInfo.BoundsMaxLng {
			delta.MapDetailsChanged = true
			msg := "boundaries changed: " + lastNotification.mapInfo.Coordinates() + " -> " + dmmd.mapInfo.Coordinates()
			changes = append(changes, msg)
		}
		if lastNotification.mapInfo.MaxErrorDistance != dmmd.mapInfo.MaxErrorDistance {
			delta.MapDetailsChanged = true
			msg := fmt.Sprintf("max error distance changed: %d -> %d", lastNotification.mapInfo.MaxErrorDistance, dmmd.mapInfo.MaxErrorDistance)
			changes = append(changes, msg)
		}
		if lastNotification.mapInfo.CoordinateCount != dmmd.mapInfo.CoordinateCount {
			delta.MapDetailsChanged = true
			msg := fmt.Sprintf("coordinate count changed: %d -> %d", lastNotification.mapInfo.CoordinateCount, dmmd.mapInfo.CoordinateCount)
			changes = append(changes, msg)
		}
		if !lastNotification.mapInfo.UpdatedAt.Equal(dmmd.mapInfo.UpdatedAt) {
			delta.MapDetailsChanged = true
			msg := fmt.Sprintf("updatedAt changed: %s -> %s", lastNotification.mapInfo.UpdatedAt.Format(time.RFC1123), dmmd.mapInfo.UpdatedAt.Format(time.RFC1123))
			changes = append(changes, msg)
		}
		if len(changes) > 0 {
			delta.Details = "Map details changed:\n- " + strings.Join(changes, "\n- ")
		}
	}

	return delta
}

func RunNotifyAboutMapChange(ctx context.Context, dbConfig config.DatabaseConfig, discordConfig config.DiscordConfig, dm config.NotifyAboutDivisionModeConfig) error {
	db := NewDB(dbConfig.URL)

	noPreviousNotifications := false
	var dmd DivisionMapDelta

	currentDivisionMap, err := db.GetCurrentDivisionMapInfo(ctx, dm.DivisionName, dm.GameMode)
	if err != nil {
		return err
	}

	lastNotification, err := db.GetLastDivisionMapNotification(ctx, dm.DivisionName, dm.GameMode)
	if err != nil {
		if errors.Is(err, ErrNoPreviousNotification) {
			noPreviousNotifications = true
		} else {
			return err
		}
	} else {
		dmd = currentDivisionMap.CompareWithNotification(lastNotification)
	}

	shouldNotify := noPreviousNotifications || dmd.MapPointerChanged || dmd.MapDetailsChanged

	if shouldNotify {
		var dm DivisionModeMapDetails
		dm.divisionName = currentDivisionMap.divisionName
		dm.gameMode = currentDivisionMap.gameMode
		dm.mapInfo = currentDivisionMap.mapInfo
		err = SendDiscordMapChangeNotification(ctx, discordConfig, dm, dmd)
		if err != nil {
			return fmt.Errorf("send discord notification: %w", err)
		}

		err = db.UpsertDivisionMapNotification(ctx, *currentDivisionMap)
		if err != nil {
			return fmt.Errorf("insert division map notification to DB: %w", err)
		}
	}

	return nil
}

func (a *Activities) InsertUserFetchHistoryActivity(ctx context.Context, userID string) (int64, error) {
	result, err := RunInsertUserFetchHistory(ctx, a.Config.Database, userID)
	if err != nil {
		return 0, err
	}

	return result, nil
}
func RunInsertUserFetchHistory(ctx context.Context, dbConfig config.DatabaseConfig, userID string) (int64, error) {
	db := NewDB(dbConfig.URL)

	fetchID, err := db.InsertUserFetchHistory(ctx, userID)
	if err != nil {
		return 0, err
	}

	return fetchID, nil
}

func (a *Activities) UserProgressFetchActivity(ctx context.Context, input UserProgressFetchInput) (UserProgressFetchOutput, error) {
	result, err := RunUserProgressFetch(ctx, a.Config.Database, a.Config.GeoguessrAPI, input.UserID, input.FetchID)

	if err != nil {
		return UserProgressFetchOutput{}, err
	}

	return result, nil
}
func RunUserProgressFetch(ctx context.Context, dbConfig config.DatabaseConfig, ggConfig config.GeoguessrAPIConfig, userID string, fetchID int64) (UserProgressFetchOutput, error) {
	api, err := NewAPIClient(ggConfig.Proxy, ggConfig.BaseURL, ggConfig.Cookie)
	if err != nil {
		return UserProgressFetchOutput{}, err
	}

	db := NewDB(dbConfig.URL)

	userProgress, apiCode, err := api.FetchUserProgress(ctx, userID)
	if err != nil {
		return UserProgressFetchOutput{}, err
	}

	_, err = db.InsertUserProgressHistory(ctx, *userProgress, userID, fetchID)
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
	result, err := RunUserStatsFetch(ctx, a.Config.Database, a.Config.GeoguessrAPI, input.UserID, input.FetchID)
	if err != nil {
		return UserStatsFetchOutput{}, err
	}

	return result, nil
}

func RunUserStatsFetch(ctx context.Context, dbConfig config.DatabaseConfig, ggConfig config.GeoguessrAPIConfig, userID string, fetchID int64) (UserStatsFetchOutput, error) {
	api, err := NewAPIClient(ggConfig.Proxy, ggConfig.BaseURL, ggConfig.Cookie)
	if err != nil {
		return UserStatsFetchOutput{}, err
	}

	db := NewDB(dbConfig.URL)

	userStats, apiCode, err := api.FetchUserStats(ctx, userID)
	if err != nil {
		return UserStatsFetchOutput{}, err
	}

	_, err = db.InsertUserStatsHistory(ctx, *userStats, userID, fetchID)
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
