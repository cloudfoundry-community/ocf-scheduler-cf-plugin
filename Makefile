# This is how we want to name the binary output
BINARY=ocf-scheduler-cf-plugin
VERSION=`./scripts/genver`
PACKAGE ?= github.com/cloudfoundry-community/ocf-scheduler-cf-plugin
TARGET="builds/${BINARY}-${VERSION}"
PREFIX="${TARGET}/${BINARY}-${VERSION}"

GO_LDFLAGS := -ldflags="-X main.Version=$(VERSION)"

REMOTE_HOST := $(shell cat .ssh-remote)
REMOTE_FOLDER := ~/programs/ocf-scheduler/ocf-scheduler-cf-plugin

build:
	go build $(GO_LDFLAGS) .

release: distclean distbuild linux darwin windows

distbuild:
	mkdir -p ${TARGET}

distclean:
	rm -rf ${TARGET}

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${GO_LDFLAGS} -o ${TARGET}/${BINARY}-linux-amd64 ${PACKAGE}
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build ${GO_LDFLAGS} -o ${TARGET}/${BINARY}-linux-arm64 ${PACKAGE}

darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build ${GO_LDFLAGS} -o ${TARGET}/${BINARY}-macos-amd64 ${PACKAGE}
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build ${GO_LDFLAGS} -o ${TARGET}/${BINARY}-macos-arm64 ${PACKAGE}
	
windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build ${GO_LDFLAGS} -o ${TARGET}/${BINARY}-windows-amd64.exe ${PACKAGE}
	


docker-build:
	docker build -t ocf-scheduler-cf-plugin .

	mkdir -p build-output

	docker container create --name build ocf-scheduler-cf-plugin
	docker container cp build:/bin/ocf-scheduler-cf-plugin ./build-output
	docker container rm build

install-remote:
	scp build-output/ocf-scheduler-cf-plugin $(REMOTE_HOST):$(REMOTE_FOLDER)/ocf-scheduler-cf-plugin
	ssh $(REMOTE_HOST) "cd $(REMOTE_FOLDER); cf uninstall-plugin OCFScheduler || true; yes | cf install-plugin ocf-scheduler-cf-plugin"

install:
	cf uninstall-plugin OCFScheduler || true
	yes | cf install-plugin ocf-scheduler-cf-plugin

acceptance-tests:
	go test -timeout 600s ./...

all: build install

run-remote: docker-build install-remote
