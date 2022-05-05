package client

import (
	"context"
	gen "github.com/hashicorp-dev-advocates/waypoint-client/pkg/waypoint"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// GetVersionInfo returns the version info from the Waypoint server
func (c *waypointImpl) GetVersionInfo(ctx context.Context) (*gen.VersionInfo, error) {
	gvr, err := c.client.GetVersionInfo(context.Background(), &emptypb.Empty{}, grpc.WaitForReady(true))
	if err != nil {
		return nil, err
	}

	return gvr.Info, nil
}


