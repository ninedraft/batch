.PHONY: all
all: test build

.PHONY: test
test:
	gotip test ./...

.PHONY: fmt
fmt:
	gotip fmt ./

.PHONY: build
build:
	gotip build ./