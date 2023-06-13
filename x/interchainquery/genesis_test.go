package interchainquery_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/persistenceOne/persistence-sdk/v2/furyapp"
	ibctesting "github.com/persistenceOne/persistence-sdk/v2/ibctesting"
	"github.com/persistenceOne/persistence-sdk/v2/x/interchainquery"
	"github.com/persistenceOne/persistence-sdk/v2/x/interchainquery/keeper"
	"github.com/persistenceOne/persistence-sdk/v2/x/interchainquery/types"
	stakingtypes "github.com/persistenceOne/persistence-sdk/v2/x/lsnative/staking/types"
)

func init() {
	ibctesting.DefaultTestingAppInit = furyapp.SetupTestingApp
}

func TestInterChainQueryTestSuite(t *testing.T) {
	suite.Run(t, new(InterChainQueryTestSuite))
}

type InterChainQueryTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	chainA *ibctesting.TestChain
	chainB *ibctesting.TestChain
	path   *ibctesting.Path
}

func (suite *InterChainQueryTestSuite) GetFuryApp(chain *ibctesting.TestChain) *furyapp.FuryApp {
	app, ok := chain.App.(*furyapp.FuryApp)
	if !ok {
		panic("not sim app")
	}

	return app
}

func (suite *InterChainQueryTestSuite) SetupTest() {
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)
	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(1))
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(2))

	suite.path = newFuryAppPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupConnections(suite.path)
}

func (suite *InterChainQueryTestSuite) TestInitGenesis() {
	bondedQuery := stakingtypes.QueryValidatorsRequest{Status: stakingtypes.BondStatusBonded}
	bz, err := bondedQuery.Marshal()
	suite.NoError(err)

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

	interchainquery.InitGenesis(suite.chainA.GetContext(), suite.GetFuryApp(suite.chainA).InterchainQueryKeeper, types.GenesisState{Queries: []types.Query{*query}})

	id := keeper.GenerateQueryHash(suite.path.EndpointB.ConnectionID, suite.chainB.ChainID, "cosmos.staking.v1beta1.Query/Validators", bz, "")
	queryResponse, found := suite.GetFuryApp(suite.chainA).InterchainQueryKeeper.GetQuery(suite.chainA.GetContext(), id)
	suite.True(found)
	suite.Equal(suite.path.EndpointB.ConnectionID, queryResponse.ConnectionId)
	suite.Equal(suite.chainB.ChainID, queryResponse.ChainId)
	suite.Equal("cosmos.staking.v1beta1.Query/Validators", queryResponse.QueryType)
	suite.Equal(sdk.NewInt(200), queryResponse.Period)
	suite.Equal(uint64(0), queryResponse.Ttl)
	suite.Equal("", queryResponse.CallbackId)
}

func newFuryAppPath(chainA, chainB *ibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = ibctesting.TransferPort
	path.EndpointB.ChannelConfig.PortID = ibctesting.TransferPort

	return path
}
