package main

import (
	"context"
	"flag"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	log "github.com/sirupsen/logrus"
	"github.com/xsoroton/entain/api/proto/racing"
	"github.com/xsoroton/entain/api/proto/sport"
	"google.golang.org/grpc"
)

var (
	apiEndpoint        = flag.String("api-endpoint", "localhost:8000", "API endpoint")
	racingGRPCEndpoint = flag.String("racing-grpc-endpoint", "localhost:9000", "racing - gRPC server endpoint")
	sportsGRPCEndpoint = flag.String("sports-grpc-endpoint", "localhost:5000", "sports - gRPC server endpoint")
)

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatalf("failed running api server: %s", err)
	}
}

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()

	if err := racing.RegisterRacingHandlerFromEndpoint(
		ctx,
		mux,
		*racingGRPCEndpoint,
		[]grpc.DialOption{grpc.WithInsecure()},
	); err != nil {
		return err
	}
	if err := sport.RegisterSportsHandlerFromEndpoint(
		ctx,
		mux,
		*sportsGRPCEndpoint,
		[]grpc.DialOption{grpc.WithInsecure()},
	); err != nil {
		return err
	}

	log.Infof("API server listening on: %s", *apiEndpoint)

	return http.ListenAndServe(*apiEndpoint, mux)
}
