package staking_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	tmjson "github.com/tendermint/tendermint/libs/json"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/nephirim/fanfury-sdk/v2/app"
	"github.com/nephirim/fanfury-sdk/v2/x/lsnative/staking/types"
)

func TestItCreatesModuleAccountOnInitBlock(t *testing.T) {
	db := dbm.NewMemDB()
	encCdc := furyapp.MakeTestEncodingConfig()
	app := furyapp.NewFuryApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, furyapp.DefaultNodeHome, 5, encCdc, furyapp.EmptyAppOptions{})

	genesisState := furyapp.GenesisStateWithSingleValidator(t, app)
	stateBytes, err := tmjson.Marshal(genesisState)
	require.NoError(t, err)

	app.InitChain(
		abcitypes.RequestInitChain{
			AppStateBytes: stateBytes,
			ChainId:       "test-chain-id",
		},
	)

	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	acc := app.AccountKeeper.GetAccount(ctx, authtypes.NewModuleAddress(types.BondedPoolName))
	require.NotNil(t, acc)

	acc = app.AccountKeeper.GetAccount(ctx, authtypes.NewModuleAddress(types.NotBondedPoolName))
	require.NotNil(t, acc)
}
