package ibctesting

import (
	"fmt"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/incubus-network/fanfury-sdk/v2/x/lsnative/staking/teststaking"
	stakingtypes "github.com/incubus-network/fanfury-sdk/v2/x/lsnative/staking/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmprotoversion "github.com/tendermint/tendermint/proto/tendermint/version"
	tmtypes "github.com/tendermint/tendermint/types"
	tmversion "github.com/tendermint/tendermint/version"

	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	commitmenttypes "github.com/cosmos/ibc-go/v6/modules/core/23-commitment/types"
	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"
	"github.com/cosmos/ibc-go/v6/modules/core/exported"
	"github.com/cosmos/ibc-go/v6/modules/core/types"
	ibctmtypes "github.com/cosmos/ibc-go/v6/modules/light-clients/07-tendermint/types"
	"github.com/incubus-network/fanfury-sdk/v2/ibctesting/furyapp"
	"github.com/incubus-network/fanfury-sdk/v2/ibctesting/mock"
)

var MaxAccounts = 10

type SenderAccount struct {
	SenderPrivKey cryptotypes.PrivKey
	SenderAccount authtypes.AccountI
}

// TestChain is a testing struct that wraps a furyapp with the last TM Header, the current ABCI
// header and the validators of the TestChain. It also contains a field called ChainID. This
// is the clientID that *other* chains use to refer to this TestChain. The SenderAccount
// is used for delivering transactions through the application state.
// NOTE: the actual application uses an empty chain-id for ease of testing.
type TestChain struct {
	*testing.T

	Coordinator   *Coordinator
	App           TestingApp
	ChainID       string
	LastHeader    *ibctmtypes.Header // header for last block height committed
	CurrentHeader tmproto.Header     // header for current block height
	QueryServer   types.QueryServer
	TxConfig      client.TxConfig
	Codec         codec.BinaryCodec

	Vals     *tmtypes.ValidatorSet
	NextVals *tmtypes.ValidatorSet

	// Signers is a map from validator address to the PrivValidator
	// The map is converted into an array that is the same order as the validators right before signing commit
	// This ensures that signers will always be in correct order even as validator powers change.
	// If a test adds a new validator after chain creation, then the signer map must be updated to include
	// the new PrivValidator entry.
	Signers map[string]tmtypes.PrivValidator

	// autogenerated sender private key
	SenderPrivKey cryptotypes.PrivKey
	SenderAccount authtypes.AccountI

	SenderAccounts []SenderAccount
}

// NewTestChainWithValSet initializes a new TestChain instance with the given validator set
// and signer array. It also initializes 10 Sender accounts with a balance of 10000000000000000000 coins of
// bond denom to use for tests.
//
// The first block height is committed to state in order to allow for client creations on
// counterparty chains. The TestChain will return with a block height starting at 2.
//
// Time management is handled by the Coordinator in order to ensure synchrony between chains.
// Each update of any chain increments the block header time for all chains by 5 seconds.
//
// NOTE: to use a custom sender privkey and account for testing purposes, replace and modify this
// constructor function.
//
// CONTRACT: Validator array must be provided in the order expected by Tendermint.
// i.e. sorted first by power and then lexicographically by address.
func NewTestChainWithValSet(t *testing.T, coord *Coordinator, chainID string, valSet *tmtypes.ValidatorSet, signers map[string]tmtypes.PrivValidator) *TestChain {
	genAccs := []authtypes.GenesisAccount{}
	genBals := []banktypes.Balance{}
	senderAccs := []SenderAccount{}

	// generate genesis accounts
	for i := 0; i < MaxAccounts; i++ {
		senderPrivKey := secp256k1.GenPrivKey()
		acc := authtypes.NewBaseAccount(senderPrivKey.PubKey().Address().Bytes(), senderPrivKey.PubKey(), uint64(i), 0)
		amount, ok := sdk.NewIntFromString("10000000000000000000")
		require.True(t, ok)

		// add sender account
		balance := banktypes.Balance{
			Address: acc.GetAddress().String(),
			Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, amount)),
		}

		genAccs = append(genAccs, acc)
		genBals = append(genBals, balance)

		senderAcc := SenderAccount{
			SenderAccount: acc,
			SenderPrivKey: senderPrivKey,
		}

		senderAccs = append(senderAccs, senderAcc)
	}

	app := SetupWithGenesisValSet(t, valSet, genAccs, chainID, sdk.DefaultPowerReduction, genBals...)

	// create current header and call begin block
	header := tmproto.Header{
		ChainID: chainID,
		Height:  1,
		Time:    coord.CurrentTime.UTC(),
	}

	txConfig := app.GetTxConfig()

	// create an account to send transactions from
	chain := &TestChain{
		T:              t,
		Coordinator:    coord,
		ChainID:        chainID,
		App:            app,
		CurrentHeader:  header,
		QueryServer:    app.GetIBCKeeper(),
		TxConfig:       txConfig,
		Codec:          app.AppCodec(),
		Vals:           valSet,
		NextVals:       valSet,
		Signers:        signers,
		SenderPrivKey:  senderAccs[0].SenderPrivKey,
		SenderAccount:  senderAccs[0].SenderAccount,
		SenderAccounts: senderAccs,
	}

	coord.CommitBlock(chain)

	return chain
}

