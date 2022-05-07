GO_LDFLAGS := -ldflags="-X main.Version=$(VERSION)"

build:
	go build $(GO_LDFLAGS) .

install:
	cf uninstall-plugin OCFScheduler || true
	yes | cf install-plugin ocf-scheduler-cf-plugin

all: build install
