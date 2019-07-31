GOMOD ?= on
GO ?= GO111MODULE=$(GOMOD) go

#Don't enable mod=vendor when GOMOD is off or else go build/install will fail
GOMODFLAG ?=-mod=vendor
ifeq ($(GOMOD), off)
GOMODFLAG=
endif

#retrieve go version details for version check
GO_VERSION     := $(shell $(GO) version | sed -e 's/^[^0-9.]*\([0-9.]*\).*/\1/')
GO_VERSION_MAJ := $(shell echo $(GO_VERSION) | cut -f1 -d'.')
GO_VERSION_MIN := $(shell echo $(GO_VERSION) | cut -f2 -d'.')

GOFMT ?= gofmt
LN = ln
RM = rm

GOBINPATH    := $(shell $(GO) env GOPATH)/bin
TAGS         := development
PROJECT_PATH := github.com/flavio/terraform-provider-susepubliccloud

CODE_DIRS    = pkg susepubliccloud

# go source files, ignore vendor directory
CODE_SRCS = $(shell find $(CODE_DIRS) -type f -name '*.go')

.PHONY: all
all: build

.PHONY: build
build: go-version-check
	$(GO) build $(GOMODFLAG)

.PHONY: clean
clean:
	$(GO) clean -i

.PHONY: distclean
distclean: clean
	$(GO) clean -i -cache -testcache -modcache

.PHONY: go-version-check
go-version-check:
	@[ $(GO_VERSION_MAJ) -ge 2 ] || \
		[ $(GO_VERSION_MAJ) -eq 1 -a $(GO_VERSION_MIN) -ge 12 ] || (echo "FATAL: Go version should be >= 1.12.x" ; exit 1 ; )

.PHONY: lint
lint:
	# explicitly enable GO111MODULE otherwise go mod will fail
	GO111MODULE=on go mod tidy && GO111MODULE=on go mod vendor && GO111MODULE=on go mod verify
	$(GO) vet ./...
	test -z `$(GOFMT) -l $(CODE_SRCS)` || { $(GOFMT) -d $(CODE_SRCS) && false; }

# tests
.PHONY: test
test:
	$(GO) test $(GOMODFLAG) -coverprofile=coverage.out $(PROJECT_PATH)/{pkg,susepubliccloud}/...

.PHONY: test-unit-coverage
test-coverage: test
	$(GO) tool cover -html=coverage.out
