package main

import (
	"fmt"
	"os"
)

func main() {

	fmt.Printf("%+v\n\n", os.Args)
	// conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", 8100),
	// 	grpc.WithInsecure())
	// defer conn.Close()
	// if err != nil {
	// 	panic(err)
	// }
	// client := model.NewBreweryClient(conn)
	// run(client, os.Args)
}
