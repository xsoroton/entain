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

	"github.com/xsoroton/entain/sport/proto/sport"
)

const (
	SportStatusOpen   = "OPEN"
	SportStatusClosed = "CLOSED"
)

// SportsRepo provides repository access to sports.
type SportsRepo interface {
	// Init will initialise our sports repository.
	Init() error

	// Get single sport by id
	GetSport(filter *sport.GetSportRequest) (*sport.Sport, error)

	// List will return a list of sports.
	List(filter *sport.ListSportsRequestFilter) ([]*sport.Sport, error)
}

type sportsRepo struct {
	db   *sql.DB
	init sync.Once
}

// NewSportsRepo creates a new sports repository.
func NewSportsRepo(db *sql.DB) SportsRepo {
	return &sportsRepo{db: db}
}

// Init prepares the sport repository dummy data.
func (r *sportsRepo) Init() error {
	var err error

	r.init.Do(func() {
		// For test/example purposes, we seed the DB with some dummy sports.
		err = r.seed()
	})

	return err
}

func (r *sportsRepo) GetSport(filter *sport.GetSportRequest) (*sport.Sport, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getSportQueries()[sportsList]
	query += fmt.Sprintf(" WHERE id = %d", filter.Id)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	sports, err := r.scanSports(rows)
	if err != nil {
		return nil, err
	}
	if len(sports) == 0 {
		return nil, nil
	}
	// just in case
	if len(sports) != 1 {
		log.Fatal("single row expected")
	}
	return sports[0], nil
}

func (r *sportsRepo) List(filter *sport.ListSportsRequestFilter) ([]*sport.Sport, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getSportQueries()[sportsList]

	query, args = r.applyFilter(query, filter)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return r.scanSports(rows)
}

func (r *sportsRepo) applyFilter(query string, filter *sport.ListSportsRequestFilter) (string, []interface{}) {
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

func (m *sportsRepo) scanSports(
	rows *sql.Rows,
) ([]*sport.Sport, error) {
	var sports []*sport.Sport
	now := time.Now()

	for rows.Next() {
		var sport sport.Sport
		var advertisedStart time.Time

		if err := rows.Scan(&sport.Id, &sport.MeetingId, &sport.Name, &sport.Number, &sport.Visible, &advertisedStart); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		ts, err := ptypes.TimestampProto(advertisedStart)
		if err != nil {
			return nil, err
		}

		sport.AdvertisedStartTime = ts

		if now.After(ts.AsTime()) {
			sport.Status = SportStatusClosed
		} else {
			sport.Status = SportStatusOpen
		}

		sports = append(sports, &sport)
	}

	return sports, nil
}
