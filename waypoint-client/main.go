package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp-dev-advocates/waypoint-client/pkg/client"
	//"main/pkg/client"

)

var token = "HZCwuUtmrrphewLAqE2kzFCmugkqtahqc8RnB3cDpnPUkK3dA2GHyNMRNsVmiXx7YRJTYPPEsVBpTz8qSEYhyG5Juie4gGLPkv9J32jZ4k9giQ5Z2uMVYEomM2WgXX4V4eQN19zmDboUiBm5sEL2M2RWPLPvwVMhmBqA"

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

	tk, err := wp.CreateToken(context.TODO())

	fmt.Println(info.Version)
	fmt.Println(info.Entrypoint)
	fmt.Println(info.Api)
	fmt.Println(tk)
}
