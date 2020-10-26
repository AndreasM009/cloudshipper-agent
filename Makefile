################################################################################
# Variables
################################################################################
export GO111MODULE ?= on
export GOPROXY ?= https://proxy.golang.org
export GOSUMDB ?= sum.golang.org
# By default, disable CGO_ENABLED. See the details on https://golang.org/cmd/cgo
CGO         ?= 0
BINARIES ?= controller runner
CONTROLLER_BINARY ?= controller
RUNNER_BINARY ?= runner 

################################################################################
# Git info
################################################################################
GIT_COMMIT  = $(shell git rev-list -1 HEAD)
GIT_VERSION = $(shell git describe --always --abbrev=7 --dirty)

################################################################################
# Release version
################################################################################
LASTEST_VERSION_TAG ?=

ifdef REL_VERSION
	AGENT_VERSION := $(REL_VERSION)
	DOCKER_TAG := $(REL_VERSION)
	DOCKER_LATEST_TAG := latest
else
	AGENT_VERSION := edge
	DOCKER_TAG := edge
	DOCKER_LATEST_TAG := edge
endif

CONTROLLER_IMAGE_NAME := cs-agent-controller
RUNNER_IMAGE_NAME := cs-agent-runner

################################################################################
# Architectue
################################################################################
LOCAL_ARCH := $(shell uname -m)
ifeq ($(LOCAL_ARCH),x86_64)
	TARGET_ARCH_LOCAL=amd64
else ifeq ($(shell echo $(LOCAL_ARCH) | head -c 5),armv8)
	TARGET_ARCH_LOCAL=arm64
else ifeq ($(shell echo $(LOCAL_ARCH) | head -c 4),armv)
	TARGET_ARCH_LOCAL=arm
else
	TARGET_ARCH_LOCAL=amd64
endif
export GOARCH ?= $(TARGET_ARCH_LOCAL)

################################################################################
# OS
################################################################################
LOCAL_OS := $(shell uname)
ifeq ($(LOCAL_OS),Linux)
   TARGET_OS_LOCAL = linux
else ifeq ($(LOCAL_OS),Darwin)
   TARGET_OS_LOCAL = darwin
else
   TARGET_OS_LOCAL ?= windows
endif
export GOOS ?= $(TARGET_OS_LOCAL)

################################################################################
# Binaries extension
################################################################################
ifeq ($(GOOS),windows)
BINARY_EXT_LOCAL:=.exe
GOLANGCI_LINT:=golangci-lint.exe
else
BINARY_EXT_LOCAL:=
GOLANGCI_LINT:=golangci-lint
endif

export BINARY_EXT ?= $(BINARY_EXT_LOCAL)

################################################################################
# GO build flags
################################################################################
BASE_PACKAGE_NAME := github.com/AndreasM009/cloudshipper-agent

DEFAULT_LDFLAGS := -X $(BASE_PACKAGE_NAME)/pkg/version.commit=$(GIT_VERSION) -X $(BASE_PACKAGE_NAME)/pkg/version.version=$(AGENT_VERSION)
ifeq ($(DEBUG),)
  BUILDTYPE_DIR:=release
  LDFLAGS:="$(DEFAULT_LDFLAGS) -s -w"
else ifeq ($(DEBUG),0)
  BUILDTYPE_DIR:=release
  LDFLAGS:="$(DEFAULT_LDFLAGS) -s -w"
else
  BUILDTYPE_DIR:=debug
  GCFLAGS:=-gcflags="all=-N -l"
  LDFLAGS:="$(DEFAULT_LDFLAGS)"
  $(info Build with debugger information)
endif

################################################################################
# output directory
################################################################################
OUT_DIR := ./dist
AGENT_OUT_DIR := $(OUT_DIR)/$(GOOS)_$(GOARCH)/$(BUILDTYPE_DIR)
AGENT_LINUX_OUT_DIR := $(OUT_DIR)/linux_$(GOARCH)/$(BUILDTYPE_DIR)

################################################################################
# Target: build-all                                                               
################################################################################
.PHONY: build-all
AGENT_BINS:=$(foreach ITEM,$(BINARIES),$(AGENT_OUT_DIR)/$(ITEM)$(BINARY_EXT))
build-all: $(AGENT_BINS)

