//+build !test

package main

import (
	"os"

	"github.com/mkuchenbecker/brewery3/brewery/cli/cli"
	model "github.com/mkuchenbecker/brewery3/brewery/model/gomodel"
	"github.com/mkuchenbecker/brewery3/brewery/utils"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := model.NewBreweryClient(conn)
	err = cli.Run(client, os.Args)
	if err != nil {
		utils.LogError(nil, err, "encountered an error")
		os.Exit(1)
	}
}
