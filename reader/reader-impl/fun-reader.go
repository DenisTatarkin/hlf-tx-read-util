package reader_impl

import (
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	"io/ioutil"
)

type FunReaderImpl struct{

}

func (reader *FunReaderImpl) GetBlocks(channel string) ([]common.Block, error){
	//todo
	var _block1,_ = ioutil.ReadFile("parser/mock/block1.pb")
	var _block2,_ = ioutil.ReadFile("parser/mock/block2.pb")
	var _block3,_ = ioutil.ReadFile("parser/mock/block3.pb")

	var block1 = &common.Block{}
	var block2 = &common.Block{}
	var block3 = &common.Block{}
	proto.Unmarshal(_block1, block1)
	proto.Unmarshal(_block2, block2)
	proto.Unmarshal(_block3, block3)
	return []common.Block{*block1,*block2,*block3}, nil
}
