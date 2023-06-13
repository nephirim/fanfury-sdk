package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/incubus-network/fanfury-sdk/v2/furyapp"
	"github.com/incubus-network/fanfury-sdk/v2/ibctesting"
	"github.com/incubus-network/fanfury-sdk/v2/x/interchainquery/keeper"
	"github.com/incubus-network/fanfury-sdk/v2/x/interchainquery/types"
	stakingtypes "github.com/incubus-network/fanfury-sdk/v2/x/lsnative/staking/types"
)

const TestOwnerAddress = "cosmos17dtl0mjt3t77kpuhg2edqzjpszulwhgzuj9ljs"

var (
	coordinator *ibctesting.Coordinator
	chainA      *ibctesting.TestChain
	chainB      *ibctesting.TestChain
	path        *ibctesting.Path
)

func init() {
	ibctesting.DefaultTestingAppInit = furyapp.SetupTestingApp
}

func GetFuryApp(chain *ibctesting.TestChain) *furyapp.FuryApp {
	app, ok := chain.App.(*furyapp.FuryApp)
	if !ok {
		panic("not sim app")
	}

	return app
}

func newFuryAppPath(chainA, chainB *ibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = ibctesting.TransferPort
	path.EndpointB.ChannelConfig.PortID = ibctesting.TransferPort

	return path
}

func TestMsgSubmitQueryResponse(t *testing.T) {
	coordinator = ibctesting.NewCoordinator(t, 2)
	chainA = coordinator.GetChain(ibctesting.GetChainID(1))
	chainB = coordinator.GetChain(ibctesting.GetChainID(2))
	path = newFuryAppPath(chainA, chainB)
	coordinator.SetupConnections(path)

	bondedQuery := stakingtypes.QueryValidatorsRequest{Status: stakingtypes.BondStatusBonded}
	bz, err := bondedQuery.Marshal()
	require.NoError(t, err)

	validators := GetFuryApp(chainB).StakingKeeper.GetBondedValidatorsByPower(chainB.GetContext())
	qvr := stakingtypes.QueryValidatorsResponse{
		Validators: ibctesting.SdkValidatorsToValidators(validators),
	}

	msg := types.MsgSubmitQueryResponse{
		ChainId:     chainB.ChainID + "-N",
		QueryId:     keeper.GenerateQueryHash(path.EndpointB.ConnectionID, chainB.ChainID, "cosmos.staking.v1beta1.Query/Validators", bz, ""),
		Result:      GetFuryApp(chainB).AppCodec().MustMarshalJSON(&qvr),
		Height:      chainB.CurrentHeader.Height,
		FromAddress: TestOwnerAddress,
	}

	require.NoError(t, msg.ValidateBasic())
	require.Equal(t, types.RouterKey, msg.Route())
	require.Equal(t, types.TypeMsgSubmitQueryResponse, msg.Type())
	require.Equal(t, TestOwnerAddress, msg.GetSigners()[0].String())
}
