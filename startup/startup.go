package startup

import (
	"flag"
	"github.com/Lawliet-Chan/yu/blockchain"
	"github.com/Lawliet-Chan/yu/common"
	"github.com/Lawliet-Chan/yu/config"
	"github.com/Lawliet-Chan/yu/node/master"
	"github.com/Lawliet-Chan/yu/state"
	"github.com/Lawliet-Chan/yu/tripod"
	"github.com/Lawliet-Chan/yu/txpool"
	"github.com/Lawliet-Chan/yu/utils/codec"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	masterCfgPath string
	chainCfgPath  string
	baseCfgPath   string
	txpoolCfgPath string
	stateCfgPath  string

	masterCfg config.MasterConf
	chainCfg  config.BlockchainConf
	baseCfg   config.BlockBaseConf
	txpoolCfg config.TxpoolConf
	stateCfg  config.StateConf
)

func StartUp(tripods ...tripod.Tripod) {
	initCfgFromFlags()
	initLog()

	codec.GlobalCodec = &codec.RlpCodec{}
	gin.SetMode(gin.ReleaseMode)

	chain, err := blockchain.NewBlockChain(&chainCfg)
	if err != nil {
		logrus.Panicf("load blockchain error: %s", err.Error())
	}
	base, err := blockchain.NewBlockBase(&baseCfg)
	if err != nil {
		logrus.Panicf("load blockbase error: %s", err.Error())
	}

	var pool txpool.ItxPool
	switch masterCfg.RunMode {
	case common.LocalNode:
		pool = txpool.LocalWithDefaultChecks(&txpoolCfg)
	case common.MasterWorker:
		logrus.Panic("no server txpool")
	}

	stateStore, err := state.NewStateStore(&stateCfg)
	if err != nil {
		logrus.Panicf("load stateKV error: %s", err.Error())
	}

	land := tripod.NewLand()
	land.SetTripods(tripods...)

	m, err := master.NewMaster(&masterCfg, chain, base, pool, stateStore, land)
	if err != nil {
		logrus.Panicf("load master error: %s", err.Error())
	}

	m.Startup()
}

func initCfgFromFlags() {
	flag.StringVar(&masterCfgPath, "m", "yu_conf/master.toml", "Master config file path")
	config.LoadConf(masterCfgPath, &masterCfg)

	flag.StringVar(&chainCfgPath, "c", "yu_conf/blockchain.toml", "blockchain config file path")
	config.LoadConf(chainCfgPath, &chainCfg)

	flag.StringVar(&baseCfgPath, "b", "yu_conf/blockbase.toml", "blockbase config file path")
	config.LoadConf(baseCfgPath, &baseCfg)

	flag.StringVar(&txpoolCfgPath, "tp", "yu_conf/txpool.toml", "txpool config file path")
	config.LoadConf(txpoolCfgPath, &txpoolCfg)

	flag.StringVar(&stateCfgPath, "s", "yu_conf/state.toml", "state config file path")
	config.LoadConf(stateCfgPath, &stateCfg)
}

func initLog() {
	formatter := &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	}
	logrus.SetFormatter(formatter)
	logrus.SetLevel(logrus.InfoLevel)
}