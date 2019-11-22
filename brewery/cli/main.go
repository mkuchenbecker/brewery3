//+build !test

package main

func main() {
	// conn, err := grpc.Dial(fmt.Sprintf("brewpi:9000"),
	// 	grpc.WithInsecure())
	// if err != nil {
	// 	panic(err)
	// }
	// defer utils.DeferErrReturn(conn.Close, &err)
	// client := model.NewBreweryClient(conn)
	// err = cli.Run(client, os.Args)
	// if err != nil {
	// 	utils.LogError(nil, err, "encountered an error")
	// 	os.Exit(1)
	// }
}
