.PHONY: all
all: test build doc

.PHONY: doc
doc:
	# requires github.com/princjef/gomarkdoc/cmd/gomarkdoc
	gomarkdoc --output=README.md ./...

.PHONY: test
test:
	gotip test -cover -covermode=atomic ./...

.PHONY: fmt
fmt:
	gotip fmt ./

.PHONY: build
build:
	gotip build ./