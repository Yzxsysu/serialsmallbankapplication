package abci

import (
	"encoding/binary"
	"log"
	"time"

	application "github.com/Yzxsysu/serialsmallbankapplication/v2/smallbankapplication/app"
	abcicode "github.com/tendermint/tendermint/abci/example/code"
	abcitypes "github.com/tendermint/tendermint/abci/types"
)

// 实现abci接口
var _ abcitypes.Application = (*SmallBankApplication)(nil)

// 定义KVStore程序的结构体
type SmallBankApplication struct {
	abcitypes.BaseApplication
	Node *application.BlockchainState
}

// 创建一个 ABCI APP
func NewSmallBankApplication(node *application.BlockchainState) *SmallBankApplication {
	return &SmallBankApplication{
		Node: node,
	}
}

func (app *SmallBankApplication) BeginBlock(req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	return abcitypes.ResponseBeginBlock{}
}

// 当新的交易被添加到Tendermint Core时，它会要求应用程序进行检查(验证格式、签名等)，当返回0时才通过
func (app SmallBankApplication) CheckTx(req abcitypes.RequestCheckTx) abcitypes.ResponseCheckTx {
	//var events []abcitypes.Event
	if app.Node.Leader {
		start := time.Now()
		app.Node.ResolveAndExecuteTx(&req.Tx)
		elapsed := time.Since(start)
		log.Printf("serail execute time: %s\n", elapsed)
	}
	return abcitypes.ResponseCheckTx{Code: abcicode.CodeTypeOK, GasUsed: 1}
}

// 这里我们创建了一个batch，它将存储block的交易。
func (app *SmallBankApplication) DeliverTx(req abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {
	if !app.Node.Leader {
		start := time.Now()
		app.Node.ResolveAndExecuteTx(&req.Tx)
		elapsed := time.Since(start)
		log.Printf("serail validate time: %s\n", elapsed)
	}
	return abcitypes.ResponseDeliverTx{Code: abcicode.CodeTypeOK}
}

func (app *SmallBankApplication) Commit() abcitypes.ResponseCommit {
	// 往数据库中提交事务，当 Tendermint core 提交区块时，会调用这个函数
	/*app.currentBatch.Commit()*/
	appHash := make([]byte, 8)
	binary.PutVarint(appHash, int64(app.Node.Height))
	app.Node.AppHash = appHash
	app.Node.Height++
	return abcitypes.ResponseCommit{Data: []byte{}}
}

func (app *SmallBankApplication) Query(reqQuery abcitypes.RequestQuery) (resQuery abcitypes.ResponseQuery) {
	return
}

func (SmallBankApplication) Info(req abcitypes.RequestInfo) abcitypes.ResponseInfo {
	return abcitypes.ResponseInfo{}
}

func (SmallBankApplication) InitChain(req abcitypes.RequestInitChain) abcitypes.ResponseInitChain {
	return abcitypes.ResponseInitChain{}
}

func (SmallBankApplication) EndBlock(req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	return abcitypes.ResponseEndBlock{}
}

func (SmallBankApplication) ListSnapshots(abcitypes.RequestListSnapshots) abcitypes.ResponseListSnapshots {
	return abcitypes.ResponseListSnapshots{}
}

func (SmallBankApplication) OfferSnapshot(abcitypes.RequestOfferSnapshot) abcitypes.ResponseOfferSnapshot {
	return abcitypes.ResponseOfferSnapshot{}
}

func (SmallBankApplication) LoadSnapshotChunk(abcitypes.RequestLoadSnapshotChunk) abcitypes.ResponseLoadSnapshotChunk {
	return abcitypes.ResponseLoadSnapshotChunk{}
}

func (SmallBankApplication) ApplySnapshotChunk(abcitypes.RequestApplySnapshotChunk) abcitypes.ResponseApplySnapshotChunk {
	return abcitypes.ResponseApplySnapshotChunk{}
}
