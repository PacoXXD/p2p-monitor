GO := go
NAME := p2p-monitor
MAIN_GO := ./cmd/p2p-monitor/main.go
BUILDFLAGS := '-extldflags "-lm -lstdc++ -static"'
BUILDTAGS := netgo
CGO_ENABLED = 0

test:
	@echo "No Test"

all:
	make submodule
	make build

submodule:
	git submodule init
	git submodule update --remote

submodule_update:
	git submodule init
	git submodule update --remote --recursive

build:
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=amd64 $(GO) build -tags $(BUILDTAGS) -ldflags $(BUILDFLAGS) -o build/$(NAME) $(MAIN_GO)

docker:
	docker build -t p2p-monitor .

.PHONY: submodule submodule_update build docker
