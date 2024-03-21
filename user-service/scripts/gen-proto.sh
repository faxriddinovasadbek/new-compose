#!/bin/bash
CURRENT_DIR=$(pwd)

protoc -I /usr/local/include \
       -I $GOPATH/src/github.com/gogo/protobuf/gogoproto \
       -I "$CURRENT_DIR/protos/" \
        --gofast_out=plugins=grpc:"$CURRENT_DIR/genproto/" \
        "$CURRENT_DIR/protos/user.proto";

# Remove ,omitempty from generated code
for module in $(find "$CURRENT_DIR/genproto" -type d); do
    if [[ "$OSTYPE" == "darwin"* ]]; then
        sed -i "" -e "s/,omitempty//g" "$module"/*.go || exit 1
    else
        sed -i -e "s/,omitempty//g" "$module"/*.go || exit 1
    fi
done

