/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package smartbft

import (
	"github.com/hyperledger/fabric/bccsp"
	"sync/atomic"

	smartbft "github.com/SmartBFT-Go/consensus/pkg/consensus"
	"github.com/SmartBFT-Go/consensus/pkg/types"
	cb "github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric/common/flogging"
	"github.com/hyperledger/fabric/common/policies"
	"github.com/hyperledger/fabric/orderer/common/cluster"
	"github.com/hyperledger/fabric/orderer/consensus"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

//go:generate counterfeiter -o mocks/mock_blockpuller.go . BlockPuller

// BlockPuller is used to pull blocks from other OSN
type BlockPuller interface {
	PullBlock(seq uint64) *cb.Block
	HeightsByEndpoints() (map[string]uint64, error)
	Close()
}

// WALConfig consensus specific configuration parameters from orderer.yaml; for SmartBFT only WALDir is relevant.
type WALConfig struct {
	WALDir            string // WAL data of <my-channel> is stored in WALDir/<my-channel>
	SnapDir           string // Snapshots of <my-channel> are stored in SnapDir/<my-channel>
	EvictionSuspicion string // Duration threshold that the node samples in order to suspect its eviction from the channel.
}

type ConfigValidator interface {
	ValidateConfig(env *cb.Envelope) error
}

type signerSerializer interface {
	// Sign a message and return the signature over the digest, or error on failure
	Sign(message []byte) ([]byte, error)

	// Serialize converts an identity to bytes
	Serialize() ([]byte, error)
}

// BFTChain implements Chain interface to wire with
// BFT smart library
type BFTChain struct {
	RuntimeConfig    *atomic.Value
	Channel          string
	Config           types.Configuration
	BlockPuller      BlockPuller
	Comm             cluster.Communicator
	SignerSerializer signerSerializer
	PolicyManager    policies.Manager
	Logger           *flogging.FabricLogger
	WALDir           string
	consensus        *smartbft.Consensus
	support          consensus.ConsenterSupport
	verifier         *Verifier
	assembler        *Assembler
	Metrics          *Metrics
	bccsp            bccsp.BCCSP
}

// NewChain creates new BFT Smart chain
func NewChain(
	cv ConfigValidator,
	selfID uint64,
	config types.Configuration,
	walDir string,
	blockPuller BlockPuller,
	comm cluster.Communicator,
	signerSerializer signerSerializer,
	policyManager policies.Manager,
	support consensus.ConsenterSupport,
	metrics *Metrics,
	bccsp bccsp.BCCSP,

) (*BFTChain, error) {

	//requestInspector := &RequestInspector{
	//	ValidateIdentityStructure: func(_ *msp.SerializedIdentity) error {
	//		return nil
	//	},
	//}

	logger := flogging.MustGetLogger("orderer.consensus.smartbft.chain").With(zap.String("channel", support.ChannelID()))

	c := &BFTChain{
		RuntimeConfig:    &atomic.Value{},
		Channel:          support.ChannelID(),
		Config:           config,
		WALDir:           walDir,
		Comm:             comm,
		support:          support,
		SignerSerializer: signerSerializer,
		PolicyManager:    policyManager,
		BlockPuller:      blockPuller,
		Logger:           logger,
		// todo: BFT metrics
		//Metrics: &Metrics{
		//	ClusterSize:          metrics.ClusterSize.With("channel", support.ChannelID()),
		//	CommittedBlockNumber: metrics.CommittedBlockNumber.With("channel", support.ChannelID()),
		//	IsLeader:             metrics.IsLeader.With("channel", support.ChannelID()),
		//	LeaderID:             metrics.LeaderID.With("channel", support.ChannelID()),
		//},
	}

	lastBlock := LastBlockFromLedgerOrPanic(support, c.Logger)
	lastConfigBlock := LastConfigBlockFromLedgerOrPanic(support, c.Logger)

	rtc := RuntimeConfig{
		logger: logger,
		id:     selfID,
	}
	rtc, err := rtc.BlockCommitted(lastConfigBlock, bccsp)
	if err != nil {
		return nil, errors.Wrap(err, "failed constructing RuntimeConfig")
	}
	rtc, err = rtc.BlockCommitted(lastBlock, bccsp)
	if err != nil {
		return nil, errors.Wrap(err, "failed constructing RuntimeConfig")
	}

	c.RuntimeConfig.Store(rtc)

	//c.verifier = buildVerifier(cv, c.RuntimeConfig, support, requestInspector, policyManager)
	//c.consensus = bftSmartConsensusBuild(c, requestInspector)

	// Setup communication with list of remotes notes for the new channel
	c.Comm.Configure(c.support.ChannelID(), rtc.RemoteNodes)

	if err := c.consensus.ValidateConfiguration(rtc.Nodes); err != nil {
		return nil, errors.Wrap(err, "failed to verify SmartBFT-Go configuration")
	}

	logger.Infof("SmartBFT-v3 is now servicing chain %s", support.ChannelID())

	return c, nil
}

func (B BFTChain) Order(env *cb.Envelope, configSeq uint64) error {
	panic("implement me")
}

func (B BFTChain) Configure(config *cb.Envelope, configSeq uint64) error {
	panic("implement me")
}

func (B BFTChain) WaitReady() error {
	panic("implement me")
}

func (B BFTChain) Errored() <-chan struct{} {
	panic("implement me")
}

func (B BFTChain) Start() {
	panic("implement me")
}

func (B BFTChain) Halt() {
	panic("implement me")
}
