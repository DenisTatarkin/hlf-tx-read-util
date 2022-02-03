package parser

import (
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
)

func TestGetTrxs(t *testing.T) {
	var _block1,_ = ioutil.ReadFile("mock/block1.pb")
	var block1 = &common.Block{}
	proto.Unmarshal(_block1, block1)

	var trxs,_ = getTrxs(block1, nil)

	require.Len(t, trxs, 1)
	require.Len(t, trxs[0].Actions, 1)
	require.NotEmpty(t, trxs[0].Actions[0].Header)
	require.NotNil(t, trxs[0].Actions[0].ChaincodeActionPayload.ChaincodeEndorsedAction.ProposalResponsePayload.ChaincodeAction)
	//require.NotNil(t, trxs[0].Actions[0].Endorsements[1])
	require.NotEmpty(t, trxs[0].Actions[0].Header)
	//require.NotEmpty(t, trxs[0].Actions[0].CCProposalPayload)
	//require.NotEmpty(t, trxs[0].Actions[0].CCResponsePayload)*/
}

func TestParseBlocks(t *testing.T) {
	var _block1,_ = ioutil.ReadFile("mock/block1.pb")
	var _block2,_ = ioutil.ReadFile("./mock/block2.pb")
	var _block3,_ = ioutil.ReadFile("./mock/block3.pb")

	var block1 = &common.Block{}
	var block2 = &common.Block{}
	var block3 = &common.Block{}
	proto.Unmarshal(_block1, block1)
	proto.Unmarshal(_block2, block2)
	proto.Unmarshal(_block3, block3)
	var blocks = []common.Block{*block1,*block2,*block3}

	var dtos,_ = parseBlocks(blocks, nil)

	require.Len(t, dtos, 3)
	require.NotNil(t, dtos[0])
	require.NotNil(t, dtos[1])
	require.NotNil(t, dtos[2])
}

type testReaderImpl struct{

}

func (reader *testReaderImpl) GetBlocks(channel string) ([]common.Block, error){
	var _block1,_ = ioutil.ReadFile("mock/block1.pb")
	var _block2,_ = ioutil.ReadFile("./mock/block2.pb")
	var _block3,_ = ioutil.ReadFile("./mock/block3.pb")

	var block1 = &common.Block{}
	var block2 = &common.Block{}
	var block3 = &common.Block{}
	proto.Unmarshal(_block1, block1)
	proto.Unmarshal(_block2, block2)
	proto.Unmarshal(_block3, block3)
	return []common.Block{*block1,*block2,*block3}, nil
}

func TestGetTransactionsJSON(t *testing.T) {
	var reader = &testReaderImpl{}

	var json,_ = GetTransactionsJSON(reader, "",nil)

	require.NotEmpty(t, json)
}