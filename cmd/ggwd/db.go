package main

import (
	"context"
	"fmt"

	"github.com/danielskowronski/geoguessrwatchdog/internal/apischema"
	db "github.com/danielskowronski/geoguessrwatchdog/internal/db/generated"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	url string
}

func NewDB(url string) *DB {
	return &DB{url: url}
}

func (d *DB) UpsertDivisionModeMapInfo(ctx context.Context, dmmi apischema.DivisionModeMapInfo) (bool, error) {
	pool, err := pgxpool.New(ctx, d.url)
	if err != nil {
		return false, err
	}
	defer pool.Close()

	q := db.New(pool)

	var divisionInfoDbId int64
	var divisionInfoChanged bool

	divisionInfoEntry, err := q.GetDivisionInfo(ctx, db.GetDivisionInfoParams{
		DivisionName: dmmi.DivisionName,
		GameMode:     dmmi.GameMode,
	})
	if err != nil {
		// division info not returned
		if err == pgx.ErrNoRows {
			// division info not stored yet - expected on DB init and after reorganizations on GeoGuessr side, insert new entry
			insertDivisionInfoParams := db.InsertDivisionInfoParams{
				DivisionName: dmmi.DivisionName,
				GameMode:     dmmi.GameMode,
				MapID:        dmmi.MapID,
				MapName:      dmmi.MapName,
			}
			divisionInfoDbId, err = q.InsertDivisionInfo(ctx, insertDivisionInfoParams)
			if err != nil {
				return false, err
			}
			divisionInfoChanged = true // new entry, consider it as "changed" to trigger any downstream processing
		} else {
			// error with DB
			return false, err
		}
	} else {
		// division info already exists, flow to check if it changed
		divisionInfoDbId = divisionInfoEntry.ID
		if divisionInfoEntry.MapID != dmmi.MapID || divisionInfoEntry.MapName != dmmi.MapName {
			divisionInfoChanged = true
		}
	}

	if divisionInfoChanged {
		tx, err := pool.Begin(ctx)
		if err != nil {
			return false, err
		}
		defer tx.Rollback(ctx)

		qTx := q.WithTx(tx)

		// FIXME: check if insert returned row
		_, err = qTx.InsertDivisionHistory(ctx, db.InsertDivisionHistoryParams{
			DivisionID: divisionInfoDbId,
			MapID:      dmmi.MapID,
			MapName:    dmmi.MapName,
		})
		if err != nil {
			return false, err
		}

		// FIXME: check if insert returned row
		_, err = qTx.UpdateDivisionInfo(ctx, db.UpdateDivisionInfoParams{
			ID:      divisionInfoDbId,
			MapID:   dmmi.MapID,
			MapName: dmmi.MapName,
		})
		if err != nil {
			return false, err
		}

		err = tx.Commit(ctx)
		if err != nil {
			return false, err
		}
	}

	return divisionInfoChanged, nil
}

func MapInfoEqualsDbEntry(mi apischema.MapInfo, dbEntry db.MapInfo) bool {
	// this may be useful to adjust logic
	return mi.Name == dbEntry.Name &&
		mi.Description == dbEntry.Description.String &&
		mi.BoundsMinLat == dbEntry.BoundsMinLat.Float64 &&
		mi.BoundsMinLng == dbEntry.BoundsMinLon.Float64 &&
		mi.BoundsMaxLat == dbEntry.BoundsMaxLat.Float64 &&
		mi.BoundsMaxLng == dbEntry.BoundsMaxLon.Float64 &&
		mi.MaxErrorDistance == dbEntry.MaxErrDistance.Int64 &&
		mi.UpdatedAt.Equal(dbEntry.ApiUpdated.Time) &&
		mi.CoordinateCount == dbEntry.LocationCount
}

