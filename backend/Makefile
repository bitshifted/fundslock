GOPATH := $(shell go env GOPATH)
TARGET_DIR := target
TIMESTAMP := $(shell date +%Y%m%d%H%M%S)
VERSION ?= "0.0.0"
GIT_COMMIT_ID ?= "unknown"
VERSION_FLAG := "github.com/bitshifted/fundslock-be/cli.version=${VERSION}"
BUILDNUM_FLAG := "github.com/bitshifted/fundslock-be/cli.buildNum=$(TIMESTAMP)"
COMMIT_ID_FLAG := "github.com/bitshifted/fundslock-be/cli.commitID=$(GIT_COMMIT_ID)"
LDFLAGS := -ldflags="-X '$(VERSION_FLAG)' -X '$(BUILDNUM_FLAG)' -X '$(COMMIT_ID_FLAG)'"

init:
	mkdir -p  ${TARGET_DIR}

clean:
	rm -rvf ${TARGET_DIR}

lint:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2
	$(GOPATH)/bin/golangci-lint run ./...
	go mod tidy

test:
	go test  -coverprofile=${TARGET_DIR}/coverage.out ./...
	go tool cover -html=${TARGET_DIR}/coverage.out -o ${TARGET_DIR}/coverage.html

build: clean init check-license-headers lint  test
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o ${TARGET_DIR}/linux-amd64/fundslock-be
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o ${TARGET_DIR}/linux-arm64/fundslock-be
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o ${TARGET_DIR}/windows-amd64/fundslock-be.exe
	GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o ${TARGET_DIR}/windows-arm64/fundslock-be.exe
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o ${TARGET_DIR}/macos-amd64/fundslock-be
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o ${TARGET_DIR}/macos-arm64/fundslock-be


add-license-headers:
	go install github.com/google/addlicense@v1.1.1
	$(GOPATH)/bin/addlicense -v -c 'Bitshift ED' -y 2026 -l mpl -s=only ./**/*.go main.go Dockerfile Makefile

check-license-headers:
	go install github.com/google/addlicense@v1.1.1
	$(GOPATH)/bin/addlicense -v -check ./**/*.go main.go Dockerfile Makefile
