package client

import (
	"context"
	gen "github.com/hashicorp-dev-advocates/waypoint-client/pkg/waypoint"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestReturnsDefaultProjectConfig(t *testing.T) {
	c := DefaultProjectConfig()

	require.Equal(t, "", c.Name)
	require.Equal(t, false, c.RemoteRunnersEnabled)
	require.Equal(t, time.Duration(0), c.GitPollInterval)
	require.Nil(t, c.WaypointHcl)
	require.Equal(t, gen.Hcl_Format(0), c.WaypointHclFormat)
	require.Equal(t, "", c.FileChangeSignal)
	require.Nil(t, c.Variables)
	require.Nil(t, c.StatusReportPoll)
}

func TestDefaultApplicationConfig(t *testing.T) {
	c := DefaultApplicationConfig()

	require.Equal(t, "", c.Name)
	require.Equal(t, "", c.FileChangeSignal)

}

func TestApplicationConfigUpdate(t *testing.T) {
	c := DefaultApplicationConfig()

	c.Name = "demo-app"
	c.FileChangeSignal = "version"

	require.Equal(t, "demo-app", c.Name)
	require.Equal(t, "version", c.FileChangeSignal)
}

func TestSetVariable(t *testing.T) {
	v := SetVariable()

	require.Equal(t, "", v.Name)
	require.Equal(t, &gen.Variable_Str{Str: ""}, v.Value)
}

func TestProjectConfigUpdate(t *testing.T) {

	var1 := SetVariable()
	var1.Name = "test-key"
	var1.Value = &gen.Variable_Str{Str: "test-value"}

	var2 := SetVariable()
	var2.Name = "who-are-we"
	var2.Value = &gen.Variable_Str{Str: "Developer Advocates"}
	var varList []*gen.Variable

	varList = append(varList, &var1, &var2)

	c := DefaultProjectConfig()
	c.Name = "dev-advocates"
	c.RemoteRunnersEnabled = true
	c.GitPollInterval = time.Duration(10)
	c.Variables = varList

	require.Equal(t, "dev-advocates", c.Name)
	require.True(t, c.RemoteRunnersEnabled)
	require.Equal(t, time.Duration(10), c.GitPollInterval)
	require.Equal(t, varList, c.Variables)
}

func TestUpsertProjectSetsNilAuth(t *testing.T) {
	p := DefaultProjectConfig()
	p.Name = "test"

	g := Git{
		Url:  "https://github.com/hashicorp/waypoint-examples",
		Path: "docker/go",
	}

	c, m := setupTests(t)

	m.On("UpsertProject", mock.Anything, mock.Anything).Return(&gen.UpsertProjectResponse{}, nil)

	_, err := c.UpsertProject(context.TODO(), p, &g, nil)
	require.NoError(t, err)

	m.AssertCalled(t, "UpsertProject", mock.Anything, mock.Anything)
	projRequest := m.Calls[0].Arguments.Get(1).(*gen.UpsertProjectRequest)

	require.Nil(t, projRequest.Project.DataSource.Source.(*gen.Job_DataSource_Git).Git.Auth)

}

func TestUpsertProjectSetsGitBasicAuth(t *testing.T) {
	p := DefaultProjectConfig()
	p.Name = "test"

	g := Git{
		Url:  "https://github.com/hashicorp/waypoint-examples",
		Path: "docker/go",
		Auth: &GitAuthBasic{
			Username: "",
			Password: "",
		},
	}

	c, m := setupTests(t)

	m.On("UpsertProject", mock.Anything, mock.Anything).Return(&gen.UpsertProjectResponse{}, nil)

	_, err := c.UpsertProject(context.TODO(), p, &g, nil)
	require.NoError(t, err)

	m.AssertCalled(t, "UpsertProject", mock.Anything, mock.Anything)
	projRequest := m.Calls[0].Arguments.Get(1).(*gen.UpsertProjectRequest)

	require.NotNil(t, projRequest.Project.DataSource.Source.(*gen.Job_DataSource_Git).Git.Auth)
	require.Equal(t, "docker/go", projRequest.Project.DataSource.Source.(*gen.Job_DataSource_Git).Git.Path)
	require.Equal(t, "https://github.com/hashicorp/waypoint-examples", projRequest.Project.DataSource.Source.(*gen.Job_DataSource_Git).Git.Url)
}
