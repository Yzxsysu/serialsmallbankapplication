package main

import (
	"fmt"
	application "github.com/Yzxsysu/serialsmallbankapplication/v2/smallbankapplication/app"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"time"
)

// tx format: 127.0.0.1:20057/broadcast_tx_commit?tx="T=3,I=1,F=1,O=3,B=156>T=1,I=2,F=2,O=1,B=190"
func main() {
	go func() {
		http.ListenAndServe(":8084", nil)
	}()
	for {
		var err error
		txNum := 1000
		txs := application.GenerateTx(float64(txNum), 1000, 1)
		txStr := make([]string, txNum)
		//l := len(txs)
		for i, tx := range txs {
			str := ""
			str += "T=" + strconv.Itoa(int(tx.T))
			str += "," + "I=" + strconv.Itoa(int(tx.I))
			str += "," + "F=" + string(tx.F)
			str += "," + "O=" + string(tx.O)
			str += "," + "B=" + strconv.Itoa(tx.B)
			txStr[i] = str
		}
		for _, tx := range txStr {
			go func(str string) {
				fmt.Println(len(str))
				request1 := "127.0.0.1:20057/broadcast_tx_commit?tx=\"" + str + "\""
				_, err = http.Get("http://" + request1)
				if err != nil {
					fmt.Println(err)
				}
			}(tx)
		}
		time.Sleep(time.Millisecond * 10000)
		//go func(str string) {
		//	fmt.Println(len(str))
		//	request1 := "127.0.0.1:20057/broadcast_tx_commit?tx=\"" + str + "\""
		//	_, err = http.Get("http://" + request1)
		//	if err != nil {
		//		fmt.Println(err)
		//	}
		//}(str)
	}
}
