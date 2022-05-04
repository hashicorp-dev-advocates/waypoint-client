package main

import (
	"context"
	"github.com/hashicorp-dev-advocates/waypoint-client/pkg/client"
	gen "github.com/hashicorp-dev-advocates/waypoint-client/pkg/waypoint"
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

	gc := client.Git{
		Url:  "https://github.com/hashicorp/waypoint-examples",
		Path: "docker/go",
	}

	projconf := client.DefaultProjectConfig()

	projconf.Name = "robbarnes"
	projconf.RemoteRunnersEnabled = false
	projconf.GitPollInterval = 30 * time.Second

	var1 := client.SetVariable()
	var1.Name = "name"
	var1.Value = &gen.Variable_Str{Str: "Devops Rob"}

	var2 := client.SetVariable()
	var2.Name = "role"
	var2.Value = &gen.Variable_Str{Str: "Developer Advocate"}

	var varList []*gen.Variable

	varList = append(varList, &var1, &var2)
	projconf.StatusReportPoll = &gen.Project_AppStatusPoll{
		Enabled:  true,
		Interval: "10m",
	}

	npr, err := wp.UpsertProject(context.TODO(), projconf, &gc, varList)
	if err != nil {
		panic(err)
	}

	pretty.Println(npr)
}
