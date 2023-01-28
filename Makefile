fmt:
	go fmt ./...
proto:
	protoc -I=./history-rewinded-regan/proto --go_out=./history-rewinded-regan/pb --go-grpc_out=./history-rewinded-regan/pb ./history-rewinded-regan/proto/lear.proto
