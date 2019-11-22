package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/kelseyhightower/envconfig"
	firestoreSink "github.com/mkuchenbecker/brewery3/data/datasink/firestore"
	"github.com/mkuchenbecker/brewery3/data/gomodel/data"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func getSettings(prefix string) *Settings {
	var s Settings
	err := envconfig.Process(prefix, &s)
	if err != nil {
		log.Fatal(context.Background(), err.Error())
	}
	return &s
}

type Settings struct {
	FirestoreCollection string `envconfig:"FIRESTORE_COLLECTION" default:"global"`
	GcpProjectID        string `envconfig:"GCP_PROJECT_ID" default:"sigma-future-259702"`
	Port                int    `envconfig:"PORT" default:"9000"`
}

func main() {
	settings := getSettings("")

	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*5))
	defer cancel()

	client, err := firestore.NewClient(ctx, settings.GcpProjectID)
	if err != nil {
		panic(err)
	}

	fc := firestoreSink.NewFirestoreClient(client)
	datasink := firestoreSink.NewStore(settings.FirestoreCollection, fc)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", settings.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	serve := grpc.NewServer()
	data.RegisterDataProcessorServer(serve, datasink)
	// Register reflection service on gRPC server.
	reflection.Register(serve)
	if err := serve.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}