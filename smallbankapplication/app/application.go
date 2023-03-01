package app

import (
	dbm "github.com/Yzxsysu/serialsmallbankapplication/v2/src/tm-db"
	"log"
)

var (
	ByteLen  int
	CycleNum int
)

type BackendType string

// These are valid backend types.
const (
	GoLevelDBBackend BackendType = "goleveldb"
	CLevelDBBackend  BackendType = "cleveldb"
	MemDBBackend     BackendType = "memdb"
	BoltDBBackend    BackendType = "boltdb"
	RocksDBBackend   BackendType = "rocksdb"
	BadgerDBBackend  BackendType = "badgerdb"
)

type BlockchainState struct {
	// SmallBank contains two types of accounts
	CheckingStore dbm.DB
	SavingStore   dbm.DB
	Height        uint32
	Leader        bool
	AppHash       []byte
}

// NewBlockchainState can choose {goleveldb, cleveldb, memdb, boltdb, rocksdb, badgerdb}
// input Corresponding name
func NewBlockchainState(DBName BackendType, leader bool, dir string) (*BlockchainState, error, error) {
	var BaseCaseState BlockchainState
	var err error
	var err1 error
	// choose DB

	switch DBName {
	case GoLevelDBBackend:
		BaseCaseState.CheckingStore, err = dbm.NewDB("CheckingStore", dbm.GoLevelDBBackend, dir+"check")
		BaseCaseState.SavingStore, err1 = dbm.NewDB("SavingStore", dbm.GoLevelDBBackend, dir+"save")
	case CLevelDBBackend:
		BaseCaseState.CheckingStore, err = dbm.NewDB("CheckingStore", dbm.CLevelDBBackend, dir+"check")
		BaseCaseState.SavingStore, err1 = dbm.NewDB("SavingStore", dbm.CLevelDBBackend, dir+"save")
	case MemDBBackend:
		BaseCaseState.CheckingStore, err = dbm.NewDB("CheckingStore", dbm.MemDBBackend, dir+"check")
		BaseCaseState.SavingStore, err1 = dbm.NewDB("SavingStore", dbm.MemDBBackend, dir+"save")
	case BoltDBBackend:
		BaseCaseState.CheckingStore, err = dbm.NewDB("CheckingStore", dbm.BoltDBBackend, dir+"check")
		BaseCaseState.SavingStore, err1 = dbm.NewDB("SavingStore", dbm.BoltDBBackend, dir+"save")
	case RocksDBBackend:
		BaseCaseState.CheckingStore, err = dbm.NewDB("CheckingStore", dbm.RocksDBBackend, dir+"check")
		BaseCaseState.SavingStore, err1 = dbm.NewDB("SavingStore", dbm.RocksDBBackend, dir+"save")
	case BadgerDBBackend:
		BaseCaseState.CheckingStore, err = dbm.NewDB("CheckingStore", dbm.BadgerDBBackend, dir+"check")
		BaseCaseState.SavingStore, err1 = dbm.NewDB("SavingStore", dbm.BadgerDBBackend, dir+"save")
	}
	if err != nil || err1 != nil {
		log.Fatalf("Create db error: %v, %v", err, err1)
	}
	BaseCaseState.Leader = leader
	return &BlockchainState{
		CheckingStore: BaseCaseState.CheckingStore,
		SavingStore:   BaseCaseState.SavingStore,
		Height:        1,
		Leader:        leader,
	}, err, err1
}

// CreateAccount can create account with saving balance and checking balance
func (BCstate *BlockchainState) CreateAccount(AccountName string, SavingBalance int, CheckingBalance int) {
	// Create two separate accounts for two DB
	var err error
	err = BCstate.SavingStore.Set([]byte(AccountName), IntToBytes(SavingBalance))
	if err != nil {
		log.Println(err)
	}
	err = BCstate.CheckingStore.Set([]byte(AccountName), IntToBytes(CheckingBalance))
	if err != nil {
		log.Println(err)
	}
}

