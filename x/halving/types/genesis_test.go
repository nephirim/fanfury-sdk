/*
 Copyright [2019] - [2021], FANFURY TECHNOLOGIES PTE. LTD. and the fanfuryCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewGenesisState(t *testing.T) {
	params := NewParams(uint64(100))
	genesisState := NewGenesisState(params)
	require.Equal(t, &GenesisState{Params: params}, genesisState)
}

func TestDefaultGenesisState(t *testing.T) {
	params := NewParams(uint64(2 * 60 * 60 * 8766 / 5))

	require.Equal(t, &GenesisState{Params: params}, DefaultGenesisState())
}

func TestValidateGenesis(t *testing.T) {
	params := NewParams(uint64(100))
	genesisState := NewGenesisState(params)
	err := ValidateGenesis(*genesisState)
	require.Equal(t, nil, err)
}
