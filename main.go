package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp-dev-advocates/waypoint-client/pkg/client"
)

var token = "HZCwuUtmrrpgkhJJuZF4ew5ektbankBX6f33w4vHw2D6emjQoaz61knyaxokhGWVj3cYTGpkFFgXEX5vyNMZgX1S3kr7EBfKsbY9yt9aDpfPZUWP5921qJw34nwRKgQxN78HAdZnhQm8agakTacnwfUs8tbq4Rh9DDxk"

func main() {

	// create a client
	conf := client.DefaultConfig()
	conf.Token = token
	conf.Address = "localhost:9701"

	wp, err := client.New(conf)
	if err != nil {
		log.Fatal(err)
	}

	info, err := wp.GetVersionInfo(context.TODO())
	if err != nil {
		panic(err)
	}

	fmt.Println(info.Version)
	fmt.Println(info.Entrypoint)
	fmt.Println(info.Api)
}
