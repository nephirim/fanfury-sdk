package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/incubus-network/fanfury-sdk/v2/x/lsnative/distribution/client/cli"
)

// ProposalHandler is the community spend proposal handler.
var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal)
)
