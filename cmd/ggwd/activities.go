package main

import (
	"context"
	"time"
)

type Activities struct {
	DatabaseURL      string
	HttpProxyURL     string
	GeoGuessrApiBase string
	GeoGuessrCookie  string
}

// TODO: split fetch and inserts activities to allow DB operations to be retried independently from API calls, and avoid hitting API rate limits on retries of DB operations

func (a *Activities) DivisionFetchActivity(ctx context.Context, input DivisionFetchInput) (DivisionFetchResult, error) {
	result, err := RunDivisionFetch(ctx, a.DatabaseURL, a.HttpProxyURL, a.GeoGuessrApiBase, a.GeoGuessrCookie)

	if err != nil {
		return DivisionFetchResult{}, err
	}

	return result, nil
}

func RunDivisionFetch(ctx context.Context, dbURL string, httpProxyURL string, apiBase string, cookie string) (DivisionFetchResult, error) {
	api, err := NewAPIClient(httpProxyURL, apiBase, cookie)
	if err != nil {
		return DivisionFetchResult{}, err
	}

	db := NewDB(dbURL)

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
	result, err := RunMapFetch(ctx, a.DatabaseURL, a.HttpProxyURL, a.GeoGuessrApiBase, a.GeoGuessrCookie, input.MapID)

	if err != nil {
		return MapFetchResult{}, err
	}

	return result, nil
}

func RunMapFetch(ctx context.Context, dbURL string, httpProxyURL string, apiBase string, cookie string, mapID string) (MapFetchResult, error) {
	api, err := NewAPIClient(httpProxyURL, apiBase, cookie)
	if err != nil {
		return MapFetchResult{}, err
	}

	db := NewDB(dbURL)

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

func (a *Activities) NotifyActivity(ctx context.Context, input NotifierInput) (NotifierResult, error) {
	// FIXME: implement this
	return NotifierResult{
		Success: true,
	}, nil
}
