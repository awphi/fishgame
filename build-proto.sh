# TODO fix this (or give up on making it work on windows)

tsproto_plugin="./client/node_modules/.bin/protoc-gen-ts_proto"

if [[ "$OSTYPE" == "cygwin" || "$OSTYPE" == "msys" || "$OSTYPE" == "win32" ]]; then
  tsproto_plugin="protoc-gen-ts_proto=.\\client\\node_modules\\.bin\\protoc-gen-ts_proto.cmd"
fi

client_out=./client/src/lib/game-client

mkdir -p $client_out
mkdir -p ./server/server
protoc -I=./proto --plugin=$tsproto_plugin --ts_proto_out=$client_out --ts_proto_opt=fileSuffix=.pb --go_out=./ ./proto/fishgame.proto