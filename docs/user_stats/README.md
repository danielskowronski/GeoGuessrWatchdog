# User stats - competitive games

## Game counts

URL is `https://www.geoguessr.com/api/v4/stats/users/USER_ID` and example response is at [counts.json](./counts.json)

API responds with various keys, most of the descrive specific competitive game mode:

- `rankedTeamDuelsStandard`
- `rankedTeamDuelsNoMove`
- `rankedTeamDuelsNmpz`
- `rankedTeamDuelsTotal`
- `battleRoyaleDistance`
- `battleRoyaleCountry`
- `competitiveCityStreaks`
- `duels` (moving)
- `duelsNoMove`
- `duelsNmpz`
- `duelsTotal`
- `unrankedDuels` (moving)
- `unrankedDuelsNoMove`
- `unrankedDuelsNmpz`
- `teamDuels`
- `teamDuelsQuickplay`

Other keys:

- `competitiveStreaksMedals`
- `battleRoyaleMedals`
- `duelsMedals`
- `totalMedals`
- `lifeTimeXpProgression`
- `quickplayFlawlessVictories`
- `perfectRounds`
- `party`

All competitive game modes will have at least following fields:

- `numGamesPlayed`
- `numWins`
- `winRatio` (float, 0-1)

Some may contain additionaly:

- `avgPosition`
- `avgGuessDistance`
- `numGuesses`
- `numFlawlessWins`

The goal of watching this API is to track daily number of games. To make it faily universal, all known competitive game modes must be tracked. It also makes sense to treat one API fetch as one target row in DB. From this perspective, we'll need 16x3 columns to match each game mode to 3 fields.

## Game stats

URL is `https://www.geoguessr.com/api/v4/ranked-system/progress/USER_ID` and example response is at [stats.json](./stats.json)

Fields are as follows:

- `divisionNumber` - matches numerical ID from [competitive_maps](../competitive_maps/); can be null
- `divisionName` (string); can be null
- `tier` (string); can be null
- `rating` - overall ELO (int) or null
- `gameModeRatings` (can be missing) - can contain following:
  - `standardDuels` - ELO (int)
  - `noMoveDuels` - ELO (int)
  - `nmpzDuels` - ELO (int)
- `guessedFirstRate` (float), defaults to `0`
- `winStreak` (int) - count of consecutive win games, `0` if last was lost
- `latestGames` - list of booleans indicating if last games was won (first is latest game)
- `bestCountries` and `worstCountries` - list of (three) strings containing lowercase ISO 3166-1 A-2 (2-letter) contry codes; can be empty

`rating` is overall ELO, and is not supposed to be average or max of individual game mode ELO scores. There are variety of reasons for that, most common is ELO roll-back after opponent was banned for cheating.

Following fields matter:

- `divisionName`; it's string, but `divisionNumber` is not stable
- `rating` - defaulting to 0
- `gameModeRatings.standardDuels` - defaulting to 0 if does not exist
- `gameModeRatings.noMoveDuels` - defaulting to 0 if does not exist
- `gameModeRatings.nmpzDuels` - defaulting to 0 if does not exist
- `guessedFirstRate`
- `bestCountries` and `worstCountries` which can be flattened to CSV; this is interesting field for further research
