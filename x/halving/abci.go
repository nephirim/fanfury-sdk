/*
 Copyright [2019] - [2021], FANFURY TECHNOLOGIES PTE. LTD. and the fanfuryCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package halving

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	mintTypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/nephirim/fanfury-sdk/v2/x/halving/types"
)

func EndBlocker(ctx sdk.Context, k Keeper) {
	params := k.GetParams(ctx)

	if params.BlockHeight != 0 && uint64(ctx.BlockHeight())%params.BlockHeight == 0 {
		mintParams := k.GetMintingParams(ctx)
		newMaxInflation := mintParams.InflationMax.QuoTruncate(sdk.NewDecFromInt(Factor))
		newMinInflation := mintParams.InflationMin.QuoTruncate(sdk.NewDecFromInt(Factor))

		if newMaxInflation.Sub(newMinInflation).LT(sdk.ZeroDec()) {
			panic(fmt.Sprintf("max inflation (%s) must be greater than or equal to min inflation (%s)", newMaxInflation.String(), newMinInflation.String()))
		}

		updatedParams := mintTypes.NewParams(mintParams.MintDenom, newMaxInflation.Sub(newMinInflation), newMaxInflation, newMinInflation, mintParams.GoalBonded, mintParams.BlocksPerYear)

		k.SetMintingParams(ctx, updatedParams)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeHalving,
				sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
				sdk.NewAttribute(types.AttributeKeyNewInflationMax, updatedParams.InflationMax.String()),
				sdk.NewAttribute(types.AttributeKeyNewInflationMin, updatedParams.InflationMin.String()),
				sdk.NewAttribute(types.AttributeKeyNewInflationRateChange, updatedParams.InflationRateChange.String()),
			),
		)
	}
}
