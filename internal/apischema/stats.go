package apischema

import (
	"encoding/json"
	"fmt"
)

// https://www.geoguessr.com/api/v4/stats/users/$USER_ID
type StatsApiResponse struct {
	RankedTeamDuelsStandard StatsModeRespose `json:"rankedTeamDuelsStandard"`
	RankedTeamDuelsNoMove   StatsModeRespose `json:"rankedTeamDuelsNoMove"`
	RankedTeamDuelsNMPZ     StatsModeRespose `json:"rankedTeamDuelsNmpz"`
	RankedTeamDulesTotal    StatsModeRespose `json:"rankedTeamDuelsTotal"`

	BattleRoyaleDistance StatsModeRespose   `json:"battleRoyaleDistance"`
	BattleRoyaleCountry  StatsModeRespose   `json:"battleRoyaleCountry"`
	BattleRoyaleMedals   MedalsModeResponse `json:"battleRoyaleMedals"`

	CompetitiveCityStreaks   StatsModeRespose   `json:"competitiveCityStreaks"`
	CompetitiveStreaksMedals MedalsModeResponse `json:"competitiveStreaksMedals"`

	Duels       StatsModeRespose   `json:"duels"`
	DuelsNoMove StatsModeRespose   `json:"duelsNoMove"`
	DuelsNMPZ   StatsModeRespose   `json:"duelsNmpz"`
	DuelsTotal  StatsModeRespose   `json:"duelsTotal"`
	DuelsMedals MedalsModeResponse `json:"duelsMedals"`

	UnrankedDuels       StatsModeRespose `json:"unrankedDuels"`
	UnrankedDuelsNoMove StatsModeRespose `json:"unrankedDuelsNoMove"`
	UnrankedDuelsNMPZ   StatsModeRespose `json:"unrankedDuelsNmpz"`
	UnrankedDuelsTotal  StatsModeRespose `json:"unrankedDuelsTotal"`

	LifetimeProgression LifetimeProgressionResponse `json:"lifetimeProgression"`
	TotalMedals         MedalsModeResponse          `json:"totalMedals"`

	TeamDuels StatsModeRespose `json:"teamDuels"`

	TeamDuelsQuickplay         StatsModeRespose `json:"teamDuelsQuickplay"`
	QuickplayFlawlessVictories uint64           `json:"quickplayFlawlessVictories"`
	PerfectRounds              uint64           `json:"perfectRounds"`

	Party PartyResponse `json:"party"`
}

type StatsModeRespose struct {
	NumGamesPlayed uint64  `json:"numGamesPlayed,omitempty"`
	NumWins        uint64  `json:"numWins,omitempty"`
	WinRatio       float64 `json:"winRatio,omitempty"`

	AveragePosition        float64 `json:"avgPosition,omitempty"`
	AveraregeGuessDistance float64 `json:"avgGuessDistance,omitempty"`
	NumGuesses             uint64  `json:"numGuesses,omitempty"`

	AverageCorrectGuesses float64 `json:"avgCorrectGuesses,omitempty"`

	NumFlawlessWins uint64 `json:"numFlawlessWins,omitempty"`
}

type MedalsModeResponse struct {
	MedalCountGold   uint64 `json:"medalCountGold,omitempty"`
	MedalCountSilver uint64 `json:"medalCountSilver,omitempty"`
	MedalCountBronze uint64 `json:"medalCountBronze,omitempty"`
}

type LifetimeProgressionResponse struct {
	XP           uint64                           `json:"xp,omitempty"`
	CurrentLevel LifetimeProgressionLevelResponse `json:"currentLevel,omitempty"`
	NextLevel    LifetimeProgressionLevelResponse `json:"nextLevel,omitempty"`
	CurrentTitle LifetimeProgressionTitleResponse `json:"currentTitle,omitempty"`
}

type LifetimeProgressionLevelResponse struct {
	Level   uint64  `json:"level,omitempty"`
	XpStart float64 `json:"xpStart,omitempty"`
}

type LifetimeProgressionTitleResponse struct {
	ID           uint64 `json:"id,omitempty"`
	TierID       uint64 `json:"tierId,omitempty"`
	MinimumLevel uint64 `json:"minimumLevel,omitempty"`
	Name         string `json:"name,omitempty"`
}

type PartyResponse struct {
	Total                 uint64 `json:"total,omitempty"`
	Duels                 uint64 `json:"duels,omitempty"`
	TeamDuels             uint64 `json:"teamDuels,omitempty"`
	BattleRoyaleCountries uint64 `json:"battleRoyaleCountries,omitempty"`
	BattleRoyaleDistance  uint64 `json:"battleRoyaleDistance,omitempty"`
	CityStreaks           uint64 `json:"cityStreaks,omitempty"`
	LiveChallenges        uint64 `json:"liveChallenges,omitempty"`
	Bullseye              uint64 `json:"bullseye,omitempty"`
	Quizzes               uint64 `json:"quizzes,omitempty"`
}

