package parser

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-protos-go/common"
	"regexp"
	"tx-read-utility/parser/filter"
	"tx-read-utility/reader"
)

func GetTransactionsJSON(reader reader.BlockReader, channel string,  filter filter.Filter) ([]byte, error) {
	var blocks, err = reader.GetBlocks(channel)
	if err != nil{
		return nil, fmt.Errorf("error while reading blocks:\n%v", err)
	}

	trxs, err := parseBlocks(blocks, filter)
	if err != nil{
		return nil, fmt.Errorf("error while parsing blocks:\n%v", err)
	}

	data, err := json.MarshalIndent(trxs, "", "   ")
	if err!= nil{
		return nil, fmt.Errorf("error while marshalling transactions to json:\n%v", err)
	}

	re, _ := regexp.Compile(`\\u\d+`)
	data = re.ReplaceAll(data, []byte(" "))

	return data, nil
}

func parseBlocks(blocks []common.Block, filter filter.Filter) ([]TrxDTO, error) {
	var dtos []TrxDTO

	for _, block := range blocks{
		trxs, err := getTrxs(&block, filter)
		if err != nil{
			return nil, fmt.Errorf("error while getting transactions from block:\n%v", err)
		}

		dtos = append(dtos, trxs...)
	}

	return dtos, nil
}
