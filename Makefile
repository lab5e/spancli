ifeq ($(GOPATH),)
GOPATH := $(HOME)/go
endif

all: test lint vet build

clean:
	@find . -name "*-wal" -delete
	@find . -name "*-shm" -delete
	@rm -f bin/*.linux

build: scl

span:
	@cd cmd/span && go build -o ../../bin/span

check: lint vet staticcheck revive

lint:
	@revive -exclude ./... 

vet:
	@go vet ./...

staticcheck:
	@staticcheck ./...

test:
	@go test ./...

test_verbose:
	@go test ./... -v

test_cover:
	@go test ./... -cover -coverprofile=unittests.cover -coverpkg=github.com/lab5e/spancli/backend/pkg/...

test_race:
	@go test ./... -race

test_all: test_cover test_race

benchmark:
	cd output && go test -bench .
