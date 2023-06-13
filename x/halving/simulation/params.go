/*
 Copyright [2019] - [2021], FANFURY TECHNOLOGIES PTE. LTD. and the fanfuryCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package simulation

import (
	"fmt"
	"math/rand"

	simulationTypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/nephirim/fanfury-sdk/v2/x/halving/types"
)

const (
	keyBlockHeight = "BlockHeight"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(_ *rand.Rand) []simulationTypes.ParamChange {
	return []simulationTypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, keyBlockHeight,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GetBlockHeight(r))
			},
		),
	}
}
