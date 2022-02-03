package filter

import "tx-read-utility/parser"

type FilterImpl struct{
	trxId 	    string
	trxStatus   string
	isTrxId     bool
	isTrxStatus bool
}

func (f *FilterImpl) ByTrxId(trxId string){
	f.trxId = trxId
	f.isTrxId = true
}

func (f *FilterImpl) ByTrxStatus(trxStatus string){
	f.trxStatus = trxStatus
	f.isTrxStatus = true
}

func (f *FilterImpl) Filter(input []parser.TrxDTO) (output []parser.TrxDTO){
	if !f.isTrxId && !f.isTrxStatus{
		output = input
		return output
	}

	for _, trx := range input{
		ok := true

		if f.isTrxId && f.trxId != trx.Id{
			ok = false
		}

		if ok && f.isTrxStatus &&
			f.trxStatus != trx.Actions[0].ChaincodeActionPayload.ChaincodeEndorsedAction.ProposalResponsePayload.ChaincodeAction.Response.Status{
			ok = false
		}

		if ok{
			output = append(output, trx)
		}
	}

	return output
}