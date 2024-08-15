GOPATH := $(shell go env GOPATH)
TARGET_DIR := target
TIMESTAMP := $(shell date +%Y%m%d%H%M%S)
VERSION ?= "0.0.0"
GIT_COMMIT_ID ?= "unknown"
VERSION_FLAG := "github.com/bitshifted/liftoff/cli.version=${VERSION}"
BUILDNUM_FLAG := "github.com/bitshifted/liftoff/cli.buildNum=$(TIMESTAMP)"
COMMIT_ID_FLAG := "github.com/bitshifted/liftoff/cli.commitID=$(GIT_COMMIT_ID)"
LDFLAGS := -ldflags="-X '$(VERSION_FLAG)' -X '$(BUILDNUM_FLAG)' -X '$(COMMIT_ID_FLAG)'"

init:
	mkdir -p  ${TARGET_DIR}/{linux,windows,macos}

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

build: clean init check-license-headers lint  test
	GOOS=linux go build $(LDFLAGS) -o ${TARGET_DIR}/linux/liftoff
	GOOS=windows go build $(LDFLAGS) -o ${TARGET_DIR}/windows/liftoff.exe
	GOOS=darwin go build $(LDFLAGS) -o ${TARGET_DIR}/macos//liftoff

add-license-headers:
	go install github.com/google/addlicense@v1.1.1
	$(GOPATH)/bin/addlicense -v -c 'Bitshift D.O.O' -y 2024 -l mpl -s=only ./**/*.go main.go Dockerfile Makefile

check-license-headers:
	go install github.com/google/addlicense@v1.1.1
	$(GOPATH)/bin/addlicense -v -check ./**/*.go main.go Dockerfile Makefile
