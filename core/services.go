package core

import "code.cloudfoundry.org/cli/v8/plugin"
import "code.cloudfoundry.org/cli/v8/cf/terminal"

type Services struct {
	CLI    plugin.CliConnection
	Client *Driver
	UI     terminal.UI
}
