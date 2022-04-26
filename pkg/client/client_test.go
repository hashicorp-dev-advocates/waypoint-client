package client

import (
	"context"
	"fmt"
	"testing"

	gen "github.com/hashicorp-dev-advocates/waypoint-client/pkg/waypoint"
	"github.com/hashicorp-dev-advocates/waypoint-client/pkg/waypoint/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupTests(t *testing.T) (*waypointImpl, *mocks.WaypointClient) {
	grpcMock := &mocks.WaypointClient{}
	client := &waypointImpl{client: grpcMock}

	return client, grpcMock
}

func TestReturnsDefaultConfig(t *testing.T) {
	c := DefaultConfig()

	require.Equal(t, "localhost:9701", c.Address)
	require.Equal(t, "", c.Token)
	require.False(t, c.UseInsecureSkipVerify)
	require.Nil(t, c.TLSConfig)
}

func TestGetVersionInfoReturnsInfo(t *testing.T) {
	// setup the client with a mock grpc client
	c, m := setupTests(t)

	// setup the mock call to return the version
	gvr := &gen.GetVersionInfoResponse{Info: &gen.VersionInfo{}}
	m.On("GetVersionInfo", mock.Anything, mock.Anything, mock.Anything).Return(gvr, nil)

	vi, err := c.GetVersionInfo(context.TODO())
	require.NoError(t, err)

	// check that the function returns the info from the GetVersionInfo call
	require.Equal(t, gvr.Info, vi)
}

func TestGetVersionInfoReturnsError(t *testing.T) {
	// setup the client with a mock grpc client
	c, m := setupTests(t)

	// setup the mock call to return the error
	m.On("GetVersionInfo", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("unable to make call"))

	_, err := c.GetVersionInfo(context.TODO())
	require.Error(t, err)
}
