package core

import "code.cloudfoundry.org/cli/plugin"
import "code.cloudfoundry.org/cli/cf/terminal"

type Services struct {
	CLI    plugin.CliConnection
	Client *Driver
    UI     terminal.UI
}
