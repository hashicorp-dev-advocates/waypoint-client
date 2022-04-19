package client

import (
	"context"
	"crypto/tls"
	"errors"

	gen "github.com/hashicorp-dev-advocates/waypoint-client/pkg/waypoint"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/emptypb"
)

var ConnectionFail error = errors.New("unable to connect to Waypoint server")

type ClientConfig struct {
	// Address of the Waypoint server
	Address string
	// Token to access the server
	Token string
	// TLSConfiguration for the server, either TLSConfig or UseInsecureSkipVerify
	// must be configured
	TLSConfig *tls.Config
	// UseInsecureSkipVerify to ignore client certificates for the server
	// either UseInsecureSkipVerify or TLSConfig must be specified
	UseInsecureSkipVerify bool
}

func DefaultConfig() ClientConfig {
	return ClientConfig{
		Address: "localhost:9701",
		Token:   "",
	}
}

// Waypoint defines and interface for the Waypoint client
type Waypoint interface {
	GRPCClient() gen.WaypointClient
	GetVersionInfo(ctx context.Context) (*gen.VersionInfo, error)
	GetProject(ctx context.Context, name string) (*gen.Project, error)
}

type waypointImpl struct {
	connection *grpc.ClientConn
	client     gen.WaypointClient
}

// New creates a new Waypoint client for the given config
//
// Upon creation the connection is established, on connection fail
// New will return an error
func New(config ClientConfig) (Waypoint, error) {
	ctx := context.Background()
	cc, err := grpc.DialContext(
		ctx,
		config.Address,
		grpc.WithPerRPCCredentials(staticToken(config.Token)),
		grpc.WithTransportCredentials(
			credentials.NewTLS(&tls.Config{InsecureSkipVerify: true}),
		),
	)

	if err != nil {
		return nil, err
	}

	for {
		s := cc.GetState()

		// If we're ready then we're done!
		if s == connectivity.Ready {
			break
		}

		// If we have a transient error and we're not retrying, then we're done.
		if s == connectivity.TransientFailure {
			cc.Close()
			return nil, ConnectionFail
		}

		if !cc.WaitForStateChange(ctx, s) {
			return nil, ConnectionFail
		}
	}

	gc := gen.NewWaypointClient(cc)

	wpc := &waypointImpl{
		connection: cc,
		client:     gc,
	}

	return wpc, nil
}

// GRPCClient returns the raw gRPC Waypoint client
func (c *waypointImpl) GRPCClient() gen.WaypointClient {
	return c.client
}

// GetVersionInfo returns the version info from the Waypoint server
func (c *waypointImpl) GetVersionInfo(ctx context.Context) (*gen.VersionInfo, error) {
	gvr, err := c.client.GetVersionInfo(context.Background(), &emptypb.Empty{}, grpc.WaitForReady(true))
	if err != nil {
		return nil, err
	}

	return gvr.Info, nil
}

// GetProject returns the project details for the given project name
func (c *waypointImpl) GetProject(ctx context.Context, name string) (*gen.Project, error) {
	gpr := &gen.GetProjectRequest{
		Project: &gen.Ref_Project{Project: name},
	}

	pr, err := c.client.GetProject(ctx, gpr)
	if err != nil {
		return nil, err
	}

	return pr.Project, nil
}

func (c *waypointImpl) CreateToken(ctx context.Context, name string) (*gen.NewTokenResponse, error) {
	gtr := &gen.LoginTokenRequest{
		User:     nil,
		Trigger:  false,
	}


	token, err := c.client.GenerateLoginToken(ctx,gtr)
	if err != nil {
		return nil, err
	}

	return token, nil
}