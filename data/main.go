package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/logging"
	"github.com/kelseyhightower/envconfig"
	firestoreSink "github.com/mkuchenbecker/brewery3/data/datasink/firestore"
	data "github.com/mkuchenbecker/brewery3/data/gomodel"
	"github.com/mkuchenbecker/brewery3/data/logger"
	stackdriver "github.com/mkuchenbecker/brewery3/data/logger/logger"
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
	fmt.Printf("Main Method Started\n")
	settings := getSettings("")

	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*5))
	defer cancel()

	logClient, err := logging.NewClient(ctx, settings.GcpProjectID)
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer logClient.Close()

	fmt.Println("created logger client")

	lg := logClient.Logger("defaultName")
	err = lg.LogSync(ctx, logging.Entry{Payload: "Logging Started"})
	if err != nil {
		l.WithError(err).Level(logger.Critical).Log(ctx, "failed to start logger")
		return
	}

	l := stackdriver.New(&stackdriver.Getter{Logger: lg})
	l.Log(ctx, "Started Logger")

	fireClient, err := firestore.NewClient(ctx, settings.GcpProjectID)
	if err != nil {
		l.WithError(err).Level(logger.Critical).Log(ctx, "failed to start service")
		return
	}

	fc := firestoreSink.NewFirestoreClient(fireClient)
	datasink := firestoreSink.NewStore(settings.FirestoreCollection, fc, l)

	ts := time.Now()
	_, err = datasink.Send(ctx, &data.DataObject{
		Key: fmt.Sprintf("debug:datasink-connect-%d", ts.Unix()),
		Fields: map[string](*data.Value){
			"timestamp": &data.Value{Value: &data.Value_Int64{Int64: ts.Unix()}},
			"message":   &data.Value{Value: &data.Value_String_{String_: "connected to data sink"}},
		},
	})

	if err != nil {
		l.WithError(err).Level(logger.Critical).Log(ctx, "failed to start service")
		panic(err)
	}

	fmt.Println("stored a piece of info")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", settings.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	serve := grpc.NewServer()
	data.RegisterDataProcessorServer(serve, datasink)
	// Register reflection service on gRPC server.
	reflection.Register(serve)
	l.Log(ctx, "Starting to Serve")
	if err := serve.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
