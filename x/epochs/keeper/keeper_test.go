package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/nephirim/fanfury-sdk/v2/app"
	"github.com/nephirim/fanfury-sdk/v2/x/epochs/types"
)

type KeeperTestSuite struct {
	furyapp.KeeperTestHelper
	queryClient types.QueryClient
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.Setup()
	suite.queryClient = types.NewQueryClient(suite.QueryHelper)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