# Generate builds for agent binaries for the target
# Params:
# $(1): the binary name for the target
# $(2): the binary main directory
# $(3): the target os
# $(4): the target arch
# $(5): the output directory
define genBinariesForTarget
.PHONY: $(5)/$(1)
$(5)/$(1):
	CGO_ENABLED=$(CGO) GOOS=$(3) GOARCH=$(4) go build $(GCFLAGS) -ldflags=$(LDFLAGS) \
	-o $(5)/$(1) \
	$(2)/main.go;
endef

# Generate binary targets
$(foreach ITEM,$(BINARIES),$(eval $(call genBinariesForTarget,$(ITEM)$(BINARY_EXT),./cmd/$(ITEM),$(GOOS),$(GOARCH),$(AGENT_OUT_DIR))))

################################################################################
# Target: build-controller                                                              
################################################################################
.PHONY: build-controller
CONTROLLER_BIN_EXT:=$(AGENT_OUT_DIR)/$(CONTROLLER_BINARY)$(BINARY_EXT)
build-controller: $(CONTROLLER_BIN_EXT)

CGO_ENABLED=$(CGO) GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(GCFLAGS) -ldflags=$(LDFLAGS) \
-o $(AGENT_OUT_DIR/$(CONTROLLER_BIN_EXT) ./cmd/$(CONTROLLER_BINARY)

################################################################################
# Target: build-runner                                                              
################################################################################
.PHONY: build-runner
RUNNER_BIN_EXT:=$(AGENT_OUT_DIR)/$(RUNNER_BINARY)$(BINARY_EXT)
build-runner: $(RUNNER_BIN_EXT)

CGO_ENABLED=$(CGO) GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(GCFLAGS) -ldflags=$(LDFLAGS) \
-o $(AGENT_OUT_DIR/$(RUNNER_BIN_EXT) ./cmd/$(RUNNER_BINARY)

################################################################################
# Target: lint                                                                
################################################################################
.PHONY: lint	
lint:
	$(GOLANGCI_LINT) run --fix

################################################################################
# Target: test
################################################################################
.PHONY: test
test:
	go test ./pkg/...

################################################################################
# Target: docker-build-...
################################################################################
.PHONY: docker-build-runner docker-build-controller docker-build-all

docker-build-runner:
	docker build -t $(RUNNER_IMAGE_NAME):$(DOCKER_TAG) -f ./Docker/runner/Dockerfile .

docker-build-controller:
	docker build -t $(CONTROLLER_IMAGE_NAME):$(DOCKER_TAG) -f ./Docker/controller/Dockerfile .

docker-build-all: docker-build-controller docker-build-runner

################################################################################
# Target: docker-publish
################################################################################
check-docker-publish-args:
ifeq ($(s),)
	$(error docker server must be set: s=<dockerserver>)
endif
ifeq ($(u),)
	$(error docker login must be set: u=<dockerusername>)
endif
ifeq ($(p),)
	$(error docker password must be set: p=<dockerpassword>)
endif

.PHONY: docker-publish-runner
docker-publish-runner: check-docker-publish-args
	docker login -p $(p) -u $(u)
	docker build -t $(s)/$(RUNNER_IMAGE_NAME):$(DOCKER_TAG) -f ./Docker/runner/Dockerfile .
	docker tag $(s)/$(RUNNER_IMAGE_NAME):$(DOCKER_TAG) $(s)/$(RUNNER_IMAGE_NAME):$(DOCKER_LATEST_TAG)
	docker push $(s)/$(RUNNER_IMAGE_NAME):$(DOCKER_TAG)
	docker push $(s)/$(RUNNER_IMAGE_NAME):$(DOCKER_LATEST_TAG)

.PHONY: docker-publish-controller
docker-publish-controller: check-docker-publish-args
	docker login -p $(p) -u $(u)
	docker build -t $(s)/$(CONTROLLER_IMAGE_NAME):$(DOCKER_TAG) -f ./Docker/controller/Dockerfile .
	docker tag $(s)/$(CONTROLLER_IMAGE_NAME):$(DOCKER_TAG) $(s)/$(CONTROLLER_IMAGE_NAME):$(DOCKER_LATEST_TAG)
	docker push $(s)/$(CONTROLLER_IMAGE_NAME):$(DOCKER_TAG)
	docker push $(s)/$(CONTROLLER_IMAGE_NAME):$(DOCKER_LATEST_TAG)

.PHONY: docker-publish-all
docker-publish-all: docker-publish-controller docker-publish-runner