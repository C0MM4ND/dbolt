BRANCH=`git rev-parse --abbrev-ref HEAD`
COMMIT=`git rev-parse --short HEAD`
GOLDFLAGS="-X main.branch $(BRANCH) -X main.commit $(COMMIT)"

race:
	go test ./tests -v -race -run="TestSimulate_(100op|1000op)"

fmt:
	!(gofmt -l -s -d $(shell find . -name '*.go') | grep '[a-z]')

imports:
	goimports -l $(shell find . -name '*.go')

external-tests:
	go test ./tests -timeout 20m -v -coverprofile cover.out -covermode atomic

test:
	go test -timeout 20m -v -coverprofile cover.out -covermode atomic
	go test -v ./cmd/dbolt

.PHONY: fmt imports test race
