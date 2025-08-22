VERSION = $(shell godzil show-version)
CURRENT_REVISION = $(shell git rev-parse --short HEAD)
BUILD_LDFLAGS = "-s -w -X github.com/Songmu/maltmill.revision=$(CURRENT_REVISION)"
ifdef update
  u=-u
endif

.PHONY: deps
deps:
	go get ${u} -v -t

.PHONY: devel-deps
devel-deps: deps
	go install github.com/Songmu/godzil/cmd/godzil@latest

.PHONY: deps
test: deps
	go test

.PHONY: build
build: deps
	go build -ldflags=$(BUILD_LDFLAGS) ./cmd/maltmill

CREDITS: devel-deps go.sum
	godzil credits -w

DIST_DIR = dist
.PHONY: crossbuild
crossbuild: devel-deps
	rm -rf $(DIST_DIR)
	godzil crossbuild -pv=v$(VERSION) -build-ldflags=$(BUILD_LDFLAGS) \
	  -d=$(DIST_DIR) ./cmd/*
	cd $(DIST_DIR) && shasum -a 256 $$(find * -type f -maxdepth 0) > SHA256SUMS

.PHONY: upload
upload:
	ghr v$(VERSION) dist
