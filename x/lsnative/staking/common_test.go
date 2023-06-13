package staking_test

import (
	"math/big"
	"testing"

	"cosmossdk.io/math"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/nephirim/fanfury-sdk/v2/app"
	"github.com/nephirim/fanfury-sdk/v2/x/lsnative/staking/keeper"
	"github.com/nephirim/fanfury-sdk/v2/x/lsnative/staking/types"
)

func init() {
	sdk.DefaultPowerReduction = sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
}

// nolint:deadcode,unused,varcheck
var (
	priv1 = secp256k1.GenPrivKey()
	addr1 = sdk.AccAddress(priv1.PubKey().Address())
	priv2 = secp256k1.GenPrivKey()
	addr2 = sdk.AccAddress(priv2.PubKey().Address())

	valKey  = ed25519.GenPrivKey()
	valAddr = sdk.AccAddress(valKey.PubKey().Address())

	commissionRates = types.NewCommissionRates(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec())

	PKs = furyapp.CreateTestPubKeys(500)
)

// getBaseFuryappWithCustomKeeper Returns a furyapp with custom StakingKeeper
// to avoid messing with the hooks.
func getBaseFuryappWithCustomKeeper(t *testing.T) (*codec.LegacyAmino, *furyapp.FuryApp, sdk.Context) {
	app := furyapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	appCodec := app.AppCodec()

	app.StakingKeeper = keeper.NewKeeper(
		appCodec,
		app.GetKey(types.StoreKey),
		app.AccountKeeper,
		app.BankKeeper,
		app.GetSubspace(types.ModuleName),
	)
	app.StakingKeeper.SetParams(ctx, types.DefaultParams())

	return codec.NewLegacyAmino(), app, ctx
}

// generateAddresses generates numAddrs of normal AccAddrs and ValAddrs
func generateAddresses(app *furyapp.FuryApp, ctx sdk.Context, numAddrs int, accAmount math.Int) ([]sdk.AccAddress, []sdk.ValAddress) {
	addrDels := furyapp.AddTestAddrsIncremental(app, ctx, numAddrs, accAmount)
	addrVals := furyapp.ConvertAddrsToValAddrs(addrDels)

	return addrDels, addrVals
}
