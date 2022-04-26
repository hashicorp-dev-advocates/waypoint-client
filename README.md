# WIP - Waypoint Client

**DISCLAIMER - THIS PACKAGE IS NOT READY FOR PRODUCTION USUAGE. BREAKING CHANGES SHOULD BE EXPECTED**

[Hashicorp Waypoint](https://www.waypointproject.io/) has a builtin gRPC API that can be used to interact with the [Waypoint Server](https://www.waypointproject.io/docs/server); however, the SDK does not ship with a client to interact with this API. This is Go client library for Waypoint that exposes some of the gRPC methods in order to interact with the server.

## Supported methods

Currently, the library supports the following methods:

1. `GetVersion` - This gets the Waypoint version information from the server.
2. `GetProject` - This method gets information about a named Waypoint project.
3. `GetUser` - This retrieves information about a Waypoint user.
4. `CreateToken` - This will create an authentication token for a names user.
5. `InviteUser` - This will create an invitation token for a new user to be onboarded to the Waypoint server.
6. `AcceptInvitation` - This method creates the user and exchanges an invitation token for an authentication token.
7. `DeleteUser` - This will delete a named user from the Waypoint server.

**More methods are constantly being added to this library. If a method you require isn't currently supported, please do open an issue on this repo or submit a Pull Request.**

## Usage

To use this client library, import the `github.com/hashicorp-dev-advocates/waypoint-client/pkg/client` package to your project. Configuring the client can be done with the following information as a minimum requirement:

1. Waypoint token
2. Waypoint server address

In the below example, the Waypoint token is fetched from an environment variable and the address is hard coded:

```go
package main

import (
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
	
}

```

### GetVersion Example

```go
info, err := wp.GetVersionInfo(context.TODO())
if err != nil {
    panic(err)
}

fmt.Println(info.Version)
fmt.Println(info.Entrypoint)
fmt.Println(info.Api)

```

### GetProject Example

```go
gpr,err := wp.GetProject(context.TODO(), "Test")
if err != nil {
    panic(err)
}

fmt.Println(gpr.Name)
fmt.Println(gpr.Variables)
fmt.Println(gpr.Applications)
```

### CreateToken Example

```go
username := client.UserId("00000000000000000000000001")
tk, err := wp.CreateToken(context.TODO(), &username)
if err != nil {
    panic(err)
}

fmt.Println(tk)
```

### GetUser Example

```go
gu, err := wp.GetUser(context.TODO(), "DevOpsRob")
if err != nil {
panic(err)
}

fmt.Printf("Username: %s \n", gu.Username)
fmt.Printf("User ID: %s \n", gu.Id)
```

### InviteUser & AcceptInvitation Example

```go
inviteUsername := "DevOpsRob"
inv, err := wp.InviteUser(context.TODO(), inviteUsername, "30s")
if err != nil {
    panic(err)
}

tok, err := wp.AcceptInvitation(context.TODO(), inv)
if err != nil {
    panic(err)
}

fmt.Printf("Token: %s \n", tok)
```

### DeleteUser Example

```go
_, err = wp.DeleteUser(context.TODO(), client.UserId("01G13MNGG5YZ6GNDF3FSXNA18X"))
if err != nil {
    log.Fatal(err)
}
```
