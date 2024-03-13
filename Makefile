GOFMT_FILES?=$$(find . -name '*.go')
export GO111MODULE=on

export TESTARGS=-race -coverprofile=coverage.txt -covermode=atomic

default: build

build:
	go build
	
build-docker:
	docker build -t eks-injector . 
	
test:
	go test ./internal/... -v $(TESTARGS) -timeout 120m -count=1

fmt:
	gofmt -w $(GOFMT_FILES)