-include .env
export
CURRENT_DIR=${PWD}

.SILENT:
.PHONY: run
run:
	go run cmd/main.go

gen-proto:
	protoc -I=${CURRENT_DIR}/protos --go_out=${CURRENT_DIR} \
		--go-grpc_out=${CURRENT_DIR} ${CURRENT_DIR}/protos/*.proto

run:
	go run cmd/main.go