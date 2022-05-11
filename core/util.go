package core

import (
	"fmt"
	"net/url"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
	models "code.cloudfoundry.org/cli/plugin/models"
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

func MySpace(services *Services) (models.SpaceFields, error) {
	space, err := services.CLI.GetCurrentSpace()
	if err != nil {
		return models.SpaceFields{}, err
	}

	return space.SpaceFields, nil
}

func MyApps(services *Services) ([]models.GetAppsModel, error) {
	return services.CLI.GetApps()
}

func AppByGUID(apps []models.GetAppsModel, guid string) (models.GetAppsModel, error) {
	for _, app := range apps {
		if app.Guid == guid {
			return app, nil
		}
	}

	return models.GetAppsModel{}, fmt.Errorf("Could not find app with GUID %s", guid)
}

type CommandLineContext struct {
	Organization string
	Space        string
	Email        string
}

func GetCurrentContext(services *Services) (*CommandLineContext, error) {
	organization, err := services.CLI.GetCurrentOrg()
	if err != nil {
		return nil, err
	}

	space, err := services.CLI.GetCurrentSpace()
	if err != nil {
		return nil, err
	}

	email, err := services.CLI.UserEmail()
	if err != nil {
		return nil, err
	}

	return &CommandLineContext{
		Organization: organization.Name,
		Space:        space.Name,
		Email:        email,
	}, nil
}

func PrintActionInProgress(services *Services, message string, args ...interface{}) error {
	cliContext, err := GetCurrentContext(services)
	if err != nil {
		return err
	}

	action := fmt.Sprintf(message, args...)
	fmt.Printf("%s in org %s / space %s as %s...\n", action, cliContext.Organization, cliContext.Space, cliContext.Email)
	return nil
}