// GetBalance can get the balance of the account including check and save store
func (BCstate *BlockchainState) GetBalance(A string) {
	AddComplexity(ByteLen, CycleNum)

	Save, err := BCstate.SavingStore.Get([]byte(A))
	if err != nil {
		log.Println(err)
	}
	_ = BytesToInt(Save)

	Check, err := BCstate.CheckingStore.Get([]byte(A))
	if err != nil {
		log.Println(err)
	}
	_ = BytesToInt(Check)

}

func (BCstate *BlockchainState) Amalgamate(A string, B string) {
	AddComplexity(ByteLen, CycleNum)

	Save, err := BCstate.SavingStore.Get([]byte(A))
	if err != nil {
		log.Println(err)
	}
	SaveInt := BytesToInt(Save)

	Check, err := BCstate.CheckingStore.Get([]byte(B))
	if err != nil {
		log.Println(err)
	}
	CheckInt := BytesToInt(Check)

	SaveInt = SaveInt + CheckInt
	err = BCstate.SavingStore.Set([]byte(A), IntToBytes(SaveInt))
	if err != nil {
		log.Println(err)
	}

	err = BCstate.CheckingStore.Set([]byte(B), IntToBytes(0))
	if err != nil {
		log.Println(err)
	}
}

func (BCstate *BlockchainState) UpdateBalance(A string, Balance int) {
	AddComplexity(ByteLen, CycleNum)

	Check, err := BCstate.CheckingStore.Get([]byte(A))
	if err != nil {
		log.Println(err)
	}
	CheckInt := BytesToInt(Check)
	CheckInt += Balance

	err = BCstate.CheckingStore.Set([]byte(A), IntToBytes(CheckInt))
	if err != nil {
		log.Println(err)
	}
}

func (BCstate *BlockchainState) UpdateSaving(A string, Balance int) {
	AddComplexity(ByteLen, CycleNum)
	Save, err := BCstate.SavingStore.Get([]byte(A))
	if err != nil {
		log.Println(err)
	}
	SaveInt := BytesToInt(Save)
	SaveInt += Balance

	err = BCstate.CheckingStore.Set([]byte(A), IntToBytes(SaveInt))
	if err != nil {
		log.Println(err)
	}
}

func (BCstate *BlockchainState) SendPayment(A string, B string, Balance int) {
	AddComplexity(ByteLen, CycleNum)

	CheckA, err := BCstate.CheckingStore.Get([]byte(A))
	if err != nil {
		log.Println(err)
	}

	CheckIntA := BytesToInt(CheckA)

	CheckB, err := BCstate.CheckingStore.Get([]byte(B))
	if err != nil {
		log.Println(err)
	}
	CheckIntB := BytesToInt(CheckB)

	CheckIntA -= Balance
	CheckIntB += Balance

	// update check value
	err = BCstate.CheckingStore.Set([]byte(A), IntToBytes(CheckIntA))
	if err != nil {
		log.Println(err)
	}

	err = BCstate.CheckingStore.Set([]byte(B), IntToBytes(CheckIntB))
	if err != nil {
		log.Println(err)
	}
}

func (BCstate *BlockchainState) WriteCheck(A string, Balance int) {
	AddComplexity(ByteLen, CycleNum)

	Save, err := BCstate.SavingStore.Get([]byte(A))
	if err != nil {
		log.Println(err)
	}
	SaveInt := BytesToInt(Save)

	Check, err := BCstate.CheckingStore.Get([]byte(A))
	if err != nil {
		log.Println(err)
	}
	CheckInt := BytesToInt(Check)

	if SaveInt+CheckInt < Balance {
		CheckInt = CheckInt - Balance - 1
	} else {
		CheckInt = CheckInt - Balance
	}
	err = BCstate.CheckingStore.Set([]byte(A), IntToBytes(CheckInt))
	if err != nil {
		log.Println(err)
	}
}
