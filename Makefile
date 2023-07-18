# This is how we want to name the binary output
APP_NAME ?= scheduler
VERSION ?= `./scripts/genver`-dev

MODULE ?= github.com/cloudfoundry-community/ocf-scheduler-cf-plugin
CMD_PATH ?= .
CGO_ENABLED ?= 0
BUILD_PATH=builds
BUILD=${APP_NAME}-${VERSION}
PACKAGE=${MODULE}/${CMD_PATH}

GO_LDFLAGS := -ldflags="-X main.Version=$(VERSION)"

REMOTE_HOST := $(shell [[ -f .ssh-remote ]] && cat .ssh-remote || echo '')
REMOTE_FOLDER := ~/programs/ocf-scheduler-cf-plugin/ocf-scheduler-cf-plugin

build: clean
	go build $(GO_LDFLAGS) .

# Cleans our project: deletes binaries
clean:
	rm -rf ${APP_NAME}

# Releases
release: distclean distbuild linux darwin windows

distbuild:
	mkdir -p ${BUILD_PATH}

distclean:
	rm -rf ${BUILD_PATH}

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${GO_LDFLAGS} -o ${BUILD_PATH}/${BUILD}-linux-amd64 ${PACKAGE}
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build ${GO_LDFLAGS} -o ${BUILD_PATH}/${BUILD}-linux-arm64 ${PACKAGE}

darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build ${GO_LDFLAGS} -o ${BUILD_PATH}/${BUILD}-macos-amd64 ${PACKAGE}
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build ${GO_LDFLAGS} -o ${BUILD_PATH}/${BUILD}-macos-arm64 ${PACKAGE}
	
windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build ${GO_LDFLAGS} -o ${BUILD_PATH}/${BUILD}-windows-amd64.exe ${PACKAGE}

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
