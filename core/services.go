package core

import "code.cloudfoundry.org/cli/plugin"

type Services struct {
	CLI    plugin.CliConnection
	Client *Driver
}
