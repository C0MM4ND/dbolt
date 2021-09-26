BRANCH=`git rev-parse --abbrev-ref HEAD`
COMMIT=`git rev-parse --short HEAD`
GOLDFLAGS="-X main.branch $(BRANCH) -X main.commit $(COMMIT)"

race:
	@TEST_FREELIST_TYPE=hashmap go test ./tests -v -race -run="TestSimulate_(100op|1000op)"
	@echo "array freelist test"
	@TEST_FREELIST_TYPE=array go test ./tests -v -race -run="TestSimulate_(100op|1000op)"

fmt:
	!(gofmt -l -s -d $(shell find . -name '*.go') | grep '[a-z]')

imports:
	goimports -l $(shell find . -name '*.go')

external-tests:
	TEST_FREELIST_TYPE=hashmap go test ./tests -timeout 20m -v -coverprofile cover.out -covermode atomic

	@echo "array freelist test"

	@TEST_FREELIST_TYPE=array go test ./tests -timeout 20m -v -coverprofile cover.out -covermode atomic

test:
	TEST_FREELIST_TYPE=hashmap go test -timeout 20m -v -coverprofile cover.out -covermode atomic
	TEST_FREELIST_TYPE=hashmap go test -v ./cmd/dbolt

	@echo "array freelist test"

	@TEST_FREELIST_TYPE=array go test -timeout 20m -v -coverprofile cover.out -covermode atomic
	@TEST_FREELIST_TYPE=array go test -v ./cmd/dbolt

.PHONY: fmt imports test race
