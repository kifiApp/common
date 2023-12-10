# common

protoc -Iproto --go_opt=module=github.com/kifiApp/common --go-grpc_opt=module=github.com/kifiApp/common  --go_out=./gen --go-grpc_out=./gen proto/*.proto