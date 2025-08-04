PROJECT_NAME := shoehorn
VERSION ?= 0.0.1
GOBIN := $(CURDIR)/bin
BUILD_DIR := $(CURDIR)/build
SOURCES := $(shell find . -name '*.go' -print0 | xargs -0)

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development
.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: gotestsum gocover-cobertura ## Run tests.
	$(GOBIN)/gotestsum --junitfile report.xml --format testname -- -coverprofile=coverage.out.tmp ./...
	grep -v "mocks/" coverage.out.tmp > coverage.out
	$(GOBIN)/gocover-cobertura < coverage.out > coverage.xml

##@ Build

.PHONY: build
build: $(BUILD_DIR)/$(PROJECT_NAME) ## Build manager binary.

.PHONY: lint
lint: golangci-lint ## Run golangci-lint against code.
	$(GOBIN)/golangci-lint run --timeout 5m

$(BUILD_DIR)/$(PROJECT_NAME): $(BUILD_DIR) $(SOURCES) ## Create build directory.
	CGO_ENBALED=0 go build -o $(BUILD_DIR)/$(PROJECT_NAME) main.go

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

##@ Tool Dependencies
GOTESTSUM_VERSION ?= v1.12.1
GOCOVER_COBERTURA_VERSION ?= v1.3.0
GOLANGCI_LINT_VERSION ?= v2.3.1

.PHONY: gotestsum
gotestsum: ## Install gotestsum
	GOBIN=$(GOBIN) go install gotest.tools/gotestsum@$(GOTESTSUM_VERSION)

.PHONY: gocover-cobertura
gocover-cobertura: ## Install gocover-cobertura
	GOBIN=$(GOBIN) go install github.com/boumenot/gocover-cobertura@$(GOCOVER_COBERTURA_VERSION)

.PHONY: golangci-lint
golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(GOBIN) $(GOLANGCI_LINT_VERSION)
