package main

import (
	"time"
)

type WorkflowInput struct {
	TriggerApiUpdates    bool `json:"triggerApiUpdates,omitempty"`
	TriggerMapUpdates    bool `json:"triggerMapUpdates,omitempty"`
	MaxParallel          int  `json:"maxParallel,omitempty"`
	TriggerNotifications bool `json:"triggerNotifications,omitempty"`
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
