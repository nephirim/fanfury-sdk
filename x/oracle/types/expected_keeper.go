package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	sdkstaking "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// StakingKeeper defines the expected interface contract defined by the x/staking
// module.
type StakingKeeper interface {
	Validator(ctx sdk.Context, address sdk.ValAddress) sdkstaking.ValidatorI
	GetBondedValidatorsByPower(ctx sdk.Context) []sdkstaking.Validator
	TotalBondedTokens(sdk.Context) math.Int
	Slash(sdk.Context, sdk.ConsAddress, int64, int64, sdk.Dec) math.Int
	Jail(sdk.Context, sdk.ConsAddress)
	ValidatorsPowerStoreIterator(ctx sdk.Context) sdk.Iterator
	MaxValidators(sdk.Context) uint32
	PowerReduction(ctx sdk.Context) (res math.Int)
}

// DistributionKeeper defines the expected interface contract defined by the
// x/distribution module.
type DistributionKeeper interface {
	AllocateTokensToValidator(ctx sdk.Context, val sdkstaking.ValidatorI, tokens sdk.DecCoins)
	GetValidatorOutstandingRewardsCoins(ctx sdk.Context, val sdk.ValAddress) sdk.DecCoins
}

// AccountKeeper defines the expected interface contract defined by the x/auth
// module.
type AccountKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, moduleName string) authtypes.ModuleAccountI

	// only used for simulation
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
}

// BankKeeper defines the expected interface contract defined by the x/bank
// module.
type BankKeeper interface {
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule string, recipientModule string, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	GetDenomMetaData(ctx sdk.Context, denom string) (banktypes.Metadata, bool)
	SetDenomMetaData(ctx sdk.Context, denomMetaData banktypes.Metadata)

	// only used for simulation
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
}
