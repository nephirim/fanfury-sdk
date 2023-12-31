package keeper_test

import (
	"math/big"
	"testing"

	"cosmossdk.io/math"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	furyapp "github.com/nephirim/fanfury-sdk/v2/app"
	"github.com/nephirim/fanfury-sdk/v2/x/lsnative/staking/keeper"
	"github.com/nephirim/fanfury-sdk/v2/x/lsnative/staking/types"

	sdkstaking "github.com/cosmos/cosmos-sdk/x/staking/types"
)

var PKs = furyapp.CreateTestPubKeys(500)

func init() {
	sdk.DefaultPowerReduction = sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
}

// createTestInput Returns a furyapp with custom StakingKeeper
// to avoid messing with the hooks.
func createTestInput(t *testing.T) (*codec.LegacyAmino, *furyapp.FuryApp, sdk.Context) {
	app := furyapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.StakingKeeper = keeper.NewKeeper(
		app.AppCodec(),
		app.GetKey(types.StoreKey),
		app.AccountKeeper,
		app.BankKeeper,
		app.GetSubspace(types.ModuleName),
	)
	return app.LegacyAmino(), app, ctx
}

// intended to be used with require/assert:  require.True(ValEq(...))
func ValEq(t *testing.T, exp, got types.Validator) (*testing.T, bool, string, types.Validator, types.Validator) {
	return t, exp.MinEqual(&got), "expected:\n%v\ngot:\n%v", exp, got
}

// generateAddresses generates numAddrs of normal AccAddrs and ValAddrs
func generateAddresses(app *furyapp.FuryApp, ctx sdk.Context, numAddrs int) ([]sdk.AccAddress, []sdk.ValAddress) {
	addrDels := furyapp.AddTestAddrsIncremental(app, ctx, numAddrs, sdk.NewInt(10000))
	addrVals := furyapp.ConvertAddrsToValAddrs(addrDels)

	return addrDels, addrVals
}

func delegateCoinsFromAccount(ctx sdk.Context, app *furyapp.FuryApp, addr sdk.AccAddress, amount math.Int, val types.Validator) error {
	// bondDenom := app.StakingKeeper.BondDenom(ctx)
	// coins := sdk.Coins{sdk.NewCoin(bondDenom, amount)}
	// app.BankKeeper.DelegateCoinsFromAccountToModule(ctx, addr, types.EpochDelegationPoolName, coins)
	_, err := app.StakingKeeper.Delegate(ctx, addr, amount, sdkstaking.Unbonded, val, true)

	return err
}