type StatsInfo struct {
	RankedTeamMovingGames uint64  `json:"rankedTeamMovingGames"`
	RankedTeamMovingWins  uint64  `json:"rankedTeamMovingWins"`
	RankedTeamMovingRatio float64 `json:"rankedTeamMovingRatio"`

	RankedTeamNoMoveGames uint64  `json:"rankedTeamNoMoveGames"`
	RankedTeamNoMoveWins  uint64  `json:"rankedTeamNoMoveWins"`
	RankedTeamNoMoveRatio float64 `json:"rankedTeamNoMoveRatio"`

	RankedTeamNMPZGames uint64  `json:"rankedTeamNMPZGames"`
	RankedTeamNMPZWins  uint64  `json:"rankedTeamNMPZWins"`
	RankedTeamNMPZRatio float64 `json:"rankedTeamNMPZRatio"`

	RankedSoloMovingGames uint64  `json:"rankedSoloMovingGames"`
	RankedSoloMovingWins  uint64  `json:"rankedSoloMovingWins"`
	RankedSoloMovingRatio float64 `json:"rankedSoloMovingRatio"`

	RankedSoloNoMoveGames uint64  `json:"rankedSoloNoMoveGames"`
	RankedSoloNoMoveWins  uint64  `json:"rankedSoloNoMoveWins"`
	RankedSoloNoMoveRatio float64 `json:"rankedSoloNoMoveRatio"`

	RankedSoloNMPZGames uint64  `json:"rankedSoloNMPZGames"`
	RankedSoloNMPZWins  uint64  `json:"rankedSoloNMPZWins"`
	RankedSoloNMPZRatio float64 `json:"rankedSoloNMPZRatio"`

	UnrankedSoloMovingGames uint64  `json:"unrankedSoloMovingGames"`
	UnrankedSoloMovingWins  uint64  `json:"unrankedSoloMovingWins"`
	UnrankedSoloMovingRatio float64 `json:"unrankedSoloMovingRatio"`

	UnrankedSoloNoMoveGames uint64  `json:"unrankedSoloNoMoveGames"`
	UnrankedSoloNoMoveWins  uint64  `json:"unrankedSoloNoMoveWins"`
	UnrankedSoloNoMoveRatio float64 `json:"unrankedSoloNoMoveRatio"`

	UnrankedSoloNMPZGames uint64  `json:"unrankedSoloNMPZGames"`
	UnrankedSoloNMPZWins  uint64  `json:"unrankedSoloNMPZWins"`
	UnrankedSoloNMPZRatio float64 `json:"unrankedSoloNMPZRatio"`
}

func (si StatsInfo) String() string {
	return fmt.Sprintf("StatsInfo{RankedTeamMovingGames=%d RankedTeamNoMoveGames=%d RankedTeamNMPZGames=%d RankedSoloMovingGames=%d RankedSoloNoMoveGames=%d RankedSoloNMPZGames=%d UnrankedSoloMovingGames=%d UnrankedSoloNoMoveGames=%d UnrankedSoloNMPZGames=%d}",
		si.RankedTeamMovingGames, si.RankedTeamNoMoveGames, si.RankedTeamNMPZGames, si.RankedSoloMovingGames, si.RankedSoloNoMoveGames, si.RankedSoloNMPZGames, si.UnrankedSoloMovingGames, si.UnrankedSoloNoMoveGames, si.UnrankedSoloNMPZGames)
}

func DecodeApiResponseStats(body []byte) (*StatsInfo, error) {
	var apiResp StatsApiResponse
	err := json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal API response: %w", err)
	}

	return &StatsInfo{
		RankedTeamMovingGames:   apiResp.RankedTeamDuelsStandard.NumGamesPlayed,
		RankedTeamMovingWins:    apiResp.RankedTeamDuelsStandard.NumWins,
		RankedTeamMovingRatio:   apiResp.RankedTeamDuelsStandard.WinRatio,
		RankedTeamNoMoveGames:   apiResp.RankedTeamDuelsNoMove.NumGamesPlayed,
		RankedTeamNoMoveWins:    apiResp.RankedTeamDuelsNoMove.NumWins,
		RankedTeamNoMoveRatio:   apiResp.RankedTeamDuelsNoMove.WinRatio,
		RankedTeamNMPZGames:     apiResp.RankedTeamDuelsNMPZ.NumGamesPlayed,
		RankedTeamNMPZWins:      apiResp.RankedTeamDuelsNMPZ.NumWins,
		RankedTeamNMPZRatio:     apiResp.RankedTeamDuelsNMPZ.WinRatio,
		RankedSoloMovingGames:   apiResp.Duels.NumGamesPlayed,
		RankedSoloMovingWins:    apiResp.Duels.NumWins,
		RankedSoloMovingRatio:   apiResp.Duels.WinRatio,
		RankedSoloNoMoveGames:   apiResp.DuelsNoMove.NumGamesPlayed,
		RankedSoloNoMoveWins:    apiResp.DuelsNoMove.NumWins,
		RankedSoloNoMoveRatio:   apiResp.DuelsNoMove.WinRatio,
		RankedSoloNMPZGames:     apiResp.DuelsNMPZ.NumGamesPlayed,
		RankedSoloNMPZWins:      apiResp.DuelsNMPZ.NumWins,
		RankedSoloNMPZRatio:     apiResp.DuelsNMPZ.WinRatio,
		UnrankedSoloMovingGames: apiResp.UnrankedDuels.NumGamesPlayed,
		UnrankedSoloMovingWins:  apiResp.UnrankedDuels.NumWins,
		UnrankedSoloMovingRatio: apiResp.UnrankedDuels.WinRatio,
		UnrankedSoloNoMoveGames: apiResp.UnrankedDuelsNoMove.NumGamesPlayed,
		UnrankedSoloNoMoveWins:  apiResp.UnrankedDuelsNoMove.NumWins,
		UnrankedSoloNoMoveRatio: apiResp.UnrankedDuelsNoMove.WinRatio,
		UnrankedSoloNMPZGames:   apiResp.UnrankedDuelsNMPZ.NumGamesPlayed,
		UnrankedSoloNMPZWins:    apiResp.UnrankedDuelsNMPZ.NumWins,
		UnrankedSoloNMPZRatio:   apiResp.UnrankedDuelsNMPZ.WinRatio,
	}, nil
}
