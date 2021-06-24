package service

import (
	"github.com/xsoroton/entain/sport/db"
	"github.com/xsoroton/entain/sport/proto/sport"
	"golang.org/x/net/context"
)

type Sport interface {
	// GetSport by ID
	GetSport(ctx context.Context, in *sport.GetSportRequest) (*sport.GetSportResponse, error)

	// ListSports will return a collection of sports.
	ListSports(ctx context.Context, in *sport.ListSportsRequest) (*sport.ListSportsResponse, error)
}

// sportService implements the Sport interface.
type sportService struct {
	sportsRepo db.SportsRepo
}

// NewSportService instantiates and returns a new sportService.
func NewSportService(sportsRepo db.SportsRepo) Sport {
	return &sportService{sportsRepo}
}

func (s *sportService) GetSport(ctx context.Context, in *sport.GetSportRequest) (*sport.GetSportResponse, error) {
	Sport, err := s.sportsRepo.GetSport(in)
	if err != nil {
		return nil, err
	}

	return &sport.GetSportResponse{Sport: Sport}, nil
}

func (s *sportService) ListSports(ctx context.Context, in *sport.ListSportsRequest) (*sport.ListSportsResponse, error) {
	sports, err := s.sportsRepo.List(in.Filter)
	if err != nil {
		return nil, err
	}

	return &sport.ListSportsResponse{Sports: sports}, nil
}
