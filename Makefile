export GO111MODULE=on
BINARY_NAME=analytics-file-export-service

all: deps build
install:
	go install cmd/$(BINARY_NAME)/$(BINARY_NAME).go
build:
	go build cmd/$(BINARY_NAME)/$(BINARY_NAME).go
test:
	go test -v ./...
clean:
	go clean -v ./...
	rm -f $(BINARY_NAME)
deps:
	go build -v ./...
upgrade:
	go get -v ./...