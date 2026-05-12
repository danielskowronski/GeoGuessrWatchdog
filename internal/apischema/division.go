package apischema

import (
	"encoding/json"
	"fmt"
	"strings"
)

// https://www.geoguessr.com/api/v4/ranked-system/divisions
type DivisionsListApiResponse struct {
	Divisions []DivisionInfoApiResponse `json:"divisions"`
}

type DivisionInfoApiResponse struct {
	// only potentially useful fields
	DivisionNumber int    `json:"divisionNumber"`
	DivisionRank   int    `json:"divisionRank"`
	Tier           string `json:"tier"`
	Name           string `json:"name"`
	// not relying on `gameModes` as this is formatted differently
	Maps map[string]TargetMapInApiResponse `json:"maps"`
}

type TargetMapInApiResponse struct {
	MapID   string `json:"mapId"`
	MapName string `json:"mapName"`
}

type DivisionModeMapInfo struct {
	DivisionName string `json:"divisionName"`
	GameMode     string `json:"gameMode"`
	MapID        string `json:"mapId"`
	MapName      string `json:"mapName"`
}

func (dmmi DivisionModeMapInfo) String() string {
	divisionName := fmt.Sprintf("%-*s", 15, strings.ReplaceAll(dmmi.DivisionName, " ", "_"))
	modeName := fmt.Sprintf("%-*s", 8, strings.ReplaceAll(dmmi.GameMode, "Duels", ""))
	return fmt.Sprintf("DivisionModeMapInfo{DivisionName=%s  GameMode=%s MapID=%s MapName=%q}", divisionName, modeName, dmmi.MapID, dmmi.MapName)
}

func DecodeApiResponseDivisionList(body []byte) ([]DivisionModeMapInfo, error) {
	var apiResp DivisionsListApiResponse
	err := json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal API response: %w", err)
	}

	var result []DivisionModeMapInfo
	for _, division := range apiResp.Divisions {
		for gameMode, targetMap := range division.Maps {
			result = append(result, DivisionModeMapInfo{
				DivisionName: division.Name,
				GameMode:     gameMode,
				MapID:        targetMap.MapID,
				MapName:      targetMap.MapName,
			})
		}
	}

	return result, nil
}
