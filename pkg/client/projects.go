package client

import (
	"context"
	gen "github.com/hashicorp-dev-advocates/waypoint-client/pkg/waypoint"
)

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

// CreateProject creates a named project on the Waypoint Server
func (c *waypointImpl) CreateProject(
	ctx context.Context,
	name string,
	remoteEnabled bool,
	) (*gen.Project, error) {

	upr := &gen.UpsertProjectRequest{
		Project: &gen.Project{
			Name:              name,
			Applications:      nil,
			RemoteEnabled:     remoteEnabled,
			DataSource:        nil,
			DataSourcePoll:    nil,
			WaypointHcl:       nil,
			WaypointHclFormat: 0,
			FileChangeSignal:  "",
			Variables:         nil,
			StatusReportPoll:  nil,
		},
	}

	cpr, err := c.client.UpsertProject(ctx, upr)
	if err != nil {
		return nil, err
	}

	return cpr.Project, nil
}
