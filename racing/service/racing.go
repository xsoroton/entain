package service

import (
	"github.com/xsoroton/entain/racing/db"
	"github.com/xsoroton/entain/racing/proto/racing"
	"golang.org/x/net/context"
)

type Racing interface {
	// GetRace by ID
	GetRace(ctx context.Context, in *racing.GetRaceRequest) (*racing.GetRaceResponse, error)

	// ListRaces will return a collection of races.
	ListRaces(ctx context.Context, in *racing.ListRacesRequest) (*racing.ListRacesResponse, error)
}

// racingService implements the Racing interface.
type racingService struct {
	racesRepo db.RacesRepo
}

// NewRacingService instantiates and returns a new racingService.
func NewRacingService(racesRepo db.RacesRepo) Racing {
	return &racingService{racesRepo}
}

func (s *racingService) GetRace(ctx context.Context, in *racing.GetRaceRequest) (*racing.GetRaceResponse, error) {
	race, err := s.racesRepo.GetRace(in)
	if err != nil {
		return nil, err
	}

	return &racing.GetRaceResponse{Race: race}, nil
}

func (s *racingService) ListRaces(ctx context.Context, in *racing.ListRacesRequest) (*racing.ListRacesResponse, error) {
	races, err := s.racesRepo.List(in.Filter)
	if err != nil {
		return nil, err
	}

	return &racing.ListRacesResponse{Races: races}, nil
}
