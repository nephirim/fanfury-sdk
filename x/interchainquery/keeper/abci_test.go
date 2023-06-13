package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/incubus-network/fanfury-sdk/v2/x/lsnative/staking/types"

	"github.com/incubus-network/fanfury-sdk/v2/ibctesting"
	"github.com/incubus-network/fanfury-sdk/v2/x/interchainquery/keeper"
)

func (suite *KeeperTestSuite) TestEndBlocker() {
	validators := suite.GetFuryApp(suite.chainB).StakingKeeper.GetBondedValidatorsByPower(suite.chainB.GetContext())

	qvr := stakingtypes.QueryValidatorsResponse{
		Validators: ibctesting.SdkValidatorsToValidators(validators),
	}

	bondedQuery := stakingtypes.QueryValidatorsRequest{Status: stakingtypes.BondStatusBonded}
	bz, err := bondedQuery.Marshal()
	suite.NoError(err)

	id := keeper.GenerateQueryHash(suite.path.EndpointB.ConnectionID, suite.chainB.ChainID, "cosmos.staking.v1beta1.Query/Validators", bz, "")

	query := suite.GetFuryApp(suite.chainA).InterchainQueryKeeper.NewQuery(
		suite.chainA.GetContext(),
		"",
		suite.path.EndpointB.ConnectionID,
		suite.chainB.ChainID,
		"cosmos.staking.v1beta1.Query/Validators",
		bz,
		sdk.NewInt(200),
		"",
		0,
	)

	// set the query
	suite.GetFuryApp(suite.chainA).InterchainQueryKeeper.SetQuery(suite.chainA.GetContext(), *query)

	// call end blocker
	suite.GetFuryApp(suite.chainA).InterchainQueryKeeper.EndBlocker(suite.chainA.GetContext())

	err = suite.GetFuryApp(suite.chainA).InterchainQueryKeeper.SetDatapointForID(
		suite.chainA.GetContext(),
		id,
		suite.GetFuryApp(suite.chainB).AppCodec().MustMarshalJSON(&qvr),
		sdk.NewInt(suite.chainB.CurrentHeader.Height),
	)
	suite.NoError(err)

	dataPoint, err := suite.GetFuryApp(suite.chainA).InterchainQueryKeeper.GetDatapointForID(suite.chainA.GetContext(), id)
	suite.NoError(err)
	suite.NotNil(dataPoint)

	// set the query
	suite.GetFuryApp(suite.chainA).InterchainQueryKeeper.DeleteQuery(suite.chainA.GetContext(), id)

	// call end blocker
	suite.GetFuryApp(suite.chainA).InterchainQueryKeeper.EndBlocker(suite.chainA.GetContext())
}
