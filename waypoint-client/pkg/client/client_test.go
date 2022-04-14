package client

import (
	"testing"

	"github.com/stretchr/testify/require"
)

//func setup() {
//	ws := pbmocks.WaypointServer{}
//
//	restartCh := make(chan struct{})
//	client := TestServer(t, m,
//		TestWithContext(ctx),
//		TestWithRestart(restartCh),
//	)
//}

func TestReturnsDefaultConfig(t *testing.T) {
	c := DefaultConfig()

	require.Equal(t, "localhost:9701", c.Address)
	require.Equal(t, "", c.Token)
	require.False(t, c.UseInsecureSkipVerify)
	require.Nil(t, c.TLSConfig)
}
