package keeper

import (
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/nephirim/fanfury-sdk/v2/x/lsnative/staking/types"
)

// UnbondingTime
func (k Keeper) UnbondingTime(ctx sdk.Context) (res time.Duration) {
	k.paramstore.Get(ctx, types.KeyUnbondingTime, &res)
	return
}

// MaxValidators - Maximum number of validators
func (k Keeper) MaxValidators(ctx sdk.Context) (res uint32) {
	k.paramstore.Get(ctx, types.KeyMaxValidators, &res)
	return
}

// MaxEntries - Maximum number of simultaneous unbonding
// delegations or redelegations (per pair/trio)
func (k Keeper) MaxEntries(ctx sdk.Context) (res uint32) {
	k.paramstore.Get(ctx, types.KeyMaxEntries, &res)
	return
}

// HistoricalEntries = number of historical info entries
// to persist in store
func (k Keeper) HistoricalEntries(ctx sdk.Context) (res uint32) {
	k.paramstore.Get(ctx, types.KeyHistoricalEntries, &res)
	return
}

// BondDenom - Bondable coin denomination
func (k Keeper) BondDenom(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyBondDenom, &res)
	return
}

// PowerReduction - is the amount of staking tokens required for 1 unit of consensus-engine power.
// Currently, this returns a global variable that the app developer can tweak.
// TODO: we might turn this into an on-chain param:
// https://github.com/cosmos/cosmos-sdk/issues/8365
func (k Keeper) PowerReduction(ctx sdk.Context) math.Int {
	return sdk.DefaultPowerReduction
}

// MinCommissionRate - Minimum validator commission rate
func (k Keeper) MinCommissionRate(ctx sdk.Context) (res sdk.Dec) {
	k.paramstore.Get(ctx, types.KeyMinCommissionRate, &res)
	return
}

// ValidatorBondFactor - validator bond factor for all validators
func (k Keeper) ValidatorBondFactor(ctx sdk.Context) (res sdk.Dec) {
	k.paramstore.Get(ctx, types.KeyValidatorBondFactor, &res)
	return
}

// Global liquid staking cap across all liquid staking providers
func (k Keeper) GlobalLiquidStakingCap(ctx sdk.Context) (res sdk.Dec) {
	k.paramstore.Get(ctx, types.KeyGlobalLiquidStakingCap, &res)
	return
}

// Check whether the validator bond factor is enabled
// A non-negative factor indicates that it is enabled
func (k Keeper) ValidatorBondFactorEnabled(ctx sdk.Context) bool {
	return k.ValidatorBondFactor(ctx).IsPositive()
}

// Check whether the global liquid staking cap is enabled
// A cap less than 100% indicates that it is enabled
func (k Keeper) GlobalLiquidStakingCapEnabled(ctx sdk.Context) bool {
	return k.GlobalLiquidStakingCap(ctx).LT(sdk.OneDec())
}

// Get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.UnbondingTime(ctx),
		k.MaxValidators(ctx),
		k.MaxEntries(ctx),
		k.HistoricalEntries(ctx),
		k.BondDenom(ctx),
		k.MinCommissionRate(ctx),
		k.ValidatorBondFactor(ctx),
		k.GlobalLiquidStakingCap(ctx),
	)
}

// set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
