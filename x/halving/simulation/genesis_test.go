/*
 Copyright [2019] - [2021], FANFURY TECHNOLOGIES PTE. LTD. and the fanfuryCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package simulation

import (
	"encoding/json"
	"math/rand"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/nephirim/fanfury-sdk/v2/x/halving/types"
	"github.com/stretchr/testify/require"
)

func TestRandomizedGenState(t *testing.T) {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)

	s := rand.NewSource(1)
	r := rand.New(s) //nolint:gosec,testfile

	simState := module.SimulationState{
		AppParams:    make(simtypes.AppParams),
		Cdc:          cdc,
		Rand:         r,
		NumBonded:    3,
		Accounts:     simtypes.RandomAccounts(r, 3),
		InitialStake: sdkmath.NewInt(1000),
		GenState:     make(map[string]json.RawMessage),
	}

	RandomizedGenState(&simState)

	var halvingGenesis types.GenesisState

	simState.Cdc.MustUnmarshalJSON(simState.GenState[types.ModuleName], &halvingGenesis)
	require.Equal(t, uint64(2*6311520), halvingGenesis.Params.BlockHeight)
}
