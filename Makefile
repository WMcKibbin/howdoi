.PHONY: build test lint clean install

BINARY := bin/howdoi
MODULE := github.com/WMcKibbin/howdoi

build:
	mkdir -p bin
	go build -o $(BINARY) .

install:
	go install .

test:
	go test ./...

lint:
	golangci-lint run ./...

clean:
	rm -f $(BINARY)

snapshot:
	goreleaser --snapshot --clean
