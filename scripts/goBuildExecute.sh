GOSRC=/home/WorkPlace
ROOT=$GOSRC/github.com/Yzxsysu/serialsmallbankapplication

mkdir -p build

go build -o build/chain $ROOT/cmd/chain

chmod +x build/*