// NewTestChain initializes a new test chain with a default of 4 validators
// Use this function if the tests do not need custom control over the validator set
func NewTestChain(t *testing.T, coord *Coordinator, chainID string) *TestChain {
	// generate validators private/public key
	var (
		validatorsPerChain = 4
		validators         []*tmtypes.Validator
		signersByAddress   = make(map[string]tmtypes.PrivValidator, validatorsPerChain)
	)

	for i := 0; i < validatorsPerChain; i++ {
		privVal := mock.NewPV()
		pubKey, err := privVal.GetPubKey()
		require.NoError(t, err)
		validators = append(validators, tmtypes.NewValidator(pubKey, 1))
		signersByAddress[pubKey.Address().String()] = privVal
	}

	// construct validator set;
	// Note that the validators are sorted by voting power
	// or, if equal, by address lexical order
	valSet := tmtypes.NewValidatorSet(validators)

	return NewTestChainWithValSet(t, coord, chainID, valSet, signersByAddress)
}

// GetContext returns the current context for the application.
func (chain *TestChain) GetContext() sdk.Context {
	return chain.App.GetBaseApp().NewContext(false, chain.CurrentHeader)
}

// GetFuryApp returns the FuryApp to allow usage ofnon-interface fields.
// CONTRACT: This function should not be called by third parties implementing
// their own FuryApp.
func (chain *TestChain) GetFuryApp() *furyapp.FuryApp {
	app, ok := chain.App.(*furyapp.FuryApp)
	require.True(chain.T, ok)

	return app
}

// QueryProof performs an abci query with the given key and returns the proto encoded merkle proof
// for the query and the height at which the proof will succeed on a tendermint verifier.
func (chain *TestChain) QueryProof(key []byte) ([]byte, clienttypes.Height) {
	return chain.QueryProofAtHeight(key, chain.App.LastBlockHeight())
}

// QueryProof performs an abci query with the given key and returns the proto encoded merkle proof
// for the query and the height at which the proof will succeed on a tendermint verifier.
func (chain *TestChain) QueryProofAtHeight(key []byte, height int64) ([]byte, clienttypes.Height) {
	res := chain.App.Query(abci.RequestQuery{
		Path:   fmt.Sprintf("store/%s/key", host.StoreKey),
		Height: height - 1,
		Data:   key,
		Prove:  true,
	})

	merkleProof, err := commitmenttypes.ConvertProofs(res.ProofOps)
	require.NoError(chain.T, err)

	proof, err := chain.App.AppCodec().Marshal(&merkleProof)
	require.NoError(chain.T, err)

	revision := clienttypes.ParseChainID(chain.ChainID)

	// proof height + 1 is returned as the proof created corresponds to the height the proof
	// was created in the IAVL tree. Tendermint and subsequently the clients that rely on it
	// have heights 1 above the IAVL tree. Thus we return proof height + 1
	return proof, clienttypes.NewHeight(revision, uint64(res.Height)+1)
}

