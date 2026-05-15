package main

import (
	"time"

	"github.com/danielskowronski/geoguessrwatchdog/internal/config"
)

type FetchDivisionsMapsWorkflowInput struct {
	TriggerApiUpdates    bool `json:"triggerApiUpdates,omitempty"`
	TriggerMapUpdates    bool `json:"triggerMapUpdates,omitempty"`
	TriggerNotifications bool `json:"triggerNotifications,omitempty"`

	TemporalOptions config.TemporalAdvanced `json:"temporalOptions,omitempty"`
}

type DivisionFetchInput struct {
}

type DivisionFetchResult struct {
	ApiResultCode int       `json:"api_result_code"`
	Changed       int       `json:"changed"`
	UniqueMaps    []string  `json:"unique_maps"`
	FetchedAt     time.Time `json:"fetched_at"`
}

type MapFetchInput struct {
	MapID string `json:"map_id"`
}

type MapFetchResult struct {
	ApiResultCode int       `json:"api_result_code"`
	Changed       int       `json:"changed"`
	Locations     int64     `json:"locations"`
	FetchedAt     time.Time `json:"fetched_at"`
}

type NotifierInput struct {
}
type NotifierResult struct {
	Success bool `json:"success"`
}

type FetchMuiltipleUsersStatsAndProgressWorkflowInput struct {
	TriggerUserStats    bool `json:"triggerUserStats,omitempty"`
	TriggerUserProgress bool `json:"triggerUserProgress,omitempty"`

	UsersList []string `json:"usersList,omitempty"`

	TemporalOptions config.TemporalAdvanced `json:"temporalOptions,omitempty"`
}
type FetchSingleUserStatsAndProgressWorkflowInput struct {
	TriggerUserStats    bool `json:"triggerUserStats,omitempty"`
	TriggerUserProgress bool `json:"triggerUserProgress,omitempty"`

	UserID string `json:"userID,omitempty"`

	TemporalOptions config.TemporalAdvanced `json:"temporalOptions,omitempty"`
}

type UserFetchInitInput struct {
	UserID string `json:"user_id"`
}

type UserFetchInitOutput struct {
	FetchID int64 `json:"fetch_id"`
}

type UserProgressFetchInput struct {
	UserID  string `json:"user_id"`
	FetchID int64  `json:"fetch_id"`
}

type UserProgressFetchOutput struct {
	ApiResultCode int       `json:"api_result_code"`
	RatingOverall uint64    `json:"rating_overall"`
	FetchedAt     time.Time `json:"fetched_at"`
}

type UserStatsFetchInput struct {
	UserID  string `json:"user_id"`
	FetchID int64  `json:"fetch_id"`
}

type UserStatsFetchOutput struct {
	ApiResultCode int       `json:"api_result_code"`
	GamesCount    uint64    `json:"games_count"`
	FetchedAt     time.Time `json:"fetched_at"`
}
