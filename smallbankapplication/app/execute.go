package app

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
)

const (
	GetBalance    uint8 = 1
	Amalgamate    uint8 = 2
	UpdateBalance uint8 = 3
	UpdateSaving  uint8 = 4
	SendPayment   uint8 = 5
	WriteCheck    uint8 = 6
)

const (
	TxTypeString  string = "T"
	TxIdString    string = "I"
	FromString    string = "F"
	ToString      string = "O"
	BalanceString string = "B"
)

func ResolveTx(request *[]byte) []SmallBankTransaction {
	txs := bytes.Split(*request, []byte(">"))
	l := len(txs)
	if l == 0 {
		//log.Println("the tx is nil")
	}
	ReceiveTx := make([]SmallBankTransaction, l)
	for i, elements := range txs {
		var tx SmallBankTransaction
		element := bytes.Split(elements, []byte(","))
		for _, e := range element {
			kv := bytes.Split(e, []byte("="))
			switch string(kv[0]) {
			case TxTypeString:
				temp, _ := strconv.ParseUint(string(kv[1]), 10, 64)
				tx.T = uint8(temp)
			case TxIdString:
				temp, _ := strconv.ParseUint(string(kv[1]), 10, 64)
				tx.I = uint16(temp)
			case FromString:
				tx.F = make([]byte, len(kv[1]))
				copy(tx.F, kv[1])
			case ToString:
				tx.O = make([]byte, len(kv[1]))
				copy(tx.O, kv[1])
			case BalanceString:
				// temp_value := string(kv[1])
				tx.B = BytesToInt(kv[1])
			}
		}
		ReceiveTx[i] = tx
	}
	return ReceiveTx
}

func (BCstate *BlockchainState) ResolveAndExecuteTx(request *[]byte) {
	// resovle and execute tx one by one
	// T=3,I=1,F=1,O=3,B=156>T=1,I=2,F=2,O=1,B=190"
	txs := bytes.Split(*request, []byte(">"))
	l := len(txs)
	if l == 0 {
		log.Println("the tx is nil")
	}
	ReceiveTx := make([]SmallBankTransaction, l)
	for i, elements := range txs {
		var tx SmallBankTransaction
		element := bytes.Split(elements, []byte(","))
		for _, e := range element {
			kv := bytes.Split(e, []byte("="))
			switch string(kv[0]) {
			case TxTypeString:
				temp, _ := strconv.ParseUint(string(kv[1]), 10, 64)
				tx.T = uint8(temp)
			case TxIdString:
				temp, _ := strconv.ParseUint(string(kv[1]), 10, 64)
				tx.I = uint16(temp)
			case FromString:
				tx.F = make([]byte, len(kv[1]))
				copy(tx.F, kv[1])
			case ToString:
				tx.O = make([]byte, len(kv[1]))
				copy(tx.O, kv[1])
			case BalanceString:
				// temp_value := string(kv[1])
				tx.B = BytesToInt(kv[1])
			}
		}
		ReceiveTx[i] = tx
	}

	var TxType uint8
	var From []byte
	var To []byte
	var Balance int

	l = len(ReceiveTx)
	for i := 0; i < l; i++ {
		tx := ReceiveTx[i]
		TxType = tx.T
		From = tx.F
		To = tx.O
		Balance = tx.B

		switch TxType {
		case GetBalance:
			BCstate.GetBalance(string(From))
		case Amalgamate:
			BCstate.Amalgamate(string(From), string(To))
		case UpdateBalance:
			BCstate.UpdateBalance(string(From), Balance)
		case UpdateSaving:
			BCstate.UpdateSaving(string(From), Balance)
		case SendPayment:
			BCstate.SendPayment(string(From), string(To), Balance)
		case WriteCheck:
			BCstate.WriteCheck(string(From), Balance)
		default:
			fmt.Println("T doesn't match")
		}
	}
}

// Original version on resolving Tx in batch
//func (BCstate *BlockchainState) ResolveAndExecuteTx(request *[]byte) {
//	// T=3,I=1,F=1,O=3,B=156>T=1,I=2,F=2,O=1,B=190"
//	txs := bytes.Split(*request, []byte(">"))
//	l := len(txs)
//	if l == 0 {
//		log.Println("the tx is nil")
//	}
//	ReceiveTx := make([]SmallBankTransaction, l)
//	/*err := json.Unmarshal(*request, &ReceiveTx)
//	if err != nil {
//		log.Println(err)
//	}*/
//	for i, elements := range txs {
//		var tx SmallBankTransaction
//		element := bytes.Split(elements, []byte(","))
//		for _, e := range element {
//			kv := bytes.Split(e, []byte("="))
//			switch string(kv[0]) {
//			case TxTypeString:
//				temp, _ := strconv.ParseUint(string(kv[1]), 10, 64)
//				tx.T = uint8(temp)
//			case TxIdString:
//				temp, _ := strconv.ParseUint(string(kv[1]), 10, 64)
//				tx.I = uint16(temp)
//			case FromString:
//				tx.F = make([]byte, len(kv[1]))
//				copy(tx.F, kv[1])
//			case ToString:
//				tx.O = make([]byte, len(kv[1]))
//				copy(tx.O, kv[1])
//			case BalanceString:
//				// temp_value := string(kv[1])
//				tx.B = BytesToInt(kv[1])
//			}
//		}
//		ReceiveTx[i] = tx
//	}
//
//	var TxType uint8
//	var From []byte
//	var To []byte
//	var Balance int
//
//	for i := 0; i < l; i++ {
//		tx := ReceiveTx[i]
//		TxType = tx.T
//		From = tx.F
//		To = tx.O
//		Balance = tx.B
//
//		switch TxType {
//		case GetBalance:
//			BCstate.GetBalance(string(From))
//		case Amalgamate:
//			BCstate.Amalgamate(string(From), string(To))
//		case UpdateBalance:
//			BCstate.UpdateBalance(string(From), Balance)
//		case UpdateSaving:
//			BCstate.UpdateSaving(string(From), Balance)
//		case SendPayment:
//			BCstate.SendPayment(string(From), string(To), Balance)
//		case WriteCheck:
//			BCstate.WriteCheck(string(From), Balance)
//		default:
//			fmt.Println("T doesn't match")
//		}
//	}
//}
