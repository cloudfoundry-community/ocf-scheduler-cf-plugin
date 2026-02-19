PROJECT        :=ocf-scheduler-cf-plugin
SHELL          :=/bin/bash
GOOS           :=$(shell go env GOOS)
GOARCH         :=$(shell go env GOARCH)
GOMODULECMD    :=main
RELEASE_ROOT   ?=releases
DEV_TEST_BUILD =./$(PROJECT)
TARGETS        ?=linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

define is_not_number 
$(shell echo ${1} | sed -e 's/[0123456789]//g')
endef

CLEAN_VERSION = $(patsubst v%,%,$(VERSION))
HAS_BUILDMETA := $(findstring +,$(CLEAN_VERSION))
VERSION_BUILDMETA := $(if $(HAS_BUILDMETA),$(lastword $(subst +, ,$(CLEAN_VERSION))),)

VERSION_AND_PRERELEASE := $(firstword $(subst +, ,$(CLEAN_VERSION)))

HAS_PRERELEASE := $(findstring -,$(VERSION_AND_PRERELEASE))
VERSION_ONLY := $(firstword $(subst -, ,$(VERSION_AND_PRERELEASE)))
VERSION_PRERELEASE := $(if $(HAS_PRERELEASE),$(patsubst $(VERSION_ONLY)-%,%,$(VERSION_AND_PRERELEASE)),)

# Output results

ifneq ($(VERSION_ONLY),)
VERSION_SPLIT:=$(subst ., ,$(VERSION_ONLY))
  ifneq ($(words $(VERSION_SPLIT)),3)
    $(error VERSION does not have 3 parts |$(words $(VERSION_SPLIT))|$(VERSION_ONLY)|$(VERSION_SPLIT)|)
  endif
else
VERSION_TAG:=$(shell git describe --tags --abbrev=0 2>/dev/null || echo 0.0.0)
CLEAN_VERSION_TAG = $(patsubst v%,%,$(VERSION_TAG))
VERSION_SPLIT:=$(subst ., ,$(CLEAN_VERSION_TAG))
  ifneq ($(words $(VERSION_SPLIT)),3)
    $(error VERSION_TAG does not have 3 parts |$(words $(VERSION_SPLIT))|$(VERSION_TAG)|$(VERSION_SPLIT)|)
  endif

  ifneq ($(words $(call is_not_number,$(word 3,$(VERSION_SPLIT)))), 0)
    $(error The VERSION_TAG patch version string contain non-numeric characters)
  endif

  VERSION_SPLIT:=$(wordlist 1, 2, $(VERSION_SPLIT)) $(shell echo $$(($(word 3,$(VERSION_SPLIT))+1)))
endif

ifneq ($(words $(call is_not_number,$(VERSION_SPLIT))), 0)
  $(error The version string contain non-numeric characters)
endif

SEMVER_MAJOR    ?=$(word 1,$(VERSION_SPLIT))
SEMVER_MINOR    ?=$(word 2,$(VERSION_SPLIT))
SEMVER_PATCH    ?=$(word 3,$(VERSION_SPLIT))
SEMVER_PRERELEASE ?=$(VERSION_PRERELEASE)
SEMVER_BUILDMETA  ?=$(VERSION_BUILDMETA)
BUILD_DATE        :=$(shell date -u -Iseconds)
BUILD_VCS_URL     :=$(shell git config --get remote.origin.url)
BUILD_VCS_ID      :=$(shell git log -n 1 --date=iso-strict-local --format="%h")
BUILD_VCS_ID_DATE :=$(shell TZ=UTC0 git log -n 1 --date=iso-strict-local --format='%ad')

build: SEMVER_PRERELEASE := dev

GO_LDFLAGS = -X '$(GOMODULECMD).SemVerMajor=$(SEMVER_MAJOR)' \
	         -X '$(GOMODULECMD).SemVerMinor=$(SEMVER_MINOR)' \
	         -X '$(GOMODULECMD).SemVerPatch=$(SEMVER_PATCH)' \
	         -X '$(GOMODULECMD).SemVerPrerelease=$(SEMVER_PRERELEASE)' \
	         -X '$(GOMODULECMD).SemVerBuild=$(SEMVER_BUILDMETA)' \
	         -X '$(GOMODULECMD).BuildDate=$(BUILD_DATE)' \
	         -X '$(GOMODULECMD).BuildVcsUrl=$(BUILD_VCS_URL)' \
	         -X '$(GOMODULECMD).BuildVcsId=$(BUILD_VCS_ID)' \
		     -X '$(GOMODULECMD).BuildVcsIdDate=$(BUILD_VCS_ID_DATE)'

