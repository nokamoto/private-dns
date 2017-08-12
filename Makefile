all: deps test install

deps:
				go get google.golang.org/grpc
				go get github.com/golang/protobuf/protoc-gen-go

install:
				protoc --go_out=plugins=grpc:. proto/service.proto
				go install ./dnscli

test:
				go test ./dnscli
