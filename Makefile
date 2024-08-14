GOPATH := $(shell go env GOPATH)
TARGET_DIR := target

init:
	mkdir ${TARGET_DIR}

clean:
	rm -rvf ${TARGET_DIR}

lint:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.1
	go install github.com/golangci/golangci-lint/cmd/golangci-lint
	$(GOPATH)/bin/golangci-lint run ./...
	go mod tidy

test:
	go test  -coverprofile=${TARGET_DIR}/coverage.out ./...
	go tool cover -html=${TARGET_DIR}/coverage.out -o ${TARGET_DIR}/coverage.html

build: clean init lint test
	GOOS=linux go build -o ${TARGET_DIR}/liftoff
	GOOS=windows go build -o ${TARGET_DIR}/liftoff.exe
	GOOS=darwin go build -o ${TARGET_DIR}/liftoff

