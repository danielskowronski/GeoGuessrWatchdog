package apischema

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	MAP_FRONTEND_URL = "https://www.geoguessr.com/maps/%s"
)

// https://www.geoguessr.com/api/maps/$MAP_ID
type MapApiResponse struct {
	ID               string               `json:"id"`
	Name             string               `json:"name"`
	Description      string               `json:"description"`
	Bounds           BoundsInApiResponse  `json:"bounds"`
	UpdatedAt        string               `json:"updatedAt"`
	MaxErrorDistance int64                `json:"maxErrorDistance"`
	MapSize          MapSizeInApiResponse `json:"mapSize"`
}

type BoundsInApiResponse struct {
	Max BoundsPointInApiResponse `json:"max"`
	Min BoundsPointInApiResponse `json:"min"`
}

type BoundsPointInApiResponse struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type MapSizeInApiResponse struct {
	CoordinateCount int64 `json:"coordinateCount"`
}

type MapInfo struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	BoundsMinLat     float64   `json:"boundsMinLat"`
	BoundsMinLng     float64   `json:"boundsMinLng"`
	BoundsMaxLat     float64   `json:"boundsMaxLat"`
	BoundsMaxLng     float64   `json:"boundsMaxLng"`
	UpdatedAt        time.Time `json:"updatedAt"`
	MaxErrorDistance int64     `json:"maxErrorDistance"`
	CoordinateCount  int64     `json:"coordinateCount"`
}

func (mi MapInfo) String() string {
	return fmt.Sprintf("MapInfo{ID=%s Name=%q Description=%q Bounds=[%f,%f - %f,%f] UpdatedAt=%s MaxErrorDistance=%d CoordinateCount=%d}",
		mi.ID, mi.Name, mi.Description, mi.BoundsMinLat, mi.BoundsMinLng, mi.BoundsMaxLat, mi.BoundsMaxLng, mi.UpdatedAt.Format(time.RFC3339), mi.MaxErrorDistance, mi.CoordinateCount)
}

func coordinateToString(value float64, neg string, pos string) string {
	if value < 0 {
		return fmt.Sprintf("%10.6f %s", -value, neg)
	} else {
		return fmt.Sprintf("%10.6f %s", value, pos)
	}
}

func (mi MapInfo) Coordinates() string {
	return fmt.Sprintf("[%s %s - %s %s]",
		coordinateToString(mi.BoundsMinLat, "S", "N"),
		coordinateToString(mi.BoundsMinLng, "W", "E"),
		coordinateToString(mi.BoundsMaxLat, "S", "N"),
		coordinateToString(mi.BoundsMaxLng, "W", "E"))
}

func DecodeApiResponseMap(body []byte) (*MapInfo, error) {
	var apiResp MapApiResponse
	err := json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal API response: %w", err)
	}

	updatedAt, err := time.Parse(time.RFC3339, apiResp.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse updatedAt: %w", err)
	}

	return &MapInfo{
		ID:               apiResp.ID,
		Name:             apiResp.Name,
		Description:      apiResp.Description,
		BoundsMinLat:     apiResp.Bounds.Min.Lat,
		BoundsMinLng:     apiResp.Bounds.Min.Lng,
		BoundsMaxLat:     apiResp.Bounds.Max.Lat,
		BoundsMaxLng:     apiResp.Bounds.Max.Lng,
		UpdatedAt:        updatedAt,
		MaxErrorDistance: apiResp.MaxErrorDistance,
		CoordinateCount:  apiResp.MapSize.CoordinateCount,
	}, nil
}
