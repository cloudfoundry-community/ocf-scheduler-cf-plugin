GO_LDFLAGS := -ldflags="-X main.Version=$(VERSION)"

REMOTE_HOST := $(shell cat .ssh-remote)
REMOTE_FOLDER := ~/programs/ocf-scheduler/ocf-scheduler-cf-plugin

build:
	go build $(GO_LDFLAGS) .

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

all: build install

run-remote: docker-build install-remote