// QueryUpgradeProof performs an abci query with the given key and returns the proto encoded merkle proof
// for the query and the height at which the proof will succeed on a tendermint verifier.
func (chain *TestChain) QueryUpgradeProof(key []byte, height uint64) ([]byte, clienttypes.Height) {
	res := chain.App.Query(abci.RequestQuery{
		Path:   "store/upgrade/key",
		Height: int64(height - 1),
		Data:   key,
		Prove:  true,
	})

	merkleProof, err := commitmenttypes.ConvertProofs(res.ProofOps)
	require.NoError(chain.T, err)

	proof, err := chain.App.AppCodec().Marshal(&merkleProof)
	require.NoError(chain.T, err)

	revision := clienttypes.ParseChainID(chain.ChainID)

	// proof height + 1 is returned as the proof created corresponds to the height the proof
	// was created in the IAVL tree. Tendermint and subsequently the clients that rely on it
	// have heights 1 above the IAVL tree. Thus we return proof height + 1
	return proof, clienttypes.NewHeight(revision, uint64(res.Height+1))
}

// QueryConsensusStateProof performs an abci query for a consensus state
// stored on the given clientID. The proof and consensusHeight are returned.
func (chain *TestChain) QueryConsensusStateProof(clientID string) ([]byte, clienttypes.Height) {
	clientState := chain.GetClientState(clientID)

	consensusHeight := clientState.GetLatestHeight().(clienttypes.Height)
	consensusKey := host.FullConsensusStateKey(clientID, consensusHeight)
	proofConsensus, _ := chain.QueryProof(consensusKey)

	return proofConsensus, consensusHeight
}

// NextBlock sets the last header to the current header and increments the current header to be
// at the next block height. It does not update the time as that is handled by the Coordinator.
// It will call Endblock and Commit and apply the validator set changes to the next validators
// of the next block being created. This follows the Tendermint protocol of applying valset changes
// returned on block `n` to the validators of block `n+2`.
// It calls BeginBlock with the new block created before returning.
func (chain *TestChain) NextBlock() {
	res := chain.App.EndBlock(abci.RequestEndBlock{Height: chain.CurrentHeader.Height})

	chain.App.Commit()

	// set the last header to the current header
	// use nil trusted fields
	chain.LastHeader = chain.CurrentTMClientHeader()

	// val set changes returned from previous block get applied to the next validators
	// of this block. See tendermint spec for details.
	chain.Vals = chain.NextVals
	chain.NextVals = ApplyValSetChanges(chain.T, chain.Vals, res.ValidatorUpdates)

	// increment the current header
	chain.CurrentHeader = tmproto.Header{
		ChainID: chain.ChainID,
		Height:  chain.App.LastBlockHeight() + 1,
		AppHash: chain.App.LastCommitID().Hash,
		// NOTE: the time is increased by the coordinator to maintain time synchrony amongst
		// chains.
		Time:               chain.CurrentHeader.Time,
		ValidatorsHash:     chain.Vals.Hash(),
		NextValidatorsHash: chain.NextVals.Hash(),
		ProposerAddress:    chain.CurrentHeader.ProposerAddress,
	}

	chain.App.BeginBlock(abci.RequestBeginBlock{Header: chain.CurrentHeader})
}

// sendMsgs delivers a transaction through the application without returning the result.
func (chain *TestChain) sendMsgs(msgs ...sdk.Msg) error {
	_, err := chain.SendMsgs(msgs...)
	return err
}

// SendMsgs delivers a transaction through the application. It updates the senders sequence
// number and updates the TestChain's headers. It returns the result and error if one
// occurred.
func (chain *TestChain) SendMsgs(msgs ...sdk.Msg) (*sdk.Result, error) {
	// ensure the chain has the latest time
	chain.Coordinator.UpdateTimeForChain(chain)

	_, r, err := furyapp.SignAndDeliver(
		chain.T,
		chain.TxConfig,
		chain.App.GetBaseApp(),
		chain.GetContext().BlockHeader(),
		msgs,
		chain.ChainID,
		[]uint64{chain.SenderAccount.GetAccountNumber()},
		[]uint64{chain.SenderAccount.GetSequence()},
		true, true, chain.SenderPrivKey,
	)
	if err != nil {
		return nil, err
	}

	// NextBlock calls app.Commit()
	chain.NextBlock()

	// increment sequence for successful transaction execution
	err = chain.SenderAccount.SetSequence(chain.SenderAccount.GetSequence() + 1)
	if err != nil {
		return nil, err
	}

	chain.Coordinator.IncrementTime()

	return r, nil
}

