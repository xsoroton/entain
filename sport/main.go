package main

import (
	"flag"
	"net"

	"database/sql"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/xsoroton/entain/sport/db"
	"github.com/xsoroton/entain/sport/proto/sport"
	"github.com/xsoroton/entain/sport/service"
)

var (
	grpcEndpoint = flag.String("grpc-endpoint", "localhost:5000", "gRPC server endpoint")
)

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatalf("failed running grpc server: %s", err)
	}
}

func run() error {
	conn, err := net.Listen("tcp", ":5000")
	if err != nil {
		return err
	}

	sportDB, err := sql.Open("sqlite3", "./db/sports.db")
	if err != nil {
		return err
	}

	sportsRepo := db.NewSportsRepo(sportDB)
	if err := sportsRepo.Init(); err != nil {
		return err
	}

	grpcServer := grpc.NewServer()

	sport.RegisterSportsServer(
		grpcServer,
		service.NewSportService(
			sportsRepo,
		),
	)

	log.Infof("gRPC server listening on: %s", *grpcEndpoint)

	if err := grpcServer.Serve(conn); err != nil {
		return err
	}

	return nil
}
