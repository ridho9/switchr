PREFIX ?= $(HOME)/.local/bin

.PHONY: all build install clean

all: install

build:
	go build -o switcher .

install: build
	mkdir -p $(PREFIX)
	install -m 755 switcher $(PREFIX)/switcher
	install -m 755 contrib/term-launcher $(PREFIX)/term-launcher

clean:
	rm -f switcher
