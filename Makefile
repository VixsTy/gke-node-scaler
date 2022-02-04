TEST?=$$(go list ./... | grep -v 'vendor')
TEST_RESULT_FOLDER=test-result


default: go-lint

testcover:
	mkdir -p ${TEST_RESULT_FOLDER}
	go test -v -coverprofile=${TEST_RESULT_FOLDER}/cover.out -cover ./...
	go tool cover -html=${TEST_RESULT_FOLDER}/cover.out -o ${TEST_RESULT_FOLDER}/coverage.html

PATH := ${PWD}/tools/bin:$(PATH)
GO_FILES = $(shell find "." ! -path "*vendor*" ! -path "*tools*" -name "*.go" -type f)

go-license:
	@export PATH=$$PATH; wwhrd check

go-format:
# @export PATH=$$PATH; gofumpt -s -w $(GO_FILES)
	@gofmt -s -w $(GO_FILES)

go-import:
	@export PATH=$$PATH; goimports -local ${LOCAL} -w $(GO_FILES)

go-lint: go-format go-import
	@export PATH=$$PATH; golangci-lint run
