package tradfricoap

import (
	"fmt"
	"sort"

	"github.com/buger/jsonparser"
)

type TradfriPlug struct {
	Id    int64
	Name  string
	State string
	Model string
}

type TradfriPlugs []TradfriPlug

func (p TradfriPlug) Describe() string {
	return fmt.Sprintf("%d: %s (%s) - %s", p.Id, p.Name, p.Model, p.State)
}

func getPlugInfo(aDevice []byte) (TradfriPlug, error) {
	var p TradfriPlug

	if value, err := jsonparser.GetString(aDevice, attrName); err == nil {
		p.Name = value
	}

	if value, err := jsonparser.GetInt(aDevice, attrId); err == nil {
		p.Id = value
	}

	_, err := jsonparser.ArrayEach(aDevice, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if res, err := jsonparser.GetInt(value, attrPlugState); err == nil {
			p.State = func() string {
				if res == 1 {
					return "On"
				} else {
					return "Off"
				}
			}()
		}
	}, attrPlugControl)
	if err != nil {
		return p, err
	}

	if value, err := jsonparser.GetString(aDevice, attrDeviceInfo, attrDeviceInfo_Model); err == nil {
		p.Model = value
	}
	return p, err
}

func GetPlug(id int64) (TradfriPlug, error) {
	var aPlug TradfriPlug

	msg, err := GetRequest(fmt.Sprintf("%s/%d", uriDevices, id))
	if err != nil {
		return aPlug, err
	}

	if _, _, _, err := jsonparser.Get(msg, attrPlugControl); err == nil {
		aPlug, err := getPlugInfo((msg))
		return aPlug, err
	} else {
		return aPlug, fmt.Errorf("device %d is not a plug", id)
	}
}

func GetPlugs() (TradfriPlugs, error) {
	result, err := GetRequest(uriDevices)
	if err != nil {
		// fmt.Println(err.Error())
		return nil, err
	}

	msg := result

	plugs := []TradfriPlug{}

	_, err = jsonparser.ArrayEach(msg, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if res, err := jsonparser.GetInt(value); err == nil {
			aPlug, err := GetPlug(res)
			if err == nil {
				plugs = append(plugs, aPlug)
			}
		}
	})
	if err != nil {
		panic(err.Error())
	}

	sort.Slice(plugs, func(i, j int) bool {
		return plugs[i].Id < plugs[j].Id
	})

	return plugs, err
}
