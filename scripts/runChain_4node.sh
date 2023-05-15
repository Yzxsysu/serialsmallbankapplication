GOSRC=/home/WorkPlace
TEST_SCENE="serial_chain"
TM_HOME="/home/.serial_tendermint"
WORKSPACE="$GOSRC/github.com/Yzxsysu/serialsmallbankapplication"
CURRENT_DATE=`date +"%Y-%m-%d-%H-%M"`
LOG_DIR="$WORKSPACE/tmplog/$TEST_SCENE-$CURRENT_DATE"
DURATION=180

rm -rf $TM_HOME

mkdir -p $TM_HOME
mkdir -p $LOG_DIR

cp -r /home/WorkPlace/github.com/Yzxsysu/serialsmallbankapplication/config/* $TM_HOME
echo "configs generated"

pkill -9 chain

./build/chain/chain -home $TM_HOME/4node/node0 -leader "true" -accountNum 1000 -coreNum 16 &> $LOG_DIR/node0.log &

./build/chain/chain -home $TM_HOME/4node/node1 -leader "false" -accountNum 1000 -coreNum 16 &> $LOG_DIR/node1.log &

./build/chain/chain -home $TM_HOME/4node/node2 -leader "false" -accountNum 1000 -coreNum 16 &> $LOG_DIR/node2.log &

./build/chain/chain -home $TM_HOME/4node/node3 -leader "false" -accountNum 1000 -coreNum 16 &> $LOG_DIR/node3.log &

echo "testnet launched"
echo "running for ${DURATION}s..."
sleep $DURATION

pkill -9 chain
echo "all done"