func (d *DB) UpsertMapInfo(ctx context.Context, mi apischema.MapInfo) (bool, error) {
	pool, err := pgxpool.New(ctx, d.url)
	if err != nil {
		return false, err
	}
	defer pool.Close()

	q := db.New(pool)

	var mapInfoDbId int64
	var mapInfoChanged bool

	mapInfoEntry, err := q.GetMapInfo(ctx, mi.ID)
	if err != nil {
		// map info not returned
		if err == pgx.ErrNoRows {
			fmt.Printf("Map info for MapID=%s not found in DB, inserting new entry\n", mi.ID) // FIXME: change to logger
			// map info not stored yet - expected on DB init and after reorganizations on GeoGuessr side, insert new entry
			insertMapInfoParams := db.InsertMapInfoParams{
				ApiID: mi.ID,
				Name:  mi.Name,
				Description: pgtype.Text{
					String: mi.Description,
					Valid:  true,
				},
				BoundsMinLat: pgtype.Float8{
					Float64: mi.BoundsMinLat,
					Valid:   true,
				},
				BoundsMinLon: pgtype.Float8{
					Float64: mi.BoundsMinLng,
					Valid:   true,
				},
				BoundsMaxLat: pgtype.Float8{
					Float64: mi.BoundsMaxLat,
					Valid:   true,
				},
				BoundsMaxLon: pgtype.Float8{
					Float64: mi.BoundsMaxLng,
					Valid:   true,
				},
				MaxErrDistance: pgtype.Int8{
					Int64: mi.MaxErrorDistance,
					Valid: true,
				},
				ApiUpdated: pgtype.Timestamptz{
					Time:  mi.UpdatedAt,
					Valid: true,
				},
				LocationCount: mi.CoordinateCount,
			}
			mapInfoDbId, err = q.InsertMapInfo(ctx, insertMapInfoParams)
			if err != nil {
				return false, err
			}
			mapInfoChanged = true // new entry, consider it as "changed" to trigger any downstream processing
		} else {
			// error with DB
			return false, err
		}
	} else {
		// map info already exists, flow to check if it changed
		mapInfoDbId = mapInfoEntry.ID
		mapInfoChanged = !MapInfoEqualsDbEntry(mi, mapInfoEntry)
	}

	if mapInfoChanged {
		tx, err := pool.Begin(ctx)
		if err != nil {
			return false, err
		}
		defer tx.Rollback(ctx)

		qTx := q.WithTx(tx)

		// FIXME: check if insert returned row
		_, err = qTx.InsertMapHistory(ctx, db.InsertMapHistoryParams{
			MapID: mapInfoDbId,
			BoundsMinLat: pgtype.Float8{
				Float64: mi.BoundsMinLat,
				Valid:   true,
			},
			BoundsMinLon: pgtype.Float8{
				Float64: mi.BoundsMinLng,
				Valid:   true,
			},
			BoundsMaxLat: pgtype.Float8{
				Float64: mi.BoundsMaxLat,
				Valid:   true,
			},
			BoundsMaxLon: pgtype.Float8{
				Float64: mi.BoundsMaxLng,
				Valid:   true,
			},
			MaxErrDistance: pgtype.Int8{
				Int64: mi.MaxErrorDistance,
				Valid: true,
			},
			LocationCount: mi.CoordinateCount,
			ApiUpdated: pgtype.Timestamptz{
				Time:  mi.UpdatedAt,
				Valid: true,
			},
		})
		if err != nil {
			return false, err
		}

		// FIXME: check if insert returned row
		_, err = qTx.UpdateMapInfo(ctx, db.UpdateMapInfoParams{
			ApiID: mi.ID,
			Name:  mi.Name,
			Description: pgtype.Text{
				String: mi.Description,
				Valid:  true,
			},
			BoundsMinLat: pgtype.Float8{
				Float64: mi.BoundsMinLat,
				Valid:   true,
			},
			BoundsMinLon: pgtype.Float8{
				Float64: mi.BoundsMinLng,
				Valid:   true,
			},
			BoundsMaxLat: pgtype.Float8{
				Float64: mi.BoundsMaxLat,
				Valid:   true,
			},
			BoundsMaxLon: pgtype.Float8{
				Float64: mi.BoundsMaxLng,
				Valid:   true,
			},
			MaxErrDistance: pgtype.Int8{
				Int64: mi.MaxErrorDistance,
				Valid: true,
			},
			ApiUpdated: pgtype.Timestamptz{
				Time:  mi.UpdatedAt,
				Valid: true,
			},
			LocationCount: mi.CoordinateCount,
		})
		if err != nil {
			return false, err
		}

		err = tx.Commit(ctx)
		if err != nil {
			return false, err
		}
	}

	return mapInfoChanged, nil
}

