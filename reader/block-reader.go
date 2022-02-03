package reader

import (
	"github.com/hyperledger/fabric-protos-go/common"
)

type BlockReader interface {
	GetBlocks(channel string) ([]common.Block, error)
}
