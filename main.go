package main

import (
	"github.com/niktrix/bitcoin-lib/account"
	"github.com/niktrix/bitcoin-lib/utils"

	"log"

	btcchain "github.com/btcsuite/btcd/chaincfg"
)

func main() {
	var ut []utils.UTXO

	//compressedKey := "cPmTTa8ctUckw7KYppLd1Vkx7jxjjRcpMqZ6Dm4n2FVfXRBRyirL" // compressed key
	unCompressedKey := "5JsjKubviP3TDfNfbE3qdxKuNqqSVCctEF3jzyw26qYzonGEgsE" //uncompressed private key
	isCompressed := false
	chain := "testnet" // testnet || mainnet
	chainConfig := &btcchain.TestNet3Params

	switch chain {
	case "testnet":
		chainConfig = &btcchain.TestNet3Params
		break
	case "mainnet":
		chainConfig = &btcchain.MainNetParams
		break
	}
	destination := "<TO BTC ADDRESS>"
	amount := int64(32000000000)
	txFee := int64(1000)
	acc, err := account.NewAccount(unCompressedKey, chainConfig, isCompressed)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("fromAddress", acc.Address)
	log.Println("toAddress", destination)
	log.Println("txFee", txFee)
	log.Println("amount", amount)

	btchelper := utils.NewBitPay(chain)
	ut, err = btchelper.GetUnspentTxs(acc.Address.String())
	if err != nil {
		log.Fatalln("Error Getting unspent Tx", err)
		return
	}

	if len(ut) == 0 {
		log.Println(" No unspent Tx available")
		return
	}

	transaction := utils.NewTx()
	transaction.SetUnspentTxs(ut)
	transaction.SetAmount(amount)
	transaction.SetFee(txFee)
	transaction.SetFrom(acc)
	transaction.SetTo(destination)
	transaction.SetConfig(chainConfig)
	transaction.SetIsCompress(isCompressed)
	err = transaction.Execute()
	if err != nil {
		log.Fatalln("Error Executing Tx", err)
	}

	rawtx := transaction.GetRaw()

	log.Println("rawTx: ", rawtx)
	response, err := btchelper.BroadCastTX(rawtx)
	log.Println(err)

	log.Println("response: ", response)

}
