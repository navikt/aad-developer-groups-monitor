BUILDTIME = $(shell date "+%s")
DATE = $(shell date "+%Y-%m-%d")
LAST_COMMIT = $(shell git rev-parse --short HEAD)
LDFLAGS := -X github.com/navikt/aad-developer-groups-monitor/pkg/version.Revision=$(LAST_COMMIT) -X github.com/navikt/aad-developer-groups-monitor/pkg/version.Date=$(DATE) -X github.com/navikt/aad-developer-groups-monitor/pkg/version.BuildUnixTime=$(BUILDTIME)

.PHONY: monitor test fmt check alpine

all: test check fmt build

build:
	go build -o bin/monitor -ldflags "-s $(LDFLAGS)" cmd/monitor/*.go

test:
	go test ./...

fmt:
	go run mvdan.cc/gofumpt -w ./

check:
	go run honnef.co/go/tools/cmd/staticcheck ./...
	go run golang.org/x/vuln/cmd/govulncheck ./...

alpine:
	go build -a -installsuffix cgo -o bin/monitor -ldflags "-s $(LDFLAGS)" cmd/monitor/main.go
