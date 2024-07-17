# project variables
PROJECT_NAME := WebIDE
PKG := ./...
MAIN := ./main.go

# go commands
BUILD := go build
CLEAN := go clean
FMT := go fmt
VET := go vet
TEST := go test
RUN := go run

# targets
.PHONY: all build clean fmt vet test run

all: fmt vet test build

build:
	$(BUILD) -o $(PROJECT_NAME) $(MAIN)

clean:
	$(CLEAN)
	rm -f $(PROJECT_NAME)

fmt:
	$(FMT) $(PKG)

vet:
	$(VET) $(PKG)

test:
	$(TEST) $(PKG)

run:
	$(RUN) $(MAIN)