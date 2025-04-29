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
lint: golangci-lint
	# explicitly enable GO111MODULE otherwise go mod will fail
	GO111MODULE=on go mod tidy && GO111MODULE=on go mod vendor && GO111MODULE=on go mod verify
	$(GO) vet ./...
	test -z `$(GOFMT) -l $(CODE_SRCS)` || { $(GOFMT) -d $(CODE_SRCS) && false; }
	$(GOLANGCI_LINT) run 

.PHONY: lint-fix
lint-fix: golangci-lint ## Run golangci-lint linter and perform fixes
	$(GOLANGCI_LINT) run --fix

# tests
.PHONY: test
test:
	$(GO) test $(GOMODFLAG) -coverprofile=coverage.out -v ./...

.PHONY: test-unit-coverage
test-coverage: test
	$(GO) tool cover -html=coverage.out


##@ Dependencies

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
GOLANGCI_LINT = $(LOCALBIN)/golangci-lint-$(GOLANGCI_LINT_VERSION)

## Tool Versions
GOLANGCI_LINT_VERSION ?= v2.1.5

.PHONY: golangci-lint
golangci-lint: $(GOLANGCI_LINT) ## Download golangci-lint locally if necessary.
$(GOLANGCI_LINT): $(LOCALBIN)
	$(call go-install-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/v2/cmd/golangci-lint,${GOLANGCI_LINT_VERSION})

# go-install-tool will 'go install' any package with custom target and name of binary, if it doesn't exist
# $1 - target path with name of binary (ideally with version)
# $2 - package url which can be installed
# $3 - specific version of package
define go-install-tool
@[ -f $(1) ] || { \
set -e; \
package=$(2)@$(3) ;\
echo "Downloading $${package}" ;\
GOBIN=$(LOCALBIN) go install $${package} ;\
mv "$$(echo "$(1)" | sed "s/-$(3)$$//")" $(1) ;\
}
endef
