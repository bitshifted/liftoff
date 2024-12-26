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
	mkdir -p  ${TARGET_DIR}/{linux-amd64,linux-arm64,windows-amd64,windows-arm64,macos-amd64,macos-arm64}

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
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o ${TARGET_DIR}/linux-amd64/liftoff
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o ${TARGET_DIR}/linux-arm64/liftoff
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o ${TARGET_DIR}/windows-amd64/liftoff.exe
	GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o ${TARGET_DIR}/windows-arm64/liftoff.exe
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o ${TARGET_DIR}/macos-amd64/liftoff
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o ${TARGET_DIR}/macos-arm64/liftoff

package: build
	mkdir ${TARGET_DIR}/dist
	tar -C ${TARGET_DIR}/linux-amd64 -czvf ${TARGET_DIR}/dist/liftoff-linux-amd64-${VERSION}.tar.gz liftoff
	tar -C ${TARGET_DIR}/linux-arm64 -czvf ${TARGET_DIR}/dist/liftoff-linux-arm64-${VERSION}.tar.gz liftoff
	tar -C ${TARGET_DIR}/macos-amd64 -czvf ${TARGET_DIR}/dist/liftoff-macos-amd64-${VERSION}.tar.gz liftoff
	tar -C ${TARGET_DIR}/macos-arm64 -czvf ${TARGET_DIR}/dist/liftoff-macos-arm64-${VERSION}.tar.gz liftoff
	zip -j  ${TARGET_DIR}/dist/liftoff-windows-amd64-${VERSION}.zip target/windows-amd64/liftoff.exe
	zip -j  ${TARGET_DIR}/dist/liftoff-windows-arm64-${VERSION}.zip target/windows-arm64/liftoff.exe

add-license-headers:
	go install github.com/google/addlicense@v1.1.1
	$(GOPATH)/bin/addlicense -v -c 'Bitshift D.O.O' -y 2024 -l mpl -s=only ./**/*.go main.go Dockerfile Makefile

check-license-headers:
	go install github.com/google/addlicense@v1.1.1
	$(GOPATH)/bin/addlicense -v -check ./**/*.go main.go Dockerfile Makefile
