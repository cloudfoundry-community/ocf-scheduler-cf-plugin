PROJECT        :=ocf-scheduler-cf-plugin
SHELL          :=/bin/bash
GOOS           :=$(shell go env GOOS)
GOARCH         :=$(shell go env GOARCH)
GOMODULECMD    :=main
RELEASE_ROOT   ?=releases
DEV_TEST_BUILD =./$(PROJECT)
TARGETS        ?=linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

# Version resolution lives in the vendored snippet. VERSION accepts the full
# semver grammar -- optional v/v- prefix, -prerelease, +buildmeta -- which the
# pipeline depends on: the semver resource hands the build job values like
# 1.2.1-rc.1. Unset, it resolves the nearest tag with the patch incremented.
include mk/snippets/version.mk

# Lifecycle verbs for the vendored tree: snippets-status (has upstream moved,
# have we edited locally), snippets-report, snippets-merge, snippets-patch,
# snippets-pin. Provenance lives in mk/snippets/snippets.sha256.
include mk/snippets/snippets.mk

# Artifact names carry major.minor.patch[-prerelease] only. No tag prefix, and
# build metadata stays out: the +os.arch segment already occupies that slot,
# and SEMVER_BUILDMETA is appended after it by the release rule below.
RELEASE_VERSION = $($(_HIDE)SEMVER_MAJOR).$($(_HIDE)SEMVER_MINOR).$($(_HIDE)SEMVER_PATCH)$(if $($(_HIDE)SEMVER_PRERELEASE),-$($(_HIDE)SEMVER_PRERELEASE))

build: $(_HIDE)SEMVER_PRERELEASE := dev

# Lazy (=) so the target-specific dev prerelease above reaches the linker.
GO_LDFLAGS = -X '$(GOMODULECMD).SemVerMajor=$($(_HIDE)SEMVER_MAJOR)' \
	         -X '$(GOMODULECMD).SemVerMinor=$($(_HIDE)SEMVER_MINOR)' \
	         -X '$(GOMODULECMD).SemVerPatch=$($(_HIDE)SEMVER_PATCH)' \
	         -X '$(GOMODULECMD).SemVerPrerelease=$($(_HIDE)SEMVER_PRERELEASE)' \
	         -X '$(GOMODULECMD).SemVerBuild=$($(_HIDE)SEMVER_BUILDMETA)' \
	         -X '$(GOMODULECMD).BuildDate=$($(_HIDE)BUILD_DATE)' \
	         -X '$(GOMODULECMD).BuildVcsUrl=$($(_HIDE)BUILD_VCS_URL)' \
	         -X '$(GOMODULECMD).BuildVcsId=$($(_HIDE)BUILD_VCS_ID)' \
		     -X '$(GOMODULECMD).BuildVcsIdDate=$($(_HIDE)BUILD_VCS_ID_DATE)'

# GMake rules generally used for local development

.PHONY: build clean install acceptance-tests

build: BUILD_GO_LDFLAGS:=-ldflags="$(GO_LDFLAGS) -X '$(GOMODULECMD).GoOs=$(GOOS)' -X '$(GOMODULECMD).GoArch=$(GOARCH)'"

build: BUILD_RULE_CMD := CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) \
	                     go build $(BUILD_GO_LDFLAGS) -o $(DEV_TEST_BUILD)

build: clean
	@echo "Building $(DEV_TEST_BUILD)"
	$(BUILD_RULE_CMD)

clean:
	@rm -f $(DEV_TEST_BUILD) || true

install: build
	cf install-plugin $(DEV_TEST_BUILD) -f || true

acceptance-tests:
	go test -timeout 600s ./...

# GMake rules for release building are below

.PHONY: distbuild require-% release-% ci-release check-version clean distclean show-releases

require-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "Environment variable $* not set"; \
		exit 1; \
	fi

RELEASES := $(foreach target,$(TARGETS),release-$(target)-$(PROJECT))

show-releases:
	@ls -lA $(RELEASE_ROOT)
	@echo ""

# version.mk soft-fails an unusable version to 0.0.0-unknown instead of
# erroring, so the release path asks for the hard failure the old in-tree
# parser used to give. Prerelease and build metadata are not judged here --
# SEMVER_VALID covers major.minor.patch, which is what must be numeric.
check-version:
	$(if $($(_HIDE)SEMVER_VALID),,$(error VERSION '$(VERSION)' is not usable))

# require-VERSION is no longer sufficient here: version.mk always resolves a
# VERSION from the nearest tag, so the variable is never empty and the guard
# could never trip. Releasing must stay deliberate, so test where the value
# came from rather than whether it exists.
ci-release:
	@if [ -z "$(filter command line environment,$(origin VERSION))" ]; then \
		echo "VERSION must be set explicitly: make ci-release VERSION=x.y.z" >&2; \
		exit 1; \
	fi
	@$(MAKE) check-version release-all

release-all: release-clean distbuild $(RELEASES) show-releases

distbuild:
	@mkdir -p $(RELEASE_ROOT)

# Arguments os,arch,build 
define build-target
release-$(1)/$(2)-$(PROJECT): RELEASE_GO_LDFLAGS:=-ldflags="$(GO_LDFLAGS) -X '$(GOMODULECMD).GoOs=$(1)' -X '$(GOMODULECMD).GoArch=$(2)'"

release-$(1)/$(2)-$(PROJECT): RELEASE_EXECUTABLE_BASE:=$(RELEASE_ROOT)/$(PROJECT)-$(RELEASE_VERSION)+$(1).$(2)$(if $(3),.$(3))

release-$(1)/$(2)-$(PROJECT): RELEASE_EXECUTABLE:=$$(RELEASE_EXECUTABLE_BASE)$(if $(patsubst windows,,$(1)),,.exe)

release-$(1)/$(2)-$(PROJECT): RELEASE_EXECUTABLE_SHA1:=$$(RELEASE_EXECUTABLE_BASE).sha1

release-$(1)/$(2)-$(PROJECT):
	@echo "Building $$(PROJECT) version $$(RELEASE_VERSION) for $(1) $(2) ..."
	@CGO_ENABLED=0 GOOS=$(1) GOARCH=$(2) go build -o $$(RELEASE_EXECUTABLE) $$(RELEASE_GO_LDFLAGS)
	@openssl sha1 -r $$(RELEASE_EXECUTABLE) > $$(RELEASE_EXECUTABLE_SHA1)
endef

$(foreach target,$(TARGETS), $(eval $(call build-target,$(word 1, $(subst /, ,$(target))),$(word 2, $(subst /, ,$(target))),$($(_HIDE)SEMVER_BUILDMETA))))

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
