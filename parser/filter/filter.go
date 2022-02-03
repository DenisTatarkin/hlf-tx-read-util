package filter

import "tx-read-utility/parser"

type Filter interface {
	Filter(input []parser.TrxDTO) (output []parser.TrxDTO)
}