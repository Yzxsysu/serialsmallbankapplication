package main

import (
	"flag"
	"fmt"
	smallbankapplication "github.com/Yzxsysu/serialsmallbankapplication/v2/smallbankapplication/abci"
	application "github.com/Yzxsysu/serialsmallbankapplication/v2/smallbankapplication/app"
	"github.com/spf13/viper"
	abciclient "github.com/tendermint/tendermint/abci/client"
	cfg "github.com/tendermint/tendermint/config"
	tmlog "github.com/tendermint/tendermint/libs/log"
	nm "github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/types"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
)

var homeDir string
var isLeader string
var accountNum, coreNum uint

func init() {
	flag.StringVar(&isLeader, "leader", "false", "Is it a leader (default: false)")
	flag.UintVar(&accountNum, "accountNum", 1000, "The account num of the SmallBank")
	flag.UintVar(&coreNum, "coreNum", 8, "control the num of cpu's cores")
	flag.StringVar(&homeDir, "home", "", "Path to the tendermint config directory (if empty, uses $HOME/.tendermint)")
}

func main() {
	go func() {
		http.ListenAndServe(":8083", nil)
	}()
	flag.Parse()
	runtime.GOMAXPROCS(int(coreNum))
	if homeDir == "" {
		homeDir = os.ExpandEnv("$HOME/.tendermint")
	}
	config := cfg.DefaultConfig()

	config.SetRoot(homeDir)

	viper.SetConfigFile(fmt.Sprintf("%s/%s", homeDir, "config/config.toml"))
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Reading config: %v", err)
	}
	if err := viper.Unmarshal(config); err != nil {
		log.Fatalf("Decoding config: %v", err)
	}
	if err := config.ValidateBasic(); err != nil {
		log.Fatalf("Invalid configuration data: %v", err)
	}
	gf, err := types.GenesisDocFromFile(config.GenesisFile())
	if err != nil {
		log.Fatalf("Loading genesis document: %v", err)
	}

	dbPath := filepath.Join(homeDir, "badger")
	db, err, err1 := application.NewBlockchainState("badgerdb", true, dbPath)
	if err != nil || err1 != nil {
		log.Fatalf("Opening database: %v", err)
	}
	defer func() {
		if err := db.SavingStore.Close(); err != nil {
			log.Fatalf("Closing SavingStore database: %v", err)
		}
		if err := db.CheckingStore.Close(); err != nil {
			log.Fatalf("Closing Checkingstore database: %v", err)
		}
	}()
	if isLeader == "true" {
		db.Leader = true
	} else if isLeader == "false" {
		db.Leader = false
	}

	app := smallbankapplication.NewSmallBankApplication(db)
	acc := abciclient.NewLocalCreator(app)

	logger := tmlog.MustNewDefaultLogger(tmlog.LogFormatPlain, tmlog.LogLevelInfo, false)
	node, err := nm.New(config, logger, acc, gf)
	if err != nil {
		log.Fatalf("Creating node: %v", err)
	}

	err = node.Start()
	if err != nil {
		log.Fatalf("Starting node: %v", err)
	}
	defer func() {
		err = node.Stop()
		if err != nil {
			log.Fatalf("Stoping node: %v", err)
		}
		node.Wait()
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
