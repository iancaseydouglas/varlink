.PHONY: build test clean

build:
    go build -o bin/varlink cmd/varlink/main.go

test:
    go test ./...

clean:
    rm -f bin/varlink
