package client

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	gen "github.com/hashicorp-dev-advocates/waypoint-client/pkg/waypoint"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	//"github.com/hashicorp/waypoint/pkg/protocolversion"
)

const (
	HeaderClientApiProtocol        = "client-api-protocol"
	HeaderClientEntrypointProtocol = "client-entrypoint-protocol"
	HeaderClientVersion            = "client-version"
)

const (
	protocolVersionApiCurrent        uint32 = 1
	protocolVersionApiMin                   = 1
	protocolVersionEntrypointCurrent uint32 = 1
	protocolVersionEntrypointMin            = 1
	currentVersion							= "0.8.1"
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

// Waypoint defines an interface for the Waypoint client
type Waypoint interface {
	GRPCClient() gen.WaypointClient
	GetVersionInfo(ctx context.Context) (*gen.VersionInfo, error)
	GetProject(ctx context.Context, name string) (*gen.Project, error)
	CreateToken(ctx context.Context) (string, error)
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
		grpc.WithUnaryInterceptor(UnaryClientInterceptor(CurrentVersion())),
		grpc.WithStreamInterceptor(StreamClientInterceptor(CurrentVersion())),

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

// CreateToken returns a waypoint token
func (c *waypointImpl) CreateToken(ctx context.Context) (string, error) {
	gtr := &gen.LoginTokenRequest{
		User:     nil,
		Trigger:  false,
	}


	token, err := c.client.GenerateLoginToken(ctx,gtr)
	if err != nil {
		return "", err
	}

	return token.Token, nil
}

func UnaryClientInterceptor(serverInfo *gen.VersionInfo) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption) error {
		ctx = metadata.AppendToOutgoingContext(ctx,
			HeaderClientApiProtocol, fmt.Sprintf(
				"%d,%d",serverInfo.Api.Minimum, serverInfo.Api.Current),
			HeaderClientEntrypointProtocol, fmt.Sprintf(
				"%d,%d",serverInfo.Entrypoint.Minimum, serverInfo.Entrypoint.Current),
			HeaderClientVersion, serverInfo.Version,
		)

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// StreamClientInterceptor returns an interceptor for the client to set
// the proper headers for stream APIs.
func StreamClientInterceptor(serverInfo *gen.VersionInfo) grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx = metadata.AppendToOutgoingContext(ctx,
			HeaderClientApiProtocol, fmt.Sprintf(
				"%d,%d",serverInfo.Api.Minimum, serverInfo.Api.Current),
			HeaderClientEntrypointProtocol, fmt.Sprintf(
				"%d,%d",serverInfo.Entrypoint.Minimum, serverInfo.Entrypoint.Current),
			HeaderClientVersion, serverInfo.Version,
		)

		return streamer(ctx, desc, cc, method, opts...)
	}
}

// Current returns the current protocol version information.
func CurrentVersion() *gen.VersionInfo {
	return &gen.VersionInfo{
		Api: &gen.VersionInfo_ProtocolVersion{
			Current: protocolVersionApiCurrent,
			Minimum: protocolVersionApiMin,
		},

		Entrypoint: &gen.VersionInfo_ProtocolVersion{
			Current: protocolVersionEntrypointCurrent,
			Minimum: protocolVersionEntrypointMin,
		},

		Version: currentVersion,
	}
}