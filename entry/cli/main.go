//+build !test

package main

import (
	"fmt"
	"os"

	"github.com/mkuchenbecker/brewery3/brewery/utils"

	"github.com/mkuchenbecker/brewery3/brewery/cli"
	model "github.com/mkuchenbecker/brewery3/brewery/model/gomodel"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", 8100),
		grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		panic(err)
	}
	client := model.NewBreweryClient(conn)
	err = cli.Run(client, os.Args)
	if err != nil {
		utils.LogError(nil, err, "encountered an error")
		os.Exit(1)
	}
}
