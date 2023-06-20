GOSRC=/home/WorkPlace
TEST_SCENE="serial_chain"
TM_HOME="/home/.serial_tendermint"
WORKSPACE="$GOSRC/github.com/Yzxsysu/serialsmallbankapplication"
CURRENT_DATE=`date +"%Y-%m-%d-%H-%M"`
LOG_DIR="$WORKSPACE/tmplog/$TEST_SCENE-$CURRENT_DATE"
DURATION=360

rm -rf $TM_HOME

mkdir -p $TM_HOME
mkdir -p $LOG_DIR

groupNum=$1
nodeId=$2
division=$3
echo "group number: $groupNum, node id: $nodeId"

cp -r /home/WorkPlace/github.com/Yzxsysu/serialsmallbankapplication/config/* $TM_HOME
echo "configs generated"

pkill -9 chain

case $nodeId in
    0)
    ./build/chain/chain -home $TM_HOME/node${nodeId} -leader "true" -accountNum 100000 -coreNum 16 &> $LOG_DIR/node${nodeId}.log &
    echo "the node Id is ${nodeId}"
    ;;
    1)
    ./build/chain/chain -home $TM_HOME/node${nodeId} -leader "false" -accountNum 100000 -coreNum 16 &> $LOG_DIR/node${nodeId}.log &
    echo "the node Id is ${nodeId}"
    ;;
    2)
    ./build/chain/chain -home $TM_HOME/node${nodeId} -leader "false" -accountNum 100000 -coreNum 16 &> $LOG_DIR/node${nodeId}.log &
    echo "the node Id is ${nodeId}"
    ;;
    3)
    ./build/chain/chain -home $TM_HOME/node${nodeId} -leader "false" -accountNum 100000 -coreNum 16 &> $LOG_DIR/node${nodeId}.log &
    echo "the node Id is ${nodeId}"
    ;;
esac

echo "testnet launched"
echo "running for ${DURATION}s..."
sleep $DURATION

pkill -9 chain
echo "all done"
