package apischema

import (
	"encoding/json"
	"fmt"
)

// https://www.geoguessr.com/api/v4/ranked-system/progress/$USER_ID
type ProgressApiResponse struct {
	DivisionNumber   uint8            `json:"divisionNumber"`
	DivisionName     string           `json:"divisionName"`
	Rating           uint16           `json:"rating"`
	Tier             string           `json:"tier"`
	GameModeRatings  *GameModeRatings `json:"gameModeRatings"`
	GuessedFirstRate float64          `json:"guessedFirstRate"`
	WinStreak        uint16           `json:"winStreak"`
	LatestGamesWon   []bool           `json:"latestGames"`
	BestCountries    []string         `json:"bestCountries"`
	WorstCountries   []string         `json:"worstCountries"`
}

type GameModeRatings struct {
	StandardDuels *uint16 `json:"standardDuels,omitempty"`
	NoMoveDuels   *uint16 `json:"noMoveDuels,omitempty"`
	NmpzDuels     *uint16 `json:"nmpzDuels,omitempty"`
}

type ProgressInfo struct {
	DivisionName   string `json:"divisionName"`
	DivisionNumber uint8  `json:"divisionNumber"`

	RatingOverall uint16 `json:"ratingOverall"`
	RatingMoving  uint16 `json:"ratingMoving"`
	RatingNoMove  uint16 `json:"ratingNoMove"`
	RatingNMPZ    uint16 `json:"ratingNMPZ"`

	GuessedFirstRate float64 `json:"guessedFirstRate"`
	BestCountries    string  `json:"bestCountries"`
	WorstCountries   string  `json:"worstCountries"`
}

func (pi ProgressInfo) String() string {
	return fmt.Sprintf("ProgressInfo{DivisionName=%q DivisionNumber=%d RatingOverall=%d RatingMoving=%d RatingNoMove=%d RatingNMPZ=%d GuessedFirstRate=%.2f BestCountries=%q WorstCountries=%q}",
		pi.DivisionName, pi.DivisionNumber, pi.RatingOverall, pi.RatingMoving, pi.RatingNoMove, pi.RatingNMPZ, pi.GuessedFirstRate, pi.BestCountries, pi.WorstCountries)
}

func (gmr *GameModeRatings) decodeGameModeRating() (uint16, uint16, uint16) {
	if gmr == nil {
		return 0, 0, 0
	}
	movingRating := uint16(0)
	noMoveRating := uint16(0)
	nmpzRating := uint16(0)

	if gmr.StandardDuels != nil {
		movingRating = *gmr.StandardDuels
	}
	if gmr.NoMoveDuels != nil {
		noMoveRating = *gmr.NoMoveDuels
	}
	if gmr.NmpzDuels != nil {
		nmpzRating = *gmr.NmpzDuels
	}

	return movingRating, noMoveRating, nmpzRating
}

func DecodeApiResponseProgress(body []byte) (*ProgressInfo, error) {
	var apiResp ProgressApiResponse
	err := json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal API response: %w", err)
	}

	movingRating, noMoveRating, nmpzRating := apiResp.GameModeRatings.decodeGameModeRating()

	return &ProgressInfo{
		DivisionName:     apiResp.DivisionName,
		DivisionNumber:   apiResp.DivisionNumber,
		RatingOverall:    apiResp.Rating,
		RatingMoving:     movingRating,
		RatingNoMove:     noMoveRating,
		RatingNMPZ:       nmpzRating,
		GuessedFirstRate: apiResp.GuessedFirstRate,
		BestCountries:    fmt.Sprintf("%v", apiResp.BestCountries),
		WorstCountries:   fmt.Sprintf("%v", apiResp.WorstCountries),
	}, nil
}
