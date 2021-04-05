package tradfricoap

import (
	"fmt"
	"sort"
	"strings"

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

func getLightInfo(aDevice []byte) (TradfriLight, error) {
	var aLight TradfriLight

	if value, err := jsonparser.GetString(aDevice, attrName); err == nil {
		aLight.Name = value
	}

	if value, err := jsonparser.GetInt(aDevice, attrId); err == nil {
		aLight.Id = value
	}

	_, err := jsonparser.ArrayEach(aDevice, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if res, err := jsonparser.GetInt(value, attrLightState); err == nil {
			aLight.State = func() string {
				if res == 1 {
					return "On"
				} else {
					return "Off"
				}
			}()
		}

		if res, err := jsonparser.GetInt(value, attrLightDimmer); err == nil {
			aLight.Dimmer = res
		}

	}, attrLightControl)
	if err != nil {
		return aLight, err
	}

	if value, err := jsonparser.GetString(aDevice, attrDeviceInfo, attrDeviceInfo_Model); err == nil {
		aLight.Model = value
		if strings.Contains(value, " CWS ") {
			aLight.Colors = cwsMap()
		} else if strings.Contains(value, " WS ") {
			aLight.Colors = cwMap()
		} else {
			aLight.Colors = nil
		}
	}

	return aLight, err
}

func GetLight(id int64) (TradfriLight, error) {
	var aLight TradfriLight

	device, err := GetRequest(fmt.Sprintf("%s/%d", uriDevices, id))
	if err != nil {
		return aLight, err
	}

	aDevice := device

	if _, _, _, err := jsonparser.Get(aDevice, attrLightControl); err == nil {
		aLight, err := getLightInfo(aDevice)
		return aLight, err
	} else {
		return aLight, fmt.Errorf("Device %d is not a light.", id)
	}
}

func GetLights() (TradfriLights, error) {
	result, err := GetRequest(uriDevices)
	if err != nil {
		// fmt.Println(err.Error())
		return nil, err
	}

	msg := result

	lights := []TradfriLight{}

	_, err = jsonparser.ArrayEach(msg, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if res, err := jsonparser.GetInt(value); err == nil {
			aLight, err := GetLight(res)
			if err == nil {
				lights = append(lights, aLight)
			}
		}
	})
	if err != nil {
		panic(err.Error())
	}

	sort.Slice(lights, func(i, j int) bool {
		return lights[i].Id < lights[j].Id
	})

	return lights, err
}

func SetState(id int64, state int) (TradfriLight, error) {
	device := TradfriLight{}

	uri := fmt.Sprintf("%s/%d", uriDevices, id)
	payload := fmt.Sprintf("{ \"%s\": [{ \"%s\": %d }] }", attrLightControl, attrLightState, state)

	_, err := PutRequest(uri, payload)
	if err != nil {
		return device, err
	}
	return GetLight(id)
}

func SetLevel(id int64, level int) (device TradfriLight, err error) {
	uri := fmt.Sprintf("%s/%d", uriDevices, id)
	payload := fmt.Sprintf("{ \"%s\": [{ \"%s\": %d, \"%s\": %d }] }", attrLightControl, attrLightDimmer, level, attrTransitionTime, 10)
	_, err = PutRequest(uri, payload)
	if err != nil {
		return device, err
	}
	return GetLight(id)
}

func SetHex(id int64, hex string) (device TradfriLight, err error) {
	uri := fmt.Sprintf("%s/%d", uriDevices, id)
	payload := fmt.Sprintf("{ \"%s\": [{ \"%s\": \"%s\" }] }", attrLightControl, attrLightHex, hex)

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
		return device, fmt.Errorf("Unknown colorlevel %d", level)
	}
}
