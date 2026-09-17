VERSION ?= dev

.PHONY: build test clean

build:
	go build -ldflags "-X main.version=$(VERSION)" -o tc .

test:
	go test ./...

clean:
	rm -f tc
