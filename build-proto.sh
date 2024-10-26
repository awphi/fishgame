# TODO fix this (or give up on making it work on windows)

tsproto_plugin="./client/node_modules/.bin/protoc-gen-ts_proto"

if [[ "$OSTYPE" == "cygwin" || "$OSTYPE" == "msys" || "$OSTYPE" == "win32" ]]; then
  tsproto_plugin="protoc-gen-ts_proto=.\\client\\node_modules\\.bin\\protoc-gen-ts_proto.cmd"
fi

mkdir -p ./client/lib/gen
mkdir -p ./server/gen
protoc -I=./proto --plugin=$tsproto_plugin --ts_proto_out=./client/lib/gen --go_out=./ ./proto/example.proto