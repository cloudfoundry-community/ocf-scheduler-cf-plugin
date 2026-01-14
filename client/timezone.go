package client

import (
	"encoding/json"

	"github.com/cloudfoundry-community/ocf-scheduler-cf-plugin/core"
	scheduler "github.com/cloudfoundry-community/ocf-scheduler/core"
	"github.com/ess/hype"
)

func ListTimeZones(driver *core.Driver) (*scheduler.TimezoneSlice, error) {
	params := hype.Params{}

	response := driver.Get("scheduler-time-zones", params)
	if !response.Okay() {
		return nil, response.Error()
	}

	data := struct {
		Resources *scheduler.TimezoneSlice `json:"resources"`
	}{}

	err := json.Unmarshal(response.Data(), &data)
	if err != nil {
		return nil, err
	}

	return data.Resources, nil
}
