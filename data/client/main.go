package main

import (
	"context"
	"fmt"

	"github.com/mkuchenbecker/brewery3/brewery/utils"
	data "github.com/mkuchenbecker/brewery3/data/gomodel"
	"google.golang.org/grpc"
)

func makeDataClient(address string) (data.DataProcessorClient, *grpc.ClientConn) {
	utils.Print(fmt.Sprintf("Connecting to client: %s", address))
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := data.NewDataProcessorClient(conn)
	return client, conn
}

func main() {

	do := &data.DataObject{
		Key: "testdata:abc",
		Fields: map[string]*data.Value{
			"id":     {Value: &data.Value_String_{String_: "testdata:abc"}},
			"double": {Value: &data.Value_Double{Double: 1}},
		},
	}
	client, conn := makeDataClient("34.70.183.94:9000")
	defer conn.Close()

	_, err := client.Send(context.Background(), do)

	if err != nil {
		panic(err)
	}

}
