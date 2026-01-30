BINARY := dotfiles-tui
PKG    := ./cmd/dotfiles-tui

.PHONY: build install clean test lint

build:
	go build -o bin/$(BINARY) $(PKG)

install:
	go install $(PKG)

clean:
	rm -rf bin/

test:
	go test -race ./...

lint:
	go vet ./...
