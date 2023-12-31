package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkstaking "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// StakingKeeper defines the expected staking keeper interface used in the
// IBC testing package
type StakingKeeper interface {
	GetHistoricalInfo(ctx sdk.Context, height int64) (sdkstaking.HistoricalInfo, bool)
}