// GetClientState retrieves the client state for the provided clientID. The client is
// expected to exist otherwise testing will fail.
func (chain *TestChain) GetClientState(clientID string) exported.ClientState {
	clientState, found := chain.App.GetIBCKeeper().ClientKeeper.GetClientState(chain.GetContext(), clientID)
	require.True(chain.T, found)

	return clientState
}

// GetConsensusState retrieves the consensus state for the provided clientID and height.
// It will return a success boolean depending on if consensus state exists or not.
func (chain *TestChain) GetConsensusState(clientID string, height exported.Height) (exported.ConsensusState, bool) {
	return chain.App.GetIBCKeeper().ClientKeeper.GetClientConsensusState(chain.GetContext(), clientID, height)
}

// GetValsAtHeight will return the validator set of the chain at a given height. It will return
// a success boolean depending on if the validator set exists or not at that height.
func (chain *TestChain) GetValsAtHeight(height int64) (*tmtypes.ValidatorSet, bool) {
	histInfo, ok := chain.App.GetStakingKeeper().GetHistoricalInfo(chain.GetContext(), height)
	if !ok {
		return nil, false
	}

	valSet := stakingtypes.Validators(SdkValidatorsToValidators(histInfo.Valset))

	tmValidators, err := teststaking.ToTmValidators(valSet, sdk.DefaultPowerReduction)
	if err != nil {
		panic(err)
	}
	return tmtypes.NewValidatorSet(tmValidators), true
}

// GetAcknowledgement retrieves an acknowledgement for the provided packet. If the
// acknowledgement does not exist then testing will fail.
func (chain *TestChain) GetAcknowledgement(packet exported.PacketI) []byte {
	ack, found := chain.App.GetIBCKeeper().ChannelKeeper.GetPacketAcknowledgement(chain.GetContext(), packet.GetDestPort(), packet.GetDestChannel(), packet.GetSequence())
	require.True(chain.T, found)

	return ack
}

// GetPrefix returns the prefix for used by a chain in connection creation
func (chain *TestChain) GetPrefix() commitmenttypes.MerklePrefix {
	return commitmenttypes.NewMerklePrefix(chain.App.GetIBCKeeper().ConnectionKeeper.GetCommitmentPrefix().Bytes())
}

// ConstructUpdateTMClientHeader will construct a valid 07-tendermint Header to update the
// light client on the source chain.
func (chain *TestChain) ConstructUpdateTMClientHeader(counterparty *TestChain, clientID string) (*ibctmtypes.Header, error) {
	return chain.ConstructUpdateTMClientHeaderWithTrustedHeight(counterparty, clientID, clienttypes.ZeroHeight())
}

// ConstructUpdateTMClientHeader will construct a valid 07-tendermint Header to update the
// light client on the source chain.
func (chain *TestChain) ConstructUpdateTMClientHeaderWithTrustedHeight(counterparty *TestChain, clientID string, trustedHeight clienttypes.Height) (*ibctmtypes.Header, error) {
	header := counterparty.LastHeader
	// Relayer must query for LatestHeight on client to get TrustedHeight if the trusted height is not set
	if trustedHeight.IsZero() {
		trustedHeight = chain.GetClientState(clientID).GetLatestHeight().(clienttypes.Height)
	}
	var (
		tmTrustedVals *tmtypes.ValidatorSet
		ok            bool
	)
	// Once we get TrustedHeight from client, we must query the validators from the counterparty chain
	// If the LatestHeight == LastHeader.Height, then TrustedValidators are current validators
	// If LatestHeight < LastHeader.Height, we can query the historical validator set from HistoricalInfo
	if trustedHeight == counterparty.LastHeader.GetHeight() {
		tmTrustedVals = counterparty.Vals
	} else {
		// NOTE: We need to get validators from counterparty at height: trustedHeight+1
		// since the last trusted validators for a header at height h
		// is the NextValidators at h+1 committed to in header h by
		// NextValidatorsHash
		tmTrustedVals, ok = counterparty.GetValsAtHeight(int64(trustedHeight.RevisionHeight + 1))
		if !ok {
			return nil, sdkerrors.Wrapf(ibctmtypes.ErrInvalidHeaderHeight, "could not retrieve trusted validators at trustedHeight: %d", trustedHeight)
		}
	}
	// inject trusted fields into last header
	// for now assume revision number is 0
	header.TrustedHeight = trustedHeight

	trustedVals, err := tmTrustedVals.ToProto()
	if err != nil {
		return nil, err
	}
	header.TrustedValidators = trustedVals

	return header, nil
}

