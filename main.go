package main

import (
	"context"
	"github.com/hashicorp-dev-advocates/waypoint-client/pkg/client"
	"github.com/kr/pretty"
	"log"
	"os"
	"time"
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

	//ga := &client.GitAuthBasic{
	//	Username: "",
	//	Password: "",
	//}
	gc := client.Git{
		Url:  "https://github.com/hashicorp/waypoint-examples",
		Path: "docker/go",
	}

	projconf := client.DefaultProjectConfig()

	projconf.Name = "NJackson"
	projconf.RemoteRunnersEnabled = false
	projconf.GitPollInterval = 30 * time.Second
	//projconf.DataSourcePoll = ""
	//projconf.WaypointHcl = []byte("")
	//projconf.WaypointHclFormat = 0
	//projconf.FileChangeSignal = ""
	//projconf.Variables = ""
	//projconf.StatusReportPoll = ""

	npr, err := wp.UpsertProject(context.TODO(), projconf, &gc)
	if err != nil {
		panic(err)
	}

	pretty.Println(npr)
}