func (d *DB) InsertUserFetchHistory(ctx context.Context, userID string) (int64, error) {
	pool, err := pgxpool.New(ctx, d.url)
	if err != nil {
		return 0, err
	}
	defer pool.Close()

	q := db.New(pool)

	fetchHistoryID, err := q.InsertUserFetchHistory(ctx, userID)
	if err != nil {
		return 0, err
	}

	return fetchHistoryID, nil
}

func (d *DB) InsertUserProgressHistory(ctx context.Context, progress apischema.ProgressInfo, userID string, fetchID int64) (bool, error) {
	pool, err := pgxpool.New(ctx, d.url)
	if err != nil {
		return false, err
	}
	defer pool.Close()

	q := db.New(pool)

	_, err = q.InsertUserProgressHistory(ctx, db.InsertUserProgressHistoryParams{
		UserID:         userID,
		FetchID:        fetchID,
		DivisionName:   progress.DivisionName,
		DivisionNumber: int32(progress.DivisionNumber),
		RatingOverall:  int32(progress.RatingOverall),
		RatingMoving:   int32(progress.RatingMoving),
		RatingNomove:   int32(progress.RatingNoMove),
		RatingNmpz:     int32(progress.RatingNMPZ),
		GuessedFirst:   progress.GuessedFirstRate,
		BestCountries:  progress.BestCountries,
		WorstCountries: progress.WorstCountries,
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (d *DB) InsertUserStatsHistory(ctx context.Context, stats apischema.StatsInfo, userID string, fetchID int64) (bool, error) {
	pool, err := pgxpool.New(ctx, d.url)
	if err != nil {
		return false, err
	}
	defer pool.Close()

	q := db.New(pool)

	_, err = q.InsertUserStatsHistory(ctx, db.InsertUserStatsHistoryParams{
		UserID:                  userID,
		FetchID:                 fetchID,
		RankedTeamMovingGames:   int64(stats.RankedTeamMovingGames),
		RankedTeamMovingWins:    int64(stats.RankedTeamMovingWins),
		RankedTeamNomoveGames:   int64(stats.RankedTeamNoMoveGames),
		RankedTeamNomoveWins:    int64(stats.RankedTeamNoMoveWins),
		RankedTeamNmpzGames:     int64(stats.RankedTeamNMPZGames),
		RankedTeamNmpzWins:      int64(stats.RankedTeamNMPZWins),
		RankedSoloMovingGames:   int64(stats.RankedSoloMovingGames),
		RankedSoloMovingWins:    int64(stats.RankedSoloMovingWins),
		RankedSoloNomoveGames:   int64(stats.RankedSoloNoMoveGames),
		RankedSoloNomoveWins:    int64(stats.RankedSoloNoMoveWins),
		RankedSoloNmpzGames:     int64(stats.RankedSoloNMPZGames),
		RankedSoloNmpzWins:      int64(stats.RankedSoloNMPZWins),
		UnrankedSoloMovingGames: int64(stats.UnrankedSoloMovingGames),
		UnrankedSoloMovingWins:  int64(stats.UnrankedSoloMovingWins),
		UnrankedSoloNomoveGames: int64(stats.UnrankedSoloNoMoveGames),
		UnrankedSoloNomoveWins:  int64(stats.UnrankedSoloNoMoveWins),
		UnrankedSoloNmpzGames:   int64(stats.UnrankedSoloNMPZGames),
		UnrankedSoloNmpzWins:    int64(stats.UnrankedSoloNMPZWins),
	})
	if err != nil {
		return false, err
	}

	return true, nil
}