// ExpireClient fast forwards the chain's block time by the provided amount of time which will
// expire any clients with a trusting period less than or equal to this amount of time.
func (chain *TestChain) ExpireClient(amount time.Duration) {
	chain.Coordinator.IncrementTimeBy(amount)
}

// CurrentTMClientHeader creates a TM header using the current header parameters
// on the chain. The trusted fields in the header are set to nil.
func (chain *TestChain) CurrentTMClientHeader() *ibctmtypes.Header {
	return chain.CreateTMClientHeader(chain.ChainID, chain.CurrentHeader.Height, clienttypes.Height{}, chain.CurrentHeader.Time, chain.Vals, chain.NextVals, nil, chain.Signers)
}

// CreateTMClientHeader creates a TM header to update the TM client. Args are passed in to allow
// caller flexibility to use params that differ from the chain.
func (chain *TestChain) CreateTMClientHeader(chainID string, blockHeight int64, trustedHeight clienttypes.Height, timestamp time.Time, tmValSet, nextVals, tmTrustedVals *tmtypes.ValidatorSet, signers map[string]tmtypes.PrivValidator) *ibctmtypes.Header {
	var (
		valSet      *tmproto.ValidatorSet
		trustedVals *tmproto.ValidatorSet
	)
	require.NotNil(chain.T, tmValSet)

	vsetHash := tmValSet.Hash()
	nextValHash := nextVals.Hash()

	tmHeader := tmtypes.Header{
		Version:            tmprotoversion.Consensus{Block: tmversion.BlockProtocol, App: 2},
		ChainID:            chainID,
		Height:             blockHeight,
		Time:               timestamp,
		LastBlockID:        MakeBlockID(make([]byte, tmhash.Size), 10_000, make([]byte, tmhash.Size)),
		LastCommitHash:     chain.App.LastCommitID().Hash,
		DataHash:           tmhash.Sum([]byte("data_hash")),
		ValidatorsHash:     vsetHash,
		NextValidatorsHash: nextValHash,
		ConsensusHash:      tmhash.Sum([]byte("consensus_hash")),
		AppHash:            chain.CurrentHeader.AppHash,
		LastResultsHash:    tmhash.Sum([]byte("last_results_hash")),
		EvidenceHash:       tmhash.Sum([]byte("evidence_hash")),
		ProposerAddress:    tmValSet.Proposer.Address, //nolint:staticcheck
	}

	hhash := tmHeader.Hash()
	blockID := MakeBlockID(hhash, 3, tmhash.Sum([]byte("part_set")))
	voteSet := tmtypes.NewVoteSet(chainID, blockHeight, 1, tmproto.PrecommitType, tmValSet)

	// MakeCommit expects a signer array in the same order as the validator array.
	// Thus we iterate over the ordered validator set and construct a signer array
	// from the signer map in the same order.
	var signerArr []tmtypes.PrivValidator   //nolint:prealloc // using prealloc here would be needlessly complex
	for _, v := range tmValSet.Validators { //nolint:staticcheck // need to check for nil validator set
		signerArr = append(signerArr, signers[v.Address.String()])
	}

	commit, err := tmtypes.MakeCommit(blockID, blockHeight, 1, voteSet, signerArr, timestamp)
	require.NoError(chain.T, err)

	signedHeader := &tmproto.SignedHeader{
		Header: tmHeader.ToProto(),
		Commit: commit.ToProto(),
	}

	if tmValSet != nil { //nolint:staticcheck
		valSet, err = tmValSet.ToProto()
		require.NoError(chain.T, err)
	}

	if tmTrustedVals != nil {
		trustedVals, err = tmTrustedVals.ToProto()
		require.NoError(chain.T, err)
	}

	// The trusted fields may be nil. They may be filled before relaying messages to a client.
	// The relayer is responsible for querying client and injecting appropriate trusted fields.
	return &ibctmtypes.Header{
		SignedHeader:      signedHeader,
		ValidatorSet:      valSet,
		TrustedHeight:     trustedHeight,
		TrustedValidators: trustedVals,
	}
}

