ifeq ($(GOPATH),)
GOPATH := $(HOME)/go
endif

ifeq ($(VERSION),)
VERSION := $(shell git tag -l --sort=-version:refname | head -n 1 | cut -c 2-)
endif

LDFLAGS := "-X github.com/lab5e/spancli/pkg/global.Version=$(VERSION)"

all: test vet span

release: all
	@cd cmd/span && GOOS=linux GOARCH=amd64 go build -ldflags=$(LDFLAGS) -o ../../bin/span.amd64-linux
	@cd cmd/span && GOOS=darwin GOARCH=amd64 go build -ldflags=$(LDFLAGS) -o ../../bin/span.amd64-macOS
	@cd cmd/span && GOOS=darwin GOARCH=arm64 go build -ldflags=$(LDFLAGS) -o ../../bin/span.arm64-macOS
	@cd cmd/span && GOOS=windows GOARCH=amd64 go build -ldflags=$(LDFLAGS) -o ../../bin/span.amd64-win.exe
	@cd cmd/span && GOOS=linux GOARCH=arm GOARM=5 go build -ldflags=$(LDFLAGS) -o ../../bin/span.arm5-rpi-linux

clean:
	@rm -rf bin

span:
	@cd cmd/span && go build  -ldflags=$(LDFLAGS) -o ../../bin/span

test:
	@go test ./...

vet:
	@go vet ./...

check: lint vet staticcheck

lint:
	@revive -exclude ./... 

staticcheck:
	@staticcheck ./...

test_verbose:
	@go test ./... -v

test_cover:
	@go test ./... -cover -coverprofile=unittests.cover -coverpkg=github.com/lab5e/spancli/backend/pkg/...

test_race:
	@go test ./... -race

test_all: test_cover test_race

benchmark:
	cd output && go test -bench .
