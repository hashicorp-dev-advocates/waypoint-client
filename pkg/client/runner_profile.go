package client

import (
	"context"
	gen "github.com/hashicorp-dev-advocates/waypoint-client/pkg/waypoint"
)

type RunnerConfig struct {
	Id                   string
	Name                 string
	TargetRunner         *gen.Ref_Runner
	OciUrl               string
	EnvironmentVariables map[string]string
	PluginType           string
	PluginConfig         []byte
	ConfigFormat         int
	Default              bool
}

func DefaultRunnerConfig() RunnerConfig {
	return RunnerConfig{
		Id:                   "",
		Name:                 "",
		TargetRunner:         &gen.Ref_Runner{Target: nil},
		OciUrl:               "",
		EnvironmentVariables: nil,
		PluginType:           "kubernetes",
		PluginConfig:         nil,
		ConfigFormat:         0,
		Default:              false,
	}
}

func (c *waypointImpl) CreateRunnerProfile(ctx context.Context, config RunnerConfig) (*gen.UpsertOnDemandRunnerConfigResponse, error) {

	odrc := &gen.OnDemandRunnerConfig{
		Id:                   config.Id,
		Name:                 config.Name,
		TargetRunner:         &gen.Ref_Runner{Target: nil},
		OciUrl:               config.OciUrl,
		EnvironmentVariables: config.EnvironmentVariables,
		PluginType:           config.PluginType,
		PluginConfig:         config.PluginConfig,
		ConfigFormat:         gen.Hcl_Format(config.ConfigFormat),
		Default:              false,
	}
	urcr := &gen.UpsertOnDemandRunnerConfigRequest{
		Config: odrc,
	}

	urc, err := c.client.UpsertOnDemandRunnerConfig(ctx, urcr)
	if err != nil {
		return nil, err
	}
	return urc, nil

}

func (c *waypointImpl) GetRunnerProfile(ctx context.Context, id string) (*gen.GetOnDemandRunnerConfigResponse, error) {

	godrc := &gen.GetOnDemandRunnerConfigRequest{
		Config: &gen.Ref_OnDemandRunnerConfig{
			Id:   id,
			Name: "",
		},
	}
	godrr, err := c.client.GetOnDemandRunnerConfig(ctx, godrc)
	if err != nil {
		return nil, err
	}

	return godrr, nil
}
