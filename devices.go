package tradfricoap

import (
	"sort"
	// "log"
	"fmt"
	// "os"
	// "strconv"
	"strings"
	// "github.com/moroen/canopus"
	// "github.com/urfave/cli"
	"github.com/buger/jsonparser"
)

func trimJSON(json string) string {
	json = strings.Trim(json, "[")
	json = strings.Trim(json, "]")
	return json
}

func GetDevices() ([]TradfriLight, error) {
	payload, err := GetRequest(uri_Devices)
	if err != nil {
		return nil, err
	}

	msg := payload.String()
	msg = strings.Trim(msg, "[")
	msg = strings.Trim(msg, "]")
	result := strings.Split(msg, ",")

	lights := []TradfriLight{}

	for i := range result {
		device, err := GetRequest(fmt.Sprintf("%s/%s", uri_Devices, result[i]))
		if err != nil {
			return nil, err
		}

		var aLight TradfriLight

		aDevice := device.GetBytes()
		// fmt.Println(string(aDevice))

		if _, _, _, err := jsonparser.Get(aDevice, attr_Light_control); err == nil {
			if value, err := jsonparser.GetString(aDevice, attr_name); err == nil {
				aLight.Name = value
			}

			if value, err := jsonparser.GetInt(aDevice, attr_id); err == nil {
				aLight.Id = value
			}

			jsonparser.ArrayEach(aDevice, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				if res, err := jsonparser.GetInt(value, attr_light_state); err == nil {
					aLight.State = func() string {
						if res == 1 {
							return "On"
						} else {
							return "Off"
						}
					}()
				}

				if res, err := jsonparser.GetInt(value, attr_light_dimmer); err == nil {
					aLight.Dimmer = res
				}

			}, attr_Light_control)

			/*
				if value, err := jsonparser.GetString(aDevice, attr_Light_control, "5850"); err == nil {
					aLight.State = string(value)
				} */
			lights = append(lights, aLight)
		}
	}

	sort.Slice(lights, func(i, j int) bool {
		return lights[i].Id < lights[j].Id
	})

	/*
		for i := range lights {
			fmt.Println(lights[i].Describe())
		}
	*/
	return lights, err
}
