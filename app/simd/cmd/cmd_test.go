package cmd_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/incubus-network/fanfury-sdk/v2/app"
	"github.com/incubus-network/fanfury-sdk/v2/app/furyd/cmd"
	"github.com/incubus-network/fanfury-sdk/v2/x/lsnative/genutil/client/cli"
)

func TestInitCmd(t *testing.T) {
	rootCmd, _ := cmd.NewRootCmd()
	rootCmd.SetArgs([]string{
		"init",         // Test the init cmd
		"furyapp-test", // Moniker
		fmt.Sprintf("--%s=%s", cli.FlagOverwrite, "true"), // Overwrite genesis.json, in case it already exists
	})

	require.NoError(t, svrcmd.Execute(rootCmd, "", furyapp.DefaultNodeHome))
}
