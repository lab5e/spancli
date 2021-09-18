ifeq ($(GOPATH),)
GOPATH := $(HOME)/go
endif

all: test span

release: all
	@cd cmd/span && GOOS=darwin go build -o ../../bin/macos/span && cd ../../bin/macos && tar czf ../span-macOS.tar.gz span
	@cd cmd/span && GOOS=linux go build -o ../../bin/linux/span  && cd ../../bin/linux && tar czf ../span-linux.tar.gz span
	@cd cmd/span && GOOS=windows go build -o ../../bin/windows/span.exe  && cd ../../bin/windows && zip ../span-windows.zip span.exe

clean:
	@find . -name "*-wal" -delete
	@find . -name "*-shm" -delete
	@rm -f bin/*.linux

span:
	@cd cmd/span && go build -o ../../bin/span

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
