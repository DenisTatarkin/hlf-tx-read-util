package parser

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/ledger/rwset"
	"github.com/hyperledger/fabric-protos-go/ledger/rwset/kvrwset"
	"github.com/hyperledger/fabric-protos-go/peer"
	"strconv"
	"strings"
	"tx-read-utility/parser/filter"
	"unicode/utf8"
)

type TrxDTO struct {
	Actions []TransactionAction `json:"actions,omitempty"`
	Id string `json:"id,omitempty"`
}

type TransactionAction struct{
	Header string `json:"header,omitempty"`
	ChaincodeActionPayload ChaincodeActionPayload `json:"chaincode_action_payload,omitempty"`
}

type ChaincodeActionPayload struct{
	ChaincodeProposalPayload ChaincodeProposalPayload `json:"chaincode_proposal_payload,omitempty"`
	ChaincodeEndorsedAction ChaincodeEndorsedAction `json:"chaincode_endorsed_action,omitempty"`
}

type ChaincodeProposalPayload struct{
	ChaincodeInvocationSpec ChaincodeInvocationSpec `json:"chaincode_invocation_spec,omitempty"`
}

type ChaincodeInvocationSpec struct{
	Type string `json:"type,omitempty"`
	ChaincodeID string `json:"chaincode_id,omitempty"`
	ChaincodeInput ChaincodeInput `json:"chaincode_input,omitempty"`
}

type ChaincodeInput struct{
	Args []string `json:"args,omitempty"`
}

type ChaincodeEndorsedAction struct{
	Endorsements []Endorsement `json:"endorsements,omitempty"`
	ProposalResponsePayload ProposalResponsePayload `json:"proposal_response_payload,omitempty"`
}

type Endorsement struct {
	Endorser string `json:"endorser,omitempty"`
	Signature string `json:"signature,omitempty"`
}

type ProposalResponsePayload struct{
	ChaincodeAction ChaincodeAction `json:"chaincode_action,omitempty"`
}

type ChaincodeAction struct{
	ChaincodeEvent ChaincodeEvent `json:"chaincode_event,omitempty"`
	Response Response `json:"response,omitempty"`
	TxReadWriteSets []TxReadWriteSet `json:"tx_read_write_sets,omitempty"`
}

type ChaincodeEvent struct{
	TxID string `json:"tx_id,omitempty"`
	EventName string `json:"event_name,omitempty"`
	Payload string `json:"payload,omitempty"`
}

type Response struct{
	Status string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Payload string `json:"payload,omitempty"`
}

type TxReadWriteSet struct{
	Reads []string `json:"reads"`
	Writes []KVWrite `json:"writes"`
}

type KVWrite struct {
	Key string `json:"key"`
	Value string `json:"value"`
	IsDelete bool `json:"is_delete"`
}

