package main

import (
	"flag"
	"log"
	"os"
	"tx-read-utility/parser"
	filter "tx-read-utility/parser/filter"
	reader_impl "tx-read-utility/reader/reader-impl"
)

func main(){
	var txId, txStatus, channel string

	if len(os.Args) < 2{
		log.Fatal("no required parameters:\n -ch=\"channel name\" is required\n-txid and -txstatus are not obligate")
	}

	flag.StringVar(&channel, "ch", "", "Channel name")
	flag.Parse()
	if !flag.Parsed(){
		log.Fatal("incorrect channel name")
	}

	if len(os.Args) > 2{
		cmd := flag.NewFlagSet("filter", flag.ExitOnError)
		cmd.StringVar(&txId, "txid", "", "Transaction Id")
		cmd.StringVar(&txStatus, "txstatus", "", "Transaction status")
	}

	var reader = &reader_impl.FunReaderImpl{}
	var filter = &filter.FilterImpl{}
	if txId != ""{
		filter.ByTrxId(txId)
	}
	if txStatus != ""{
		filter.ByTrxStatus(txStatus)
	}

	defer func() {
		if err := recover(); err != nil {
		log.Fatalf("recovered from panic \n%v", err)
	}}()

	var json, err = parser.GetTransactionsJSON(reader, channel, filter)
	if err != nil{
		log.Fatalf("error while getting transactions in json:\n%v", err)
	}

	file, err := os.Create("transactions.json")
	if err != nil{
		log.Fatalf("error while creating output file:\n%v", err)
	}
	defer func() {
		file.Close()
	}()

	_,err = file.Write(json)
	if err != nil{
		log.Fatalf("error while writing data in output file:\n%v\ncontent of json:\n%s",err, json)
	}
}
