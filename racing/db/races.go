package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	_ "github.com/mattn/go-sqlite3"

	"github.com/xsoroton/entain/racing/proto/racing"
)

const (
	RaceStatusOpen   = "OPEN"
	RaceStatusClosed = "CLOSED"
)

// RacesRepo provides repository access to races.
type RacesRepo interface {
	// Init will initialise our races repository.
	Init() error

	// Get single race by id
	GetRace(filter *racing.GetRaceRequest) (*racing.Race, error)

	// List will return a list of races.
	List(filter *racing.ListRacesRequestFilter) ([]*racing.Race, error)
}

type racesRepo struct {
	db   *sql.DB
	init sync.Once
}

// NewRacesRepo creates a new races repository.
func NewRacesRepo(db *sql.DB) RacesRepo {
	return &racesRepo{db: db}
}

// Init prepares the race repository dummy data.
func (r *racesRepo) Init() error {
	var err error

	r.init.Do(func() {
		// For test/example purposes, we seed the DB with some dummy races.
		err = r.seed()
	})

	return err
}

func (r *racesRepo) GetRace(filter *racing.GetRaceRequest) (*racing.Race, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getRaceQueries()[racesList]
	query += fmt.Sprintf(" WHERE id = %d", filter.Id)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	races, err := r.scanRaces(rows)
	if err != nil {
		return nil, err
	}
	if len(races) == 0 {
		return nil, nil
	}
	// just in case
	if len(races) != 1 {
		log.Fatal("single row expected")
	}
	return races[0], nil
}

func (r *racesRepo) List(filter *racing.ListRacesRequestFilter) ([]*racing.Race, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getRaceQueries()[racesList]

	query, args = r.applyFilter(query, filter)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return r.scanRaces(rows)
}

func (r *racesRepo) applyFilter(query string, filter *racing.ListRacesRequestFilter) (string, []interface{}) {
	var (
		clauses []string
		args    []interface{}
	)

	if filter == nil {
		return query, args
	}

	if len(filter.MeetingIds) > 0 {
		clauses = append(clauses, "meeting_id IN ("+strings.Repeat("?,", len(filter.MeetingIds)-1)+"?)")

		for _, meetingID := range filter.MeetingIds {
			args = append(args, meetingID)
		}
	}

	if filter.VisibleOnly {
		clauses = append(clauses, "visible IS TRUE")
	}

	if len(clauses) != 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	switch filter.SortByFieldName {
	case "id":
		query += " ORDER BY id"
	case "meeting_id":
		query += " ORDER BY meeting_id"
	case "name":
		query += " ORDER BY name"
	case "number":
		query += " ORDER BY number"
	case "visible":
		query += " ORDER BY visible"
	case "advertised_start_time":
		query += " ORDER BY advertised_start_time"
	default:
		query += " ORDER BY advertised_start_time"
	}

	return query, args
}

func (m *racesRepo) scanRaces(
	rows *sql.Rows,
) ([]*racing.Race, error) {
	var races []*racing.Race
	now := time.Now()

	for rows.Next() {
		var race racing.Race
		var advertisedStart time.Time

		if err := rows.Scan(&race.Id, &race.MeetingId, &race.Name, &race.Number, &race.Visible, &advertisedStart); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		ts, err := ptypes.TimestampProto(advertisedStart)
		if err != nil {
			return nil, err
		}

		race.AdvertisedStartTime = ts

		if now.After(ts.AsTime()) {
			race.Status = RaceStatusClosed
		} else {
			race.Status = RaceStatusOpen
		}

		races = append(races, &race)
	}

	return races, nil
}
