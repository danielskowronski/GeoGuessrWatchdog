-- name: GetUserFetchCombinedHistoryDaily :many
SELECT DISTINCT ON (date_trunc('day', fetch_timestamp))
    fetch_id,
    fetch_timestamp,

    division_name,
    division_number,

    rating_overall,
    rating_moving,
    rating_nomove,
    rating_nmpz,
    guessed_first,
    best_countries,
    worst_countries,

    ranked_team_moving_games,
    ranked_team_moving_wins,
    ranked_team_nomove_games,
    ranked_team_nomove_wins,
    ranked_team_nmpz_games,
    ranked_team_nmpz_wins,

    ranked_solo_moving_games,
    ranked_solo_moving_wins,
    ranked_solo_nomove_games,
    ranked_solo_nomove_wins,
    ranked_solo_nmpz_games,
    ranked_solo_nmpz_wins,

    unranked_solo_moving_games,
    unranked_solo_moving_wins,
    unranked_solo_nomove_games,
    unranked_solo_nomove_wins,
    unranked_solo_nmpz_games,
    unranked_solo_nmpz_wins,

    singleplayer_games_played,
    singleplayer_rounds_played
FROM user_fetch_combined_history
WHERE user_id = $1
ORDER BY date_trunc('day', fetch_timestamp), fetch_timestamp ASC;