package apischema

import (
	"encoding/json"
	"fmt"
)

// https://www.geoguessr.com/api/v3/users/$USER_ID/stats
type SingleplayerStatsApiResponse struct {
	GamesPlayed  uint64 `json:"gamesPlayed,omitempty"`
	RoundsPlayed uint64 `json:"roundsPlayed,omitempty"`
}
type SingleplayerStatsInfo SingleplayerStatsApiResponse

func (si SingleplayerStatsInfo) String() string {
	return fmt.Sprintf("SingleplayerStatsInfo{GamesPlayed=%d RoundsPlayed=%d}",
		si.GamesPlayed, si.RoundsPlayed)
}
func DecodeApiResponseSingleplayerStats(body []byte) (*SingleplayerStatsInfo, error) {
	var apiResp SingleplayerStatsApiResponse
	err := json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal API response: %w", err)
	}

	return &SingleplayerStatsInfo{
		GamesPlayed:  apiResp.GamesPlayed,
		RoundsPlayed: apiResp.RoundsPlayed,
	}, nil
}