# The build meta data is added when the build is done
#
SEMVER_VERSION := $(if $(SEMVER_MAJOR),$(SEMVER_MAJOR),$(error Missing SEMVER_MAJOR))
SEMVER_VERSION := $(SEMVER_VERSION)$(if $(SEMVER_MINOR),.$(SEMVER_MINOR),$(error Missing SEMVER_MINOR))
SEMVER_VERSION := $(SEMVER_VERSION)$(if $(SEMVER_PATCH),.$(SEMVER_PATCH),$(error Missing SEMVER_PATCH))
SEMVER_VERSION := $(SEMVER_VERSION)$(if $(SEMVER_PRERELEASE),-$(SEMVER_PRERELEASE))

# GMake rules generally used for local development

.PHONY: generate build clean install acceptance-tests

generate:
	go generate ./...

build: BUILD_GO_LDFLAGS:=-ldflags="$(GO_LDFLAGS) -X '$(GOMODULECMD).GoOs=$(GOOS)' -X '$(GOMODULECMD).GoArch=$(GOARCH)'"

build: BUILD_RULE_CMD := CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) \
	                     go build $(BUILD_GO_LDFLAGS) -o $(DEV_TEST_BUILD)

build: clean generate
	@echo "Building $(DEV_TEST_BUILD)"
	$(BUILD_RULE_CMD)

clean:
	@rm -f $(DEV_TEST_BUILD) cron_expression_styles.go || true

install: build
	cf install-plugin $(DEV_TEST_BUILD) -f || true

acceptance-tests:
	go test -timeout 600s ./...

# GMake rules for release building are below

.PHONY: distbuild require-% release-% ci-release clean distclean show-releases

require-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "Environment variable $* not set"; \
		exit 1; \
	fi

RELEASES := $(foreach target,$(TARGETS),release-$(target)-$(PROJECT))

show-releases:
	@ls -lA $(RELEASE_ROOT)
	@echo ""

ci-release: require-VERSION release-all

release-all: release-clean generate distbuild $(RELEASES) show-releases

distbuild:
	@mkdir -p $(RELEASE_ROOT)

# Arguments os,arch,build 
define build-target
release-$(1)/$(2)-$(PROJECT): RELEASE_GO_LDFLAGS:=-ldflags="$(GO_LDFLAGS) -X '$(GOMODULECMD).GoOs=$(1)' -X '$(GOMODULECMD).GoArch=$(2)'"

release-$(1)/$(2)-$(PROJECT): RELEASE_EXECUTABLE_BASE:=$(RELEASE_ROOT)/$(PROJECT)-$(SEMVER_VERSION)+$(1).$(2)$(if $(3),.$(3))

release-$(1)/$(2)-$(PROJECT): RELEASE_EXECUTABLE:=$$(RELEASE_EXECUTABLE_BASE)$(if $(patsubst windows,,$(1)),,.exe)

release-$(1)/$(2)-$(PROJECT): RELEASE_EXECUTABLE_SHA1:=$$(RELEASE_EXECUTABLE_BASE).sha1

release-$(1)/$(2)-$(PROJECT):
	@echo "Building $$(PROJECT) version $$(SEMVER_VERSION) for $(1) $(2) ..."
	@CGO_ENABLED=0 GOOS=$(1) GOARCH=$(2) go build -o $$(RELEASE_EXECUTABLE) $$(RELEASE_GO_LDFLAGS)
	@openssl sha1 -r $$(RELEASE_EXECUTABLE) > $$(RELEASE_EXECUTABLE_SHA1)
endef

$(foreach target,$(TARGETS), $(eval $(call build-target,$(word 1, $(subst /, ,$(target))),$(word 2, $(subst /, ,$(target))),$(SEMVER_BUILDMETA))))

release-clean:
	@rm -f $(RELEASE_ROOT)/$(PROJECT)-* || true
	@[[ ! -d $(RELEASE_ROOT) ]] || rmdir -p $(RELEASE_ROOT)

distclean: clean release-clean

# REMOTE_HOST := $(shell [ -f .ssh-remote ] && cat .ssh-remote || echo '')
# REMOTE_FOLDER := ~/programs/ocf-scheduler-cf-plugin/ocf-scheduler-cf-plugin

# docker-build:
# 	docker build -t ocf-scheduler-cf-plugin .
#
# 	mkdir -p build-output
#
# 	docker container create --name build ocf-scheduler-cf-plugin
# 	docker container cp build:/bin/ocf-scheduler-cf-plugin ./build-output
# 	docker container rm build

# install-remote:
# 	scp build-output/ocf-scheduler-cf-plugin $(REMOTE_HOST):$(REMOTE_FOLDER)/ocf-scheduler-cf-plugin
# 	ssh $(REMOTE_HOST) "cd $(REMOTE_FOLDER); cf uninstall-plugin OCFScheduler || true; yes | cf install-plugin ocf-scheduler-cf-plugin"

# run-remote: docker-build install-remote
#

.DEFAULT_GOAL := ci-release
