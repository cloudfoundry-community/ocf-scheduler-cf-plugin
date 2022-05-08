package core

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
	models "code.cloudfoundry.org/cli/plugin/models"
	"github.com/ess/hype"
	scheduler "github.com/starkandwayne/scheduler-for-ocf/core"
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

func AllJobs(services *Services, space models.SpaceFields) ([]*scheduler.Job, error) {
	output := make([]*scheduler.Job, 0)

	params := hype.Params{}
	params.Set("space_guid", space.Guid)

	response := services.Client.Get("jobs", params)
	if !response.Okay() {
		return output, response.Error()
	}

	data := struct {
		Resources []*scheduler.Job `json:"resources"`
	}{}

	err := json.Unmarshal(response.Data(), &data)
	if err != nil {
		return output, err
	}

	return data.Resources, nil
}

func JobNamed(services *Services, space models.SpaceFields, name string) (*scheduler.Job, error) {
	all, err := AllJobs(services, space)
	if err != nil {
		return nil, err
	}

	for _, candidate := range all {
		if candidate.Name == name {
			return candidate, nil
		}
	}

	return nil, fmt.Errorf("could not find job named %s", name)
}

func ExecuteJob(services *Services, job *scheduler.Job) error {
	response := services.Client.Post("jobs/"+job.GUID+"/execute", nil, nil)

	if !response.Okay() {
		return response.Error()
	}

	return nil
}

type schedulePayload struct {
	Enabled        bool   `json:"enabled"`
	Expression     string `json:"expression"`
	ExpressionType string `json:"expression_type"`
}

func ScheduleJob(services *Services, job *scheduler.Job, expression string) (*scheduler.Schedule, error) {
	schedule := &schedulePayload{
		Enabled:        true,
		Expression:     expression,
		ExpressionType: "cron_expression",
	}

	sdata, err := json.Marshal(schedule)
	if err != nil {
		return nil, err
	}

	response := services.Client.Post("jobs/"+job.GUID+"/schedules", nil, sdata)

	if !response.Okay() {
		return nil, response.Error()
	}

	output := &scheduler.Schedule{}

	err = json.Unmarshal(response.Data(), output)
	if err != nil {
		return nil, err
	}

	return output, nil
}