func getTrxs (block *common.Block, filter filter.Filter) ([]TrxDTO, error){
	var envelopes []common.Envelope
	var dtos []TrxDTO
	for _, raw := range block.Data.Data{
		env := &common.Envelope{}
		err := proto.Unmarshal(raw, env)
		if err != nil{
			return nil, fmt.Errorf("error while envelope unmarshalling:\n%v", err)
		}

		envelopes = append(envelopes, *env)
	}
	if len(envelopes) == 0{
		return nil, errors.New("no envelopes in block")
	}

	for _, envelope := range envelopes{
		payl := &common.Payload{}
		err := proto.Unmarshal(envelope.Payload, payl)
		if err != nil{
			return nil, fmt.Errorf("error while payload unmarshalling:\n%v", err)
		}

		header := &common.ChannelHeader{}
		err = proto.Unmarshal(payl.Header.ChannelHeader, header)
		if err != nil{
			return nil, fmt.Errorf("error while channel header unmarshalling:\n%v", err)
		}

		trx := &peer.Transaction{}
		err = proto.Unmarshal(payl.Data, trx)
		if err != nil {
			return nil, fmt.Errorf("error while transaction unmarshalling:\n%v", err)
		}

		var trxActions []TransactionAction
		for _, trxAction := range trx.Actions{
			var transactionAction = &TransactionAction{}

			transactionAction.Header = formatString(trxAction.Header)
			

			ccActionPayload := &peer.ChaincodeActionPayload{}
			err := proto.Unmarshal(trxAction.Payload, ccActionPayload)
			if err != nil{
				return nil, fmt.Errorf("error while chaincode action payload unmarshalling:\n%v", err)
			}

			for _, endorsement := range ccActionPayload.Action.Endorsements{
				transactionAction.ChaincodeActionPayload.ChaincodeEndorsedAction.Endorsements = append(
					transactionAction.ChaincodeActionPayload.ChaincodeEndorsedAction.Endorsements, Endorsement{
						Endorser:  formatString (endorsement.Endorser),
						Signature: formatString(endorsement.Signature),
					})
			}


			var ccProposalResponsePayload = &peer.ProposalResponsePayload{}
			err = proto.Unmarshal(ccActionPayload.Action.ProposalResponsePayload, ccProposalResponsePayload)
			if err != nil{
				return nil, fmt.Errorf("error while chaincode proposal response payload unmarshalling:\n%v", err)
			}

			var ccAction = &peer.ChaincodeAction{}
			err = proto.Unmarshal(ccProposalResponsePayload.Extension, ccAction)
			if err != nil{
				return nil, fmt.Errorf("error while chaincode proposal response chaincode action unmarshalling:\n%v", err)
			}

			var event = &peer.ChaincodeEvent{}
			err = proto.Unmarshal(ccAction.Events, event)
			if err != nil{
				return nil, fmt.Errorf("error while chaincode event unmarshalling:\n%v", err)
			}

			transactionAction.ChaincodeActionPayload.ChaincodeEndorsedAction.
				ProposalResponsePayload.ChaincodeAction.ChaincodeEvent.Payload = formatString(event.Payload)
			transactionAction.ChaincodeActionPayload.ChaincodeEndorsedAction.
				ProposalResponsePayload.ChaincodeAction.ChaincodeEvent.TxID = event.TxId
			transactionAction.ChaincodeActionPayload.ChaincodeEndorsedAction.
				ProposalResponsePayload.ChaincodeAction.ChaincodeEvent.EventName = event.EventName

			transactionAction.ChaincodeActionPayload.ChaincodeEndorsedAction.
				ProposalResponsePayload.ChaincodeAction.Response.Payload = formatString(ccAction.Response.Payload)
			transactionAction.ChaincodeActionPayload.ChaincodeEndorsedAction.
				ProposalResponsePayload.ChaincodeAction.Response.Status = strconv.Itoa((int)(ccAction.Response.Status))
			transactionAction.ChaincodeActionPayload.ChaincodeEndorsedAction.
				ProposalResponsePayload.ChaincodeAction.Response.Message = ccAction.Response.Message

			var txRWSet = &rwset.TxReadWriteSet{}
			err = proto.Unmarshal(ccAction.Results, txRWSet)
			if err != nil{
				return nil, fmt.Errorf("error while rwset unmarshalling:\n%v", err)
			}

			for _, nsrwset := range txRWSet.NsRwset{
				kvrwset := &kvrwset.KVRWSet{}
				err = proto.Unmarshal(nsrwset.Rwset, kvrwset)
				if err != nil{
					return nil, fmt.Errorf("error while kvrwset unmarshalling:\n%v", err)
				}
				
				var reads []string
				for _, read := range kvrwset.Reads{
					reads = append(reads, read.Key)
				}
				
				var writes []KVWrite
				for _, write := range kvrwset.Writes{
					writes = append(writes, KVWrite{
						Key:      write.Key,
						Value:    formatString(write.Value),
						IsDelete: write.IsDelete,
					})
				}
				
				transactionAction.ChaincodeActionPayload.ChaincodeEndorsedAction.
					ProposalResponsePayload.ChaincodeAction.TxReadWriteSets = append(
						transactionAction.ChaincodeActionPayload.ChaincodeEndorsedAction.
						ProposalResponsePayload.ChaincodeAction.TxReadWriteSets, TxReadWriteSet{
						Reads:  reads,
						Writes: writes,
					})
			}

			var ccProposalPayload = &peer.ChaincodeProposalPayload{}
			err = proto.Unmarshal(ccActionPayload.ChaincodeProposalPayload, ccProposalPayload)
			if err != nil{
				return nil, fmt.Errorf("error while chaincode proposal payload unmarshalling:\n%v", err)
			}
			
			var ccInput = &peer.ChaincodeInput{}
			err = proto.Unmarshal(ccProposalPayload.Input, ccInput)
			if err != nil{
				return nil, fmt.Errorf("error while cc input unmarshalling:\n%v", err)
			}
			
			for _, arg := range ccInput.Args{
				transactionAction.ChaincodeActionPayload.ChaincodeProposalPayload.
					ChaincodeInvocationSpec.ChaincodeInput.Args = append (
						transactionAction.ChaincodeActionPayload.ChaincodeProposalPayload.ChaincodeInvocationSpec.ChaincodeInput.Args, 
						formatString(arg))
			}
			
			trxActions = append(trxActions, *transactionAction)
		}

		dtos = append(dtos, TrxDTO{
			Actions: trxActions,
			Id: header.TxId,
		})
	}

	if filter != nil{
		dtos = filter.Filter(dtos)
	}

	return dtos, nil
}

func formatString(input []byte) string{
	var output = bytes.Map(notUTF8, input)

	strings.ReplaceAll("\n", string(output), string(output))
	return string(output)
}

func notUTF8(r rune) rune{
	if r == utf8.RuneError || r == '\n' || r == '\t' || r == '"'{
		return -1
	}

	return r
}