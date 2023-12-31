package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkstaking "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/suite"

	"github.com/nephirim/fanfury-sdk/v2/app"
	"github.com/nephirim/fanfury-sdk/v2/ibctesting"
	"github.com/nephirim/fanfury-sdk/v2/x/interchainquery/keeper"
	icqtypes "github.com/nephirim/fanfury-sdk/v2/x/interchainquery/types"
	stakingtypes "github.com/nephirim/fanfury-sdk/v2/x/lsnative/staking/types"
)

const TestOwnerAddress = "cosmos17dtl0mjt3t77kpuhg2edqzjpszulwhgzuj9ljs"

func init() {
	ibctesting.DefaultTestingAppInit = furyapp.SetupTestingApp
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

type KeeperTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	chainA *ibctesting.TestChain
	chainB *ibctesting.TestChain
	path   *ibctesting.Path
}

func (suite *KeeperTestSuite) GetFuryApp(chain *ibctesting.TestChain) *furyapp.FuryApp {
	app, ok := chain.App.(*furyapp.FuryApp)
	if !ok {
		panic("not sim app")
	}

	return app
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)
	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(1))
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(2))

	suite.path = newFuryAppPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupConnections(suite.path)
}

func (suite *KeeperTestSuite) TestMakeRequest() {
	bondedQuery := stakingtypes.QueryValidatorsRequest{Status: sdkstaking.BondStatusBonded}
	bz, err := bondedQuery.Marshal()
	suite.NoError(err)

	suite.GetFuryApp(suite.chainA).InterchainQueryKeeper.MakeRequest(
		suite.chainA.GetContext(),
		suite.path.EndpointB.ConnectionID,
		suite.chainB.ChainID,
		"cosmos.staking.v1beta1.Query/Validators",
		bz,
		sdk.NewInt(200),
		"",
		"",
		0,
	)

	id := keeper.GenerateQueryHash(suite.path.EndpointB.ConnectionID, suite.chainB.ChainID, "cosmos.staking.v1beta1.Query/Validators", bz, "")
	query, found := suite.GetFuryApp(suite.chainA).InterchainQueryKeeper.GetQuery(suite.chainA.GetContext(), id)
	suite.True(found)
	suite.Equal(suite.path.EndpointB.ConnectionID, query.ConnectionId)
	suite.Equal(suite.chainB.ChainID, query.ChainId)
	suite.Equal("cosmos.staking.v1beta1.Query/Validators", query.QueryType)
	suite.Equal(sdk.NewInt(200), query.Period)
	suite.Equal("", query.CallbackId)

	suite.GetFuryApp(suite.chainA).InterchainQueryKeeper.MakeRequest(
		suite.chainA.GetContext(),
		suite.path.EndpointB.ConnectionID,
		suite.chainB.ChainID,
		"cosmos.staking.v1beta1.Query/Validators",
		bz,
		sdk.NewInt(200),
		"",
		"",
		0,
	)
}

func (suite *KeeperTestSuite) TestSubmitQueryResponse() {
	bondedQuery := stakingtypes.QueryValidatorsRequest{Status: sdkstaking.BondStatusBonded}
	bz, err := bondedQuery.Marshal()
	suite.NoError(err)

	validators := suite.GetFuryApp(suite.chainB).StakingKeeper.GetBondedValidatorsByPower(suite.chainB.GetContext())
	qvr := stakingtypes.QueryValidatorsResponse{
		Validators: ibctesting.SdkValidatorsToValidators(validators),
	}

	tests := []struct {
		query       *icqtypes.Query
		setQuery    bool
		expectError error
	}{
		{
			suite.GetFuryApp(suite.chainA).InterchainQueryKeeper.
				NewQuery(
					suite.chainA.GetContext(),
					"",
					suite.path.EndpointB.ConnectionID,
					suite.chainB.ChainID,
					"cosmos.staking.v1beta1.Query/Validators",
					bz,
					sdk.NewInt(200),
					"",
					0,
				),
			true,
			nil,
		},
		{
			suite.GetFuryApp(suite.chainA).InterchainQueryKeeper.
				NewQuery(
					suite.chainA.GetContext(),
					"",
					suite.path.EndpointB.ConnectionID,
					suite.chainB.ChainID,
					"cosmos.staking.v1beta1.Query/Validators",
					bz,
					sdk.NewInt(200),
					"",
					10,
				),
			true,
			nil,
		},
		{
			suite.GetFuryApp(suite.chainA).InterchainQueryKeeper.
				NewQuery(
					suite.chainA.GetContext(),
					"",
					suite.path.EndpointB.ConnectionID,
					suite.chainB.ChainID,
					"cosmos.staking.v1beta1.Query/Validators",
					bz,
					sdk.NewInt(-200),
					"",
					0,
				),
			true,
			nil,
		},
		{
			suite.GetFuryApp(suite.chainA).InterchainQueryKeeper.
				NewQuery(
					suite.chainA.GetContext(),
					"",
					suite.path.EndpointB.ConnectionID,
					suite.chainB.ChainID,
					"cosmos.staking.v1beta1.Query/Validators",
					bz,
					sdk.NewInt(100),
					"",
					0,
				),
			false,
			nil,
		},
	}

	for _, tc := range tests {
		// set the query
		if tc.setQuery {
			suite.GetFuryApp(suite.chainA).InterchainQueryKeeper.SetQuery(suite.chainA.GetContext(), *tc.query)
		}

		icqmsgSrv := keeper.NewMsgServerImpl(suite.GetFuryApp(suite.chainA).InterchainQueryKeeper)

		qmsg := icqtypes.MsgSubmitQueryResponse{
			ChainId:     suite.chainB.ChainID,
			QueryId:     keeper.GenerateQueryHash(tc.query.ConnectionId, tc.query.ChainId, tc.query.QueryType, bz, ""),
			Result:      suite.GetFuryApp(suite.chainB).AppCodec().MustMarshalJSON(&qvr),
			Height:      suite.chainB.CurrentHeader.Height,
			FromAddress: TestOwnerAddress,
		}

		_, err = icqmsgSrv.SubmitQueryResponse(sdk.WrapSDKContext(suite.chainA.GetContext()), &qmsg)
		suite.Equal(tc.expectError, err)
	}
}

func (suite *KeeperTestSuite) TestDataPoints() {
	bondedQuery := stakingtypes.QueryValidatorsRequest{Status: sdkstaking.BondStatusBonded}
	bz, err := bondedQuery.Marshal()
	suite.NoError(err)

	validators := suite.GetFuryApp(suite.chainB).StakingKeeper.GetBondedValidatorsByPower(suite.chainB.GetContext())
	qvr := stakingtypes.QueryValidatorsResponse{
		Validators: ibctesting.SdkValidatorsToValidators(validators),
	}

	id := keeper.GenerateQueryHash(suite.path.EndpointB.ConnectionID, suite.chainB.ChainID, "cosmos.staking.v1beta1.Query/Validators", bz, "")

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

	suite.GetFuryApp(suite.chainA).InterchainQueryKeeper.DeleteDatapoint(suite.chainA.GetContext(), id)
}

func newFuryAppPath(chainA, chainB *ibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = ibctesting.TransferPort
	path.EndpointB.ChannelConfig.PortID = ibctesting.TransferPort

	return path
}
