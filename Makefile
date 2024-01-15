PROTO_DIR = proto
RELEASE_TAG = $(shell git tag --abbrev=0)

ifeq ($(OS), Windows_NT)
	OS = windows
	SHELL := powershell.exe
	.SHELLFLAGS := -NoProfile -Command
	PACKAGE = $(shell Get-Content go.mod -head 1 | Foreach-Object { $$data = $$_ -split " "; "{0}" -f $$data[1]})
	BIN = proto-go-course.exe
else
	UNAME := $(shell uname -s)
	ifeq ($(UNAME),Darwin)
		OS = macos
	else ifeq ($(UNAME),Linux)
		OS = linux
	else
    	$(error OS not supported by this Makefile)
	endif
	PACKAGE = $(shell head -1 go.mod | awk '{print $$2}')
	BIN = proto-go-course
endif

generate:
	protoc protoc -Iproto --go_opt=module=github.com/kifiApp/common --go-grpc_opt=module=github.com/kifiApp/common  --go_out=./gen --go-grpc_out=./gen proto/*.proto

modsync: generate
	go get -u ./...
	go mod tidy

clean:
	rm ${PROTO_DIR}/*.pb.go
	rm ${BIN}
