package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp-dev-advocates/waypoint-client/pkg/client"
)

var token = "rHhYzVXQBcKnfcTUrfnFCQtxT8HiAJmQNrfiEvydhgc448hwP4rHHBBSWkoaXwpYjnYK1G6325owyNRCQoc54WmM1ZR7c1hmdjEbNdjP9Um1Fb3Qwb7h7AM1XdrQoyzt2qUBj87qPf7hhuwvAXnHK5B4VkV7y6PRa"

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

	username := client.UserId("00000000000000000000000001")
	tk, err := wp.CreateToken(context.TODO(), &username)
	if err != nil {
		panic(err)
	}


	fmt.Println(info.Version)
	fmt.Println(info.Entrypoint)
	fmt.Println(info.Api)
	fmt.Println(tk)
}
