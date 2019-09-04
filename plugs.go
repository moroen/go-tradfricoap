package tradfricoap

import (
	"fmt"

	"github.com/bradfitz/slice"
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

func (p *TradfriPlug) getInfo(aDevice []byte) error {

	if value, err := jsonparser.GetString(aDevice, attr_name); err == nil {
		p.Name = value
	}

	if value, err := jsonparser.GetInt(aDevice, attr_id); err == nil {
		p.Id = value
	}

	jsonparser.ArrayEach(aDevice, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if res, err := jsonparser.GetInt(value, attr_plug_state); err == nil {
			p.State = func() string {
				if res == 1 {
					return "On"
				} else {
					return "Off"
				}
			}()
		}
	}, attr_plug_control)

	if value, err := jsonparser.GetString(aDevice, attr_DeviceInfo, attr_DeviceInfo_Model); err == nil {
		p.Model = value
	}

	return nil
}

func GetPlug(id int64) (TradfriPlug, error) {
	var aPlug TradfriPlug

	msg, err := GetRequest(fmt.Sprintf("%s/%d", uri_Devices, id))
	if err != nil {
		return aPlug, err
	}

	if _, _, _, err := jsonparser.Get(msg.Payload, attr_plug_control); err == nil {
		err = aPlug.getInfo(msg.Payload)
		return aPlug, err
	} else {
		return aPlug, fmt.Errorf("device %d is not a plug", id)
	}
}

func GetPlugs() (TradfriPlugs, error) {
	result, err := GetRequest(uri_Devices)
	if err != nil {
		// fmt.Println(err.Error())
		return nil, err
	}

	msg := result.Payload

	plugs := []TradfriPlug{}

	jsonparser.ArrayEach(msg, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if res, err := jsonparser.GetInt(value); err == nil {
			aPlug, err := GetPlug(res)
			if err == nil {
				plugs = append(plugs, aPlug)
			}
		}
	})

	slice.Sort(plugs, func(i, j int) bool {
		return plugs[i].Id < plugs[j].Id
	})

	return plugs, err
}
