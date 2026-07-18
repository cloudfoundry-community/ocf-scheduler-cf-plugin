package core

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	models "code.cloudfoundry.org/cli/plugin/models"
)

// CliCurler is the subset of plugin.CliConnection needed to issue `cf curl`
// requests. The plugin RPC GetApp/GetApps helpers are backed by the CC v2
// API, which modern foundations no longer serve, so app lookups must go
// through the v3 API instead.
type CliCurler interface {
	CliCommandWithoutTerminalOutput(args ...string) ([]string, error)
}

type appsV3Response struct {
	Pagination struct {
		Next *struct {
			Href string `json:"href"`
		} `json:"next"`
	} `json:"pagination"`
	Resources []struct {
		Guid string `json:"guid"`
		Name string `json:"name"`
	} `json:"resources"`
	Errors []struct {
		Title  string `json:"title"`
		Detail string `json:"detail"`
	} `json:"errors"`
}

// AppsV3 lists all apps in the given space via the CC v3 API.
func AppsV3(cli CliCurler, spaceGUID string) ([]models.GetAppsModel, error) {
	path := fmt.Sprintf("/v3/apps?space_guids=%s&per_page=5000", url.QueryEscape(spaceGUID))
	return appsV3Fetch(cli, path)
}

// AppV3ByName resolves a single app by name in the given space via the CC
// v3 API.
func AppV3ByName(cli CliCurler, spaceGUID, name string) (models.GetAppsModel, error) {
	path := fmt.Sprintf("/v3/apps?names=%s&space_guids=%s&per_page=5000",
		url.QueryEscape(name), url.QueryEscape(spaceGUID))

	apps, err := appsV3Fetch(cli, path)
	if err != nil {
		return models.GetAppsModel{}, err
	}
	for _, app := range apps {
		if app.Name == name {
			return app, nil
		}
	}

	return models.GetAppsModel{}, fmt.Errorf("could not find app with name %s", name)
}

func appsV3Fetch(cli CliCurler, path string) ([]models.GetAppsModel, error) {
	apps := []models.GetAppsModel{}

	for path != "" {
		lines, err := cli.CliCommandWithoutTerminalOutput("curl", path)
		if err != nil {
			return nil, err
		}

		var parsed appsV3Response
		if err := json.Unmarshal([]byte(strings.Join(lines, "")), &parsed); err != nil {
			return nil, fmt.Errorf("could not parse CC /v3/apps response: %s", err)
		}
		if len(parsed.Errors) > 0 {
			return nil, fmt.Errorf("CC API error: %s: %s", parsed.Errors[0].Title, parsed.Errors[0].Detail)
		}

		for _, resource := range parsed.Resources {
			apps = append(apps, models.GetAppsModel{Guid: resource.Guid, Name: resource.Name})
		}

		path = ""
		if parsed.Pagination.Next != nil {
			next, err := url.Parse(parsed.Pagination.Next.Href)
			if err != nil {
				return nil, fmt.Errorf("could not parse CC pagination link: %s", err)
			}
			path = next.Path
			if next.RawQuery != "" {
				path += "?" + next.RawQuery
			}
		}
	}

	return apps, nil
}