// MakeBlockID copied unimported test functions from tmtypes to use them here
func MakeBlockID(hash []byte, partSetSize uint32, partSetHash []byte) tmtypes.BlockID {
	return tmtypes.BlockID{
		Hash: hash,
		PartSetHeader: tmtypes.PartSetHeader{
			Total: partSetSize,
			Hash:  partSetHash,
		},
	}
}

// CreatePortCapability binds and claims a capability for the given portID if it does not
// already exist. This function will fail testing on any resulting error.
// NOTE: only creation of a capability for a transfer or mock port is supported
// Other applications must bind to the port in InitGenesis or modify this code.
func (chain *TestChain) CreatePortCapability(scopedKeeper capabilitykeeper.ScopedKeeper, portID string) {
	// check if the portId is already binded, if not bind it
	_, ok := chain.App.GetScopedIBCKeeper().GetCapability(chain.GetContext(), host.PortPath(portID))
	if !ok {
		// create capability using the IBC capability keeper
		cap, err := chain.App.GetScopedIBCKeeper().NewCapability(chain.GetContext(), host.PortPath(portID))
		require.NoError(chain.T, err)

		// claim capability using the scopedKeeper
		err = scopedKeeper.ClaimCapability(chain.GetContext(), cap, host.PortPath(portID))
		require.NoError(chain.T, err)
	}

	chain.NextBlock()
}

// GetPortCapability returns the port capability for the given portID. The capability must
// exist, otherwise testing will fail.
func (chain *TestChain) GetPortCapability(portID string) *capabilitytypes.Capability {
	cap, ok := chain.App.GetScopedIBCKeeper().GetCapability(chain.GetContext(), host.PortPath(portID))
	require.True(chain.T, ok)

	return cap
}

// CreateChannelCapability binds and claims a capability for the given portID and channelID
// if it does not already exist. This function will fail testing on any resulting error. The
// scoped keeper passed in will claim the new capability.
func (chain *TestChain) CreateChannelCapability(scopedKeeper capabilitykeeper.ScopedKeeper, portID, channelID string) {
	capName := host.ChannelCapabilityPath(portID, channelID)
	// check if the portId is already binded, if not bind it
	_, ok := chain.App.GetScopedIBCKeeper().GetCapability(chain.GetContext(), capName)
	if !ok {
		cap, err := chain.App.GetScopedIBCKeeper().NewCapability(chain.GetContext(), capName)
		require.NoError(chain.T, err)
		err = scopedKeeper.ClaimCapability(chain.GetContext(), cap, capName)
		require.NoError(chain.T, err)
	}

	chain.NextBlock()
}

// GetChannelCapability returns the channel capability for the given portID and channelID.
// The capability must exist, otherwise testing will fail.
func (chain *TestChain) GetChannelCapability(portID, channelID string) *capabilitytypes.Capability {
	cap, ok := chain.App.GetScopedIBCKeeper().GetCapability(chain.GetContext(), host.ChannelCapabilityPath(portID, channelID))
	require.True(chain.T, ok)

	return cap
}

// GetTimeoutHeight is a convenience function which returns a IBC packet timeout height
// to be used for testing. It returns the current IBC height + 100 blocks
func (chain *TestChain) GetTimeoutHeight() clienttypes.Height {
	return clienttypes.NewHeight(clienttypes.ParseChainID(chain.ChainID), uint64(chain.GetContext().BlockHeight())+100)
}
