package tradfricoap

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/bradfitz/slice"
	"github.com/buger/jsonparser"
)

type TradfriLight struct {
	Id     int64
	Name   string
	State  string
	Dimmer int64
	Model  string
	Colors ColorMap
}

type TradfriLights []TradfriLight

func (l TradfriLight) Describe() string {
	return fmt.Sprintf("%d: %s (%s) - %s (%d)", l.Id, l.Name, l.Model, l.State, l.Dimmer)
}

func GetLight(id int64) (TradfriLight, error) {
	var aLight TradfriLight

	device, err := GetRequest(fmt.Sprintf("%s/%d", uri_Devices, id))
	if err != nil {
		return aLight, err
	}

	aDevice := device.Payload

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

		if value, err := jsonparser.GetString(aDevice, attr_DeviceInfo, attr_DeviceInfo_Model); err == nil {
			aLight.Model = value
			if strings.Contains(value, " CWS ") {
				aLight.Colors = cwsMap()
			} else if strings.Contains(value, " WS ") {
				aLight.Colors = cwMap()
			} else {
				aLight.Colors = nil
			}
		}
	} else {
		err := errors.New(fmt.Sprintf("Device %s is not a light.", id))
		return aLight, err
	}
	return aLight, err
}

func GetDevices() (TradfriLights, error) {
	result, err := GetRequest(uri_Devices)
	if err != nil {
		// fmt.Println(err.Error())
		return nil, err
	}

	log.Println("Have result")

	msg := result.Payload

	lights := []TradfriLight{}

	jsonparser.ArrayEach(msg, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if res, err := jsonparser.GetInt(value); err == nil {
			aLight, err := GetLight(res)
			if err == nil {
				lights = append(lights, aLight)
			}
		}
	})

	slice.Sort(lights, func(i, j int) bool {
		return lights[i].Id < lights[j].Id
	})

	return lights, err
}

func SetState(id int64, state int) (TradfriLight, error) {
	device := TradfriLight{}

	uri := fmt.Sprintf("%s/%d", uri_Devices, id)
	payload := fmt.Sprintf("{ \"%s\": [{ \"%s\": %d }] }", attr_Light_control, attr_light_state, state)
	_, err := PutRequest(uri, payload)
	if err != nil {
		return device, err
	}
	return GetLight(id)
}

func SetLevel(id int64, level int) (device TradfriLight, err error) {
	uri := fmt.Sprintf("%s/%d", uri_Devices, id)
	payload := fmt.Sprintf("{ \"%s\": [{ \"%s\": %d }] }", attr_Light_control, attr_light_dimmer, level)
	_, err = PutRequest(uri, payload)
	if err != nil {
		return device, err
	}
	return GetLight(id)
}

func SetHex(id int64, hex string) (device TradfriLight, err error) {
	uri := fmt.Sprintf("%s/%d", uri_Devices, id)
	payload := fmt.Sprintf("{ \"%s\": [{ \"%s\": \"%s\" }] }", attr_Light_control, attr_light_hex, hex)

	_, err = PutRequest(uri, payload)
	if err != nil {
		return device, err
	}
	return GetLight(id)
}

func SetHexForLevel(id int64, level int) (device TradfriLight, err error) {
	device, err = GetLight(id)
	if err != nil {
		return device, err
	}

	if hex, keyExists := device.Colors[level]["Hex"]; keyExists {
		return SetHex(id, hex)
	} else {
		return device, errors.New(fmt.Sprintf("Unknown colorlevel %d", level))
	}
}
