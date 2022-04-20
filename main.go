package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp-dev-advocates/waypoint-client/pkg/client"
)

var token = "HZCwuUtmrrpQ842t6eGMknRhErt2svw42Yt7KitMnvJdUiUu1tsDQS5ntrbKRkAMkye8Tk6eLi48K3FCsQp8GSstXAti9zJBLgF6v1yGnpSdfkVck322LxVqy3hFWGTDF764tmMA85TBdy7PQG6hmSxqz9i5z6Xp8NSC"

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

	//username := client.UserId("00000000000000000000000001")
	//tk, err := wp.CreateToken(context.TODO(), &username)
	//if err != nil {
	//	panic(err)
	//}

	// Invite User
	//inviteUsername := "Robert"
	//inv, err := wp.InviteUser(context.TODO(), inviteUsername)
	//if err != nil {
	//	panic(err)
	//}

	// Accept User
	//tok, err := wp.AcceptInvitation(context.TODO(), inv)
	//

	// Delete User
	du, err := wp.DeleteUser(context.TODO(), client.UserId("01G13MNGG5YZ6GNDF3FSXNA18X"))

	fmt.Println(info.Version)
	fmt.Println(info.Entrypoint)
	fmt.Println(info.Api)
	//fmt.Println(tk)
	fmt.Println(du)
}
