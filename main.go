package main

import (
	"context"
	"fmt"
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

	npr, err := wp.CreateProject(context.TODO(), "lash", false)

	fmt.Println(npr)
}
