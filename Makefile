GO_LDFLAGS := -ldflags="-X main.Version=$(VERSION)"

build:
	go build $(GO_LDFLAGS) .

install:
	cf uninstall-plugin OCFScheduler || true
	yes | cf install-plugin ocf-scheduler-cf-plugin

acceptance-tests:
	go test -timeout 600s ./...

all: build install
