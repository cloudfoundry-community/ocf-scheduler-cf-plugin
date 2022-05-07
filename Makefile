GO_LDFLAGS := -ldflags="-X main.Version=$(VERSION)"

build:
	go build $(GO_LDFLAGS) .

install:
	cf uninstall-plugin ocf-scheduler || true
	yes | cf install-plugin cf-plugin-*

all: build install
