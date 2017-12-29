package tradfricoap

import (
	"sort"
	// "log"

	// "os"
	// "strconv"
	"strings"
	// "github.com/moroen/canopus"
	// "github.com/urfave/cli"
)

func trimJSON(json string) string {
	json = strings.Trim(json, "[")
	json = strings.Trim(json, "]")
	return json
}

func GetDevices() (TradfriLights, TradfriGroups, error) {
	payload, err := GetRequest(uri_Devices)
	if err != nil {
		return nil, nil, err
	}

	msg := payload.String()
	msg = strings.Trim(msg, "[")
	msg = strings.Trim(msg, "]")
	result := strings.Split(msg, ",")

	lights := []TradfriLight{}

	for i := range result {
		aLight, err := GetLight(result[i])
		if err == nil {
			lights = append(lights, aLight)
		}
	}

	sort.Slice(lights, func(i, j int) bool {
		return lights[i].Id < lights[j].Id
	})

	return lights, nil, err
}
