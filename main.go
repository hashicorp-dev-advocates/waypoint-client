package main

import (
	"context"
	"github.com/hashicorp-dev-advocates/waypoint-client/pkg/client"
	"github.com/kr/pretty"
	"log"
	"os"
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

	//gc := client.Git{
	//	Url:  "https://github.com/hashicorp/waypoint-examples",
	//	Path: "docker/go",
	//}
	//
	//projconf := client.DefaultProjectConfig()
	//
	//projconf.Name = "robbarnes"
	//projconf.RemoteRunnersEnabled = false
	//projconf.GitPollInterval = 30 * time.Second
	//
	//var1 := client.SetVariable()
	//var1.Name = "name"
	//var1.Value = &gen.Variable_Str{Str: "Devops Rob"}
	//
	//var2 := client.SetVariable()
	//var2.Name = "role"
	//var2.Value = &gen.Variable_Str{Str: "Developer Advocate"}
	//
	//var varList []*gen.Variable
	//
	//varList = append(varList, &var1, &var2)
	//projconf.StatusReportPoll = 0 * time.Second
	//
	//npr, err := wp.UpsertProject(context.TODO(), projconf, &gc, varList)
	//if err != nil {
	//	panic(err)
	//}
	//
	////gpr, err := wp.GetProject(context.TODO(), "robbarnes")
	////if err != nil {
	////	fmt.Println(err)
	////}
	//
	//prl, err := wp.ListProject(context.TODO())
	//pretty.Println(prl)
	//pretty.Println(npr)
	var envVarsMap = map[string]string{
		"VAULT_ADDR": "http://localhost:8200",
	}

	dCon := client.DefaultRunnerConfig()
	dCon.Name = "mercedes"
	dCon.OciUrl = "hashicorp/waypoint-odr:latest"
	dCon.EnvironmentVariables = envVarsMap
	dCon.Default = false
	dCon.PluginType = "docker"

	urp, err := wp.CreateRunnerProfile(context.TODO(), dCon)

	//grc, err := wp.GetRunnerProfile(context.TODO(), "01G5BS0SAKR1TVVN01QTKV7FXC")
	pretty.Println(urp.Config)
}
