package core

import (
	"fmt"
	"net/url"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

func GetBearer(cli plugin.CliConnection) (string, error) {
	good, err := cli.IsLoggedIn()
	if err != nil || !good {
		return "", fmt.Errorf("You must be logged in.")
	}

	raw, err := cli.AccessToken()
	if err != nil {
		return "", fmt.Errorf("You do not have an access token.")
	}

	parts := strings.Split(raw, " ")

	// we have to normalize this until the scheduler API is updated to
	// be less picky
	if strings.ToLower(parts[0]) == "bearer" {
		parts[0] = "Bearer"
	}

	return strings.Join(parts, " "), nil
}

func GetScheduler(cli plugin.CliConnection) (string, error) {
	good, err := cli.HasAPIEndpoint()
	if err != nil || !good {
		return "", fmt.Errorf("Your client is not configured for the CF API.")
	}

	api, err := cli.ApiEndpoint()
	if err != nil {
		return "", fmt.Errorf("Could not get the CF API endpoint.")
	}

	u, err := url.Parse(api)
	if err != nil {
		return "", fmt.Errorf("Could not calculate the Scheduler API endpoint.")
	}

	parts := strings.Split(u.Host, ".")
	parts[0] = "scheduler"
	u.Host = strings.Join(parts, ".")

	return u.String(), nil
}
