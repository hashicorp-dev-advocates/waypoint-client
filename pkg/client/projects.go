package client

import (
	"context"
	gen "github.com/hashicorp-dev-advocates/waypoint-client/pkg/waypoint"
	"time"
)

type ProjectConfig struct {
	// Name of Waypoint project
	Name string
	// List of applications associated with Waypoint project
	Applications []*gen.Application
	// Whether remote runners are enabled or not
	RemoteRunnersEnabled bool
	// Where data is sourced for remote operations
	DataSource *gen.Job_DataSource
	// Polling settings.
	// Polling will trigger a "waypoint up" whenever a new
	// data is detected on the data source.
	GitPollInterval time.Duration
	// The contents of a default waypoint.hcl file.
	// This will be used ONLY IF this project does not have
	// a waypoint.hcl file when an operation is executed.
	WaypointHcl       []byte
	WaypointHclFormat gen.Hcl_Format
	// Indicates signal to be sent to any applications
	// when their config files change.
	FileChangeSignal string
	// Variable values stored on the server.
	Variables []*gen.Variable
	// Application polling settings.
	// Polling will trigger a "StatusFunc" for collecting
	// a report on the current status of the application.
	StatusReportPoll *gen.Project_AppStatusPoll
}

func DefaultProjectConfig() ProjectConfig {
	return ProjectConfig{
		Name:                 "",
		Applications:         nil,
		RemoteRunnersEnabled: false,
		//DataSource:        nil,
		GitPollInterval:   0,
		WaypointHcl:       nil,
		WaypointHclFormat: 0,
		FileChangeSignal:  "",
		Variables:         nil,
		StatusReportPoll:  nil,
	}
}

type ApplicationConfig struct {
	Name             string
	FileChangeSignal string
}

func DefaultApplicationConfig() ApplicationConfig {
	return ApplicationConfig{
		Name:             "",
		FileChangeSignal: "",
	}
}

type DataSourceConfig struct {
	JobDataSource *gen.Job_DataSource
}

type Local string

func DefaultDataSourceConfig() DataSourceConfig {

	return DataSourceConfig{
		JobDataSource: &gen.Job_DataSource{
			Source: &gen.Job_DataSource_Local{Local: &gen.Job_Local{}},
		},
	}
}

type Git struct {
	Url                      string
	Path                     string
	IgnoreChangesOutsidePath bool
	Ref                      string
	Auth                     Auth
}

type Auth interface {
	getProto() interface{}
}

type GitAuthBasic struct {
	Username string
	Password string
}

func (g *GitAuthBasic) getProto() interface{} {
	return &gen.Job_Git_Basic_{Basic: &gen.Job_Git_Basic{
		Username: g.Username,
		Password: g.Password,
	}}
}

type GitAuthSsh struct {
	PrivateKeyPem []byte
	Password      string
	User          string
}

func (g *GitAuthSsh) getProto() interface{} {
	return &gen.Job_Git_Ssh{Ssh: &gen.Job_Git_SSH{
		PrivateKeyPem: g.PrivateKeyPem,
		Password:      g.Password,
		User:          g.User,
	}}
}

type DataSourceRef interface {
	Ref() string
}

type DataSourceGit string

func (ds *DataSourceGit) Ref() string {
	return string(*ds)
}

type DataSourceLocal string

func (ds *DataSourceLocal) Ref() string {
	return string(*ds)
}

// UpsertProject creates or updates a named project on the Waypoint Server
// gc := &Git{Url: "blah",Auth: &GitAuthSsh{PrivateKeyPem: "blah"}}
func (c *waypointImpl) UpsertProject(
	ctx context.Context,
	projectConfig ProjectConfig,
	gitConfig *Git,
) (*gen.Project, error) {

	jobGit := &gen.Job_Git{
		Url:                      gitConfig.Url,
		Ref:                      gitConfig.Ref,
		Path:                     gitConfig.Path,
		IgnoreChangesOutsidePath: gitConfig.IgnoreChangesOutsidePath,
	}

	if gitConfig.Auth != nil {
		switch t := gitConfig.Auth.getProto().(type) {
		case *gen.Job_Git_Basic_:
			jobGit.Auth = t
		case *gen.Job_Git_Ssh:
			jobGit.Auth = t
		}
	}

	var datasource *gen.Job_DataSource

	datasource = &gen.Job_DataSource{
		Source: &gen.Job_DataSource_Git{
			Git: jobGit,
		},
	}

	poll := &gen.Project_Poll{
		Enabled:  projectConfig.GitPollInterval > 0,
		Interval: projectConfig.GitPollInterval.String(),
	}

	upr := &gen.UpsertProjectRequest{
		Project: &gen.Project{
			Name:              projectConfig.Name,
			RemoteEnabled:     projectConfig.RemoteRunnersEnabled,
			DataSource:        datasource,
			DataSourcePoll:    poll,
			WaypointHcl:       projectConfig.WaypointHcl,
			WaypointHclFormat: projectConfig.WaypointHclFormat,
			FileChangeSignal:  projectConfig.FileChangeSignal,
			Variables:         projectConfig.Variables,
			StatusReportPoll:  projectConfig.StatusReportPoll,
		},
	}

	cpr, err := c.client.UpsertProject(ctx, upr)
	if err != nil {
		return nil, err
	}

	return cpr.Project, nil
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
