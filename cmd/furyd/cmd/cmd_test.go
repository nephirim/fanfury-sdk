package cmd_test

import (
	"fmt"
	"testing"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/persistenceOne/persistence-sdk/v2/x/lsnative/genutil/client/cli"
	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/persistence-sdk/v2/ibctesting/furyapp"
	"github.com/persistenceOne/persistence-sdk/v2/ibctesting/furyapp/furyd/cmd"
)

func TestInitCmd(t *testing.T) {
	rootCmd, _ := cmd.NewRootCmd()
	rootCmd.SetArgs([]string{
		"init",         // Test the init cmd
		"furyapp-test", // Moniker
		fmt.Sprintf("--%s=%s", cli.FlagOverwrite, "true"), // Overwrite genesis.json, in case it already exists
	})

	require.NoError(t, svrcmd.Execute(rootCmd, "furyd", furyapp.DefaultNodeHome))
}
