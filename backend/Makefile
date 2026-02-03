BIN=opentrans

VERSION=$(shell git describe --tags --abbrev=0 2>/dev/null || echo "dev")
DATE=$(shell date +%Y%m%d%H%M%S)
COMMIT:=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

LDFLAGS=-s -w
LDFLAGS+= -X github.com/fdddf/opentrans/cmd.Version=$(VERSION)
LDFLAGS+= -X github.com/fdddf/opentrans/cmd.Commit=$(COMMIT)
LDFLAGS+= -X github.com/fdddf/opentrans/cmd.Build=$(DATE)

all: binary

binary:
	go build -v -ldflags="$(LDFLAGS)" -o $(BIN) main.go

ui:
	cd web && npm install && npm run build

clean:
	rm -f $(BIN)

test:
	go test -v

install:
	go install -ldflags="$(LDFLAGS)"


.PHONY: all binary ui clean test install
