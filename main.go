package main

import (
	"context"
	"log"
	"os"

	"github.com/hashicorp-dev-advocates/waypoint-client/pkg/client"
)

var token string

func main() {

	token = os.Getenv("WAYPOINT_TOKEN")
	if token == "" {
		log.Fatal("WAYPOINT_TOKEN environment variable not set")
	}

	// create a client
	conf := client.DefaultConfig()
	conf.Token = token
	conf.Address = "localhost:9701"

	wp, err := client.New(conf)
	if err != nil {
		log.Fatal(err)
	}

	gu, err := wp.GetUser(context.TODO(), "DevOpsRob")
	if err != nil {
		panic(err)
	}


	_, err = wp.DeleteUser(context.TODO(), client.UserId(gu.Id))
	if err != nil {
		log.Fatal(err)
	}


}
