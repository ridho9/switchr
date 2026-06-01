PREFIX ?= $(HOME)/.local/bin

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS  = -X main.Version=$(VERSION) -X main.Commit=$(COMMIT)

.PHONY: all build install clean

all: install

build:
	go build -ldflags "$(LDFLAGS)" -o switcher .

install: build
	mkdir -p $(PREFIX)
	install -m 755 switcher $(PREFIX)/switcher
	install -m 755 contrib/term-launcher $(PREFIX)/term-launcher

clean:
	rm -f switcher
