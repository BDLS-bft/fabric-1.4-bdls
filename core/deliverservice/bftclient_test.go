/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package deliverservice

import (
	"crypto/x509"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/orderer"
	"github.com/hyperledger/fabric/common/flogging"
	"github.com/hyperledger/fabric/core/deliverservice/mocks"
	"github.com/hyperledger/fabric/gossip/util"
	"github.com/hyperledger/fabric/internal/pkg/comm"
	"github.com/hyperledger/fabric/internal/pkg/peer/blocksprovider/fake"
	"github.com/hyperledger/fabric/internal/pkg/peer/orderers"
	mocks2 "github.com/hyperledger/fabric/orderer/consensus/smartbft/mocks"
	"github.com/hyperledger/fabric/protoutil"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

var endpoints = []*orderers.Endpoint{
	{Address: "localhost:5611"},
	{Address: "localhost:5612"},
	{Address: "localhost:5613"},
	{Address: "localhost:5614"},
}

var endpoints3 = []*orderers.Endpoint{
	{Address: "localhost:5615"},
	{Address: "localhost:5616"},
	{Address: "localhost:5617"},
}

var endpointMap = map[int]*orderers.Endpoint{
	5611: {Address: "localhost:5611"},
	5612: {Address: "localhost:5612"},
	5613: {Address: "localhost:5613"},
	5614: {Address: "localhost:5614"},
}

var endpointMap3 = map[int]*orderers.Endpoint{
	5615: {Address: "localhost:5615"},
	5616: {Address: "localhost:5616"},
	5617: {Address: "localhost:5617"},
}

const goRoutineTestWaitTimeout = time.Second * 15

// Scenario: create an deliver adapter.
func TestNewBftDeliverAdapter(t *testing.T) {
	fakeLedgerInfo := &fake.LedgerInfo{}
	fakeBlockVerifier := &fake.BlockVerifier{}

	grpcClient, err := comm.NewGRPCClient(comm.ClientConfig{
		SecOpts: comm.SecureOptions{
			UseTLS: true,
		},
	})
	require.NoError(t, err)

	logger = flogging.NewFabricLogger(zap.NewNop())
	conf := &Config{
		CryptoSvc:         fakeBlockVerifier,
		Gossip:            nil,
		OrdererSource:     orderers.NewConnectionSource(logger, nil),
		Signer:            nil,
		DeliverGRPCClient: grpcClient,
	}

	da := NewBftDeliverAdapter("test-chain", fakeLedgerInfo, conf, &fake.Dialer{})
	require.NotNil(t, da)
}

// Scenario: create an delivery client.
func TestNewBFTDeliveryClient(t *testing.T) {
	flogging.ActivateSpec("bftDeliveryClient=DEBUG")

	grpcClient, err := comm.NewGRPCClient(comm.ClientConfig{
		SecOpts: comm.SecureOptions{
			UseTLS: true,
		},
	})
	require.NoError(t, err)

	fakeLedgerInfo := &fake.LedgerInfo{}
	fakeBlockVerifier := &fake.BlockVerifier{}
	mockSignerSerializer := &mocks2.SignerSerializer{}
	fakeDialer := &fake.Dialer{}
	fakeDialer.DialStub = func(string, *x509.CertPool) (*grpc.ClientConn, error) {
		cc, err := grpc.Dial("", grpc.WithInsecure())
		require.Nil(t, err)
		require.NotEqual(t, cc.GetState(), connectivity.Shutdown)
		return cc, nil
	}

	conn, err := fakeDialer.Dial("", x509.NewCertPool())
	require.Nil(t, err)
	require.NotNil(t, conn)

	bc, err := NewBFTDeliveryClient( "test-chain", endpoints, fakeLedgerInfo, fakeBlockVerifier, mockSignerSerializer, grpcClient, fakeDialer)
	require.NotNil(t, bc)
	require.Nil(t, err)
}

// Scenario: create a client against a set of orderer mocks. Receive several blocks and check block & header reception.
func Test_bftDeliveryClient_Recv(t *testing.T) {
	flogging.ActivateSpec("bftDeliveryClient=DEBUG")

	osArray := make([]*mocks.Orderer, 0)
	for port := range endpointMap {
		osArray = append(osArray, mocks.NewOrderer(port, t))
	}
	for _, os := range osArray {
		os.SetNextExpectedSeek(5)
	}

	grpcClient, err := comm.NewGRPCClient(comm.ClientConfig{
		SecOpts: comm.SecureOptions{
			UseTLS: true,
		},
	})
	require.NoError(t, err)

	fakeLedgerInfo := &fake.LedgerInfo{}
	fakeLedgerInfo.LedgerHeightReturns(5, nil)
	fakeBlockVerifier := &fake.BlockVerifier{}
	mockSignerSerializer := &mocks2.SignerSerializer{}
	mockSignerSerializer.On("Sign", mock.Anything).Return([]byte{1, 2, 3}, nil)
	mockSignerSerializer.On("Serialize", mock.Anything).Return([]byte{0, 2, 4, 6}, nil)
	fakeDialer := &fake.Dialer{}
	fakeDialer.DialCalls(func(ep string, cp *x509.CertPool) (*grpc.ClientConn, error) {
		cc, err := grpc.Dial(ep, grpc.WithInsecure(), grpc.WithBlock())
		require.Nil(t, err)
		require.NotEqual(t, cc.GetState(), connectivity.Shutdown)
		return cc, nil
	})

	bc, err := NewBFTDeliveryClient( "test-chain", endpoints, fakeLedgerInfo, fakeBlockVerifier, mockSignerSerializer, grpcClient, fakeDialer)
	require.NotNil(t, bc)
	require.Nil(t, err)

	go func() {
		for {
			resp, err := bc.Recv()
			if err != nil {
				assert.EqualError(t, err, errClientClosing.Error())
				return
			}
			block := resp.GetBlock()
			assert.NotNil(t, block)
			if block == nil {
				return
			}
			bc.UpdateReceived(block.Header.Number)
		}
	}()

	// all orderers send something: block/header
	beforeSend := time.Now()
	for seq := uint64(5); seq < uint64(10); seq++ {
		for _, os := range osArray {
			os.SendBlock(seq)
		}
	}

	time.Sleep(time.Second)
	bc.Close()

	lastN, lastT := bc.GetNextBlockNumTime()
	assert.Equal(t, uint64(10), lastN)
	assert.True(t, lastT.After(beforeSend))

	headersNum, headerT, headerErr := bc.GetHeadersBlockNumTime()
	for i, n := range headersNum {
		assert.Equal(t, uint64(10), n)
		assert.True(t, headerT[i].After(beforeSend))
		assert.NoError(t, headerErr[i])
	}

	for _, os := range osArray {
		os.Shutdown()
	}
}

// Scenario: block censorship by orderer. Create a client against a set of orderer mocks.
// Receive one block. Then, the orderer sending blocks stops sending but headers keep coming.
// The client should switch to another orderer and seek from the new height. Check block & header reception.
func TestBFTDeliverClient_Censorship(t *testing.T) {
	flogging.ActivateSpec("bftDeliveryClient,deliveryClient=DEBUG")
	viper.Set("peer.deliveryclient.bft.blockCensorshipTimeout", 2*time.Second)
	defer viper.Reset()

	assert.Equal(t, util.GetDurationOrDefault("peer.deliveryclient.bft.blockCensorshipTimeout", bftBlockCensorshipTimeout), 2*time.Second)

	osMap := make(map[string]*mocks.Orderer, len(endpointMap))
	for port, ep := range endpointMap {
		osMap[ep.Address] = mocks.NewOrderer(port, t)
	}
	for _, os := range osMap {
		os.SetNextExpectedSeek(5)
	}

	grpcClient, err := comm.NewGRPCClient(comm.ClientConfig{
		SecOpts: comm.SecureOptions{
			UseTLS: true,
		},
	})
	require.NoError(t, err)

	var ledgerHeight uint64 = 5

	fakeLedgerInfo := &fake.LedgerInfo{}
	fakeLedgerInfo.LedgerHeightStub = func() (uint64, error) {
		return atomic.LoadUint64(&ledgerHeight), nil
	}
	fakeBlockVerifier := &fake.BlockVerifier{}
	fakeBlockVerifier.VerifyHeaderReturns(nil)
	mockSignerSerializer := &mocks2.SignerSerializer{}
	mockSignerSerializer.On("Sign", mock.Anything).Return([]byte{1, 2, 3}, nil)
	mockSignerSerializer.On("Serialize", mock.Anything).Return([]byte{0, 2, 4, 6}, nil)
	fakeDialer := &fake.Dialer{}
	fakeDialer.DialCalls(func(ep string, cp *x509.CertPool) (*grpc.ClientConn, error) {
		cc, err := grpc.Dial(ep, grpc.WithInsecure(), grpc.WithBlock())
		require.Nil(t, err)
		require.NotEqual(t, cc.GetState(), connectivity.Shutdown)
		return cc, nil
	})

	bc, err := NewBFTDeliveryClient( "test-chain", endpoints, fakeLedgerInfo, fakeBlockVerifier, mockSignerSerializer, grpcClient, fakeDialer)
	require.NotNil(t, bc)
	require.Nil(t, err)

	go func() {
		for {
			resp, err := bc.Recv()
			if err != nil {
				assert.EqualError(t, err, "client is closing")
				return
			}
			block := resp.GetBlock()
			assert.NotNil(t, block)
			if block == nil {
				return
			}
			atomic.StoreUint64(&ledgerHeight, block.Header.Number+1)
			bc.UpdateReceived(block.Header.Number)
		}
	}()

	blockEP, err := waitForBlockEP(bc)
	assert.NoError(t, err)
	require.NotNil(t, blockEP)
	osnMocks, err := detectOSNConnections(true, osnMapValues(osMap)...)
	assert.NoError(t, err)
	assert.Equal(t, strings.Split(blockEP.Address, ":")[1], strings.Split(osnMocks[0].Addr().String(), ":")[1])

	// one normal block
	beforeSend := time.Now()
	for _, os := range osMap {
		os.SendBlock(5)
	}
	for _, os := range osMap {
		os.SetNextExpectedSeek(6)
	}
	time.Sleep(time.Second)
	lastN, lastT := bc.GetNextBlockNumTime()
	assert.Equal(t, uint64(6), lastN)
	assert.True(t, lastT.After(beforeSend))
	verifyHeaderReceivers(t, bc, 3, 5, beforeSend, 0)

	// only headers
	beforeSend = time.Now()
	for seq := uint64(6); seq < uint64(10); seq++ {
		for _, os := range osMap {
			if strings.Split(blockEP.Address, ":")[1] == strings.Split(os.Addr().String(), ":")[1] { // censorship
				continue
			}
			os.SendBlock(seq)
		}
	}

	time.Sleep(3 * time.Second)

	// the client detected the censorship and switched
	blockEP2, err := waitForBlockEP(bc)
	assert.NoError(t, err)
	require.NotNil(t, blockEP2)
	assert.True(t, blockEP.Address != blockEP2.Address)

	for seq := uint64(6); seq < uint64(10); seq++ {
		for _, os := range osMap {
			if strings.Split(blockEP2.Address, ":")[1] == strings.Split(os.Addr().String(), ":")[1] {
				os.SendBlock(seq)
			}
		}
	}

	time.Sleep(1 * time.Second)

	lastN, lastT = bc.GetNextBlockNumTime()
	assert.Equal(t, uint64(10), lastN)
	assert.True(t, lastT.After(beforeSend))
	verifyHeaderReceivers(t, bc, 3, 9, beforeSend, 1)

	bc.Close()

	for _, os := range osMap {
		os.Shutdown()
	}
}

// Scenario: server fail-over. Create a client against a set of orderer mocks.
// Receive one block. Then, the orderer sending blocks fails.
// The client should switch to another orderer and seek from the new height. Send a few blocks.
// Check block & header reception.
func TestBFTDeliverClient_Failover(t *testing.T) {
	flogging.ActivateSpec("bftDeliveryClient,deliveryClient=DEBUG")
	viper.Set("peer.deliveryclient.bft.blockRcvTotalBackoffDelay", time.Second)
	viper.Set("peer.deliveryclient.bft.maxBackoffDelay", time.Second)
	viper.Set("peer.deliveryclient.connTimeout", 100*time.Millisecond)
	defer viper.Reset()

	osArray := make([]*mocks.Orderer, 0)
	for port := range endpointMap {
		osArray = append(osArray, mocks.NewOrderer(port, t))
	}
	for _, os := range osArray {
		os.SetNextExpectedSeek(5)
	}

	grpcClient, err := comm.NewGRPCClient(comm.ClientConfig{
		SecOpts: comm.SecureOptions{
			UseTLS: true,
		},
	})
	require.NoError(t, err)

	var ledgerHeight uint64 = 5

	fakeLedgerInfo := &fake.LedgerInfo{}
	fakeLedgerInfo.LedgerHeightStub = func() (uint64, error) {
		return atomic.LoadUint64(&ledgerHeight), nil
	}
	fakeBlockVerifier := &fake.BlockVerifier{}
	fakeBlockVerifier.VerifyHeaderReturns(nil)
	mockSignerSerializer := &mocks2.SignerSerializer{}
	mockSignerSerializer.On("Sign", mock.Anything).Return([]byte{1, 2, 3}, nil)
	mockSignerSerializer.On("Serialize", mock.Anything).Return([]byte{0, 2, 4, 6}, nil)
	fakeDialer := &fake.Dialer{}
	fakeDialer.DialCalls(func(ep string, cp *x509.CertPool) (*grpc.ClientConn, error) {
		cc, err := grpc.Dial(ep, grpc.WithInsecure(), grpc.WithBlock())
		require.Nil(t, err)
		require.NotEqual(t, cc.GetState(), connectivity.Shutdown)
		return cc, nil
	})

	bc, err := NewBFTDeliveryClient( "test-chain", endpoints, fakeLedgerInfo, fakeBlockVerifier, mockSignerSerializer, grpcClient, fakeDialer)
	require.NotNil(t, bc)
	require.Nil(t, err)

	go func() {
		for {
			resp, err := bc.Recv()
			if err != nil {
				assert.EqualError(t, err, "client is closing")
				return
			}
			block := resp.GetBlock()
			assert.NotNil(t, block)
			if block == nil {
				return
			}
			atomic.StoreUint64(&ledgerHeight, block.Header.Number+1)
			bc.UpdateReceived(block.Header.Number)
		}
	}()

	time.Sleep(time.Second)
	blockEP, err := waitForBlockEP(bc)
	assert.NoError(t, err)

	osnMocks, err := detectOSNConnections(true, osArray...)
	assert.NoError(t, err)
	assert.Equal(t, strings.Split(blockEP.Address, ":")[1], strings.Split(osnMocks[0].Addr().String(), ":")[1])

	// one normal block
	beforeSend := time.Now()
	for _, os := range osArray {
		os.SendBlock(5)
	}
	for _, os := range osArray {
		os.SetNextExpectedSeek(6)
	}
	time.Sleep(time.Second)

	for _, os := range osArray {
		if strings.Split(blockEP.Address, ":")[1] == strings.Split(os.Addr().String(), ":")[1] {
			os.Shutdown()
			bftLogger.Infof("TEST: shutting down: %s", os.Addr().String())
		}
	}

	time.Sleep(10 * time.Second)

	// the client detected the failure and switched
	blockEP2, err := waitForBlockEP(bc)
	assert.NoError(t, err)
	assert.True(t, blockEP != blockEP2)

	beforeSend = time.Now()
	for seq := uint64(6); seq < uint64(10); seq++ {
		for _, os := range osArray {
			if strings.Split(blockEP.Address, ":")[1] == strings.Split(os.Addr().String(), ":")[1] {
				continue
			}
			os.SendBlock(seq)
		}
	}

	time.Sleep(2 * time.Second)

	lastN, lastT := bc.GetNextBlockNumTime()
	assert.Equal(t, uint64(10), lastN)
	assert.True(t, lastT.After(beforeSend))
	verifyHeaderReceivers(t, bc, 3, 9, beforeSend, 1)

	//restart the orderer
	//for port, ep := range endpointMap {
	//	if ep == blockEP {
	//		os := mocks.NewOrderer(port, t)
	//		os.SetNextExpectedSeek(10)
	//		osArray[ep.Endpoint] = os
	//		bftLogger.Infof("TEST: restarting: %s", ep)
	//	}
	//}

	//time.Sleep(2 * time.Second)
	//beforeSend = time.Now()
	//for _, os := range osArray {
	//	os.SendBlock(10)
	//}
	//time.Sleep(1 * time.Second)
	//
	//lastN, lastT = bc.GetNextBlockNumTime()
	//assert.Equal(t, uint64(11), lastN)
	//assert.True(t, lastT.After(beforeSend))
	//verifyHeaderReceivers(t, bc, 3, 10, beforeSend, 0)

	bc.Close()

	for _, os := range osArray {
		os.Shutdown()
	}
}

func waitForBlockEP(bc *bftDeliveryClient) (*orderers.Endpoint, error) {
	for i := int64(0); ; i++ {

		blockEP := bc.GetEndpoint()
		if blockEP != nil {
			return blockEP, nil
		}
		time.Sleep(10 * time.Millisecond)
		if time.Duration(i*10*time.Millisecond.Nanoseconds()) > 5*time.Second {
			return nil, errors.New("timeout: no block receiver")
		}
	}
}

func detectOSNConnections(isBFT bool, osSet ...*mocks.Orderer) ([]*mocks.Orderer, error) {
	var first int
	r := make([]*mocks.Orderer, len(osSet))

	if !isBFT {
		for i := int64(0); ; i++ {
			first = -1
			for j, os := range osSet {
				if os.ConnCount() > 0 {
					first = j
					break
				}
			}

			if first >= 0 {
				break //found
			}

			time.Sleep(time.Millisecond * 10)
			if time.Duration(i*10*time.Millisecond.Nanoseconds()) > 5*time.Second {
				return nil, errors.New("timeout: no connections")
			}
		}
	} else {
		cType := make([]orderer.SeekInfo_SeekContentType, len(osSet))
		for i := int64(0); ; i++ {
			var cCnt int
			first = -1
			for j, os := range osSet {
				c1, t1 := os.ConnCountType()
				if c1 > 0 {
					cCnt++
				}
				cType[j] = t1
				if t1 == orderer.SeekInfo_BLOCK {
					first = j
				}
			}

			if cCnt == len(osSet) && first >= 0 {
				break //found
			}

			time.Sleep(time.Millisecond * 10)
			if time.Duration(i*10*time.Millisecond.Nanoseconds()) > 5*time.Second {
				return nil, errors.New("timeout: no connections")
			}
		}
	}

	r[0] = osSet[first]
	n := 1
	for j, os := range osSet {
		if j == first {
			continue
		}
		r[n] = os
		n++
	}
	return r, nil
}

func verifyHeaderReceivers(t *testing.T, client *bftDeliveryClient, expectedNumHdr int, expectedBlockNum uint64, expectedBlockTime time.Time, expectedNumErr int) {
	lastHdrN, lastHdrT, lastHdrErr := client.GetHeadersBlockNumTime()
	numErr := 0
	assert.Len(t, lastHdrN, expectedNumHdr)
	for i, n := range lastHdrN {
		if lastHdrErr[i] == nil {
			assert.Equal(t, uint64(expectedBlockNum), n)
			assert.True(t, lastHdrT[i].After(expectedBlockTime))
		} else {
			numErr++
		}
	}
	assert.Equal(t, expectedNumErr, numErr)
	return
}

type mockOrderer struct {
	t *testing.T
}

func (*mockOrderer) Broadcast(orderer.AtomicBroadcast_BroadcastServer) error {
	panic("not implemented")
}

func (o *mockOrderer) Deliver(stream orderer.AtomicBroadcast_DeliverServer) error {
	envelope, _ := stream.Recv()
	inspectTLSBinding := comm.NewBindingInspector(true, func(msg proto.Message) []byte {
		env, isEnvelope := msg.(*common.Envelope)
		if !isEnvelope || env == nil {
			assert.Fail(o.t, "not an envelope")
		}
		ch, err := protoutil.ChannelHeader(env)
		assert.NoError(o.t, err)
		return ch.TlsCertHash
	})

	err := inspectTLSBinding(stream.Context(), envelope)
	assert.NoError(o.t, err, "orderer rejected TLS binding")

	payload, err := protoutil.UnmarshalPayload(envelope.Payload)
	assert.NoError(o.t, err, "cannot unmarshal payload")
	seekInfo := &orderer.SeekInfo{}
	err = proto.Unmarshal(payload.Data, seekInfo)
	assert.NoError(o.t, err, "cannot unmarshal payload.Data")

	block := &common.Block{
		Header: &common.BlockHeader{
			Number: 100,
		},
		Data: &common.BlockData{
			Data: [][]byte{{1, 2}, {3, 4}},
		},
		Metadata: &common.BlockMetadata{
			Metadata: [][]byte{{5, 6}, {7, 8}},
		},
	}

	if oldest := seekInfo.Start.GetOldest(); oldest != nil {
		block.Header.Number = 0
	}

	if seekInfo.ContentType == orderer.SeekInfo_HEADER_WITH_SIG {
		block.Data = nil
	}

	response := &orderer.DeliverResponse{
		Type: &orderer.DeliverResponse_Block{
			Block: block,
		},
	}

	stream.Send(response)
	return nil
}

func osnMapValues(m map[string]*mocks.Orderer) []*mocks.Orderer {
	s := make([]*mocks.Orderer, 0)
	for _, v := range m {
		s = append(s, v)
	}
	return s
}