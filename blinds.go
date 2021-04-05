package tradfricoap

import (
	"fmt"
	"log"

	// "sort"

	"github.com/buger/jsonparser"
)

type TradfriBlind struct {
	Id    int64
	Name  string
	State float64
	Model string
}

type TradfriBlinds []TradfriBlind

func (p TradfriBlind) Describe() string {
	return fmt.Sprintf("%d: %s (%s) - %.1f", p.Id, p.Name, p.Model, p.State)
}

func getBlindInfo(aDevice []byte) (TradfriBlind, error) {
	var p TradfriBlind

	if value, err := jsonparser.GetString(aDevice, attrName); err == nil {
		p.Name = value
	}

	if value, err := jsonparser.GetInt(aDevice, attrId); err == nil {
		p.Id = value
	}

	_, err := jsonparser.ArrayEach(aDevice, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if res, err := jsonparser.GetFloat(value, attrBlindPosition); err == nil {
			p.State = res
		}
	}, attrBlindControl)
	if err != nil {
		return p, err
	}

	if value, err := jsonparser.GetString(aDevice, attrDeviceInfo, attrDeviceInfo_Model); err == nil {
		p.Model = value
	}
	return p, err
}

func getBlind(id int64) (TradfriBlind, error) {
	var aBlind TradfriBlind

	device, err := GetRequest(fmt.Sprintf("%s/%d", uriDevices, id))
	if err != nil {
		return aBlind, err
	}

	aDevice := device

	if _, _, _, err := jsonparser.Get(aDevice, attrBlindControl); err == nil {
		aBlind, err := getBlindInfo(aDevice)
		return aBlind, err
	} else {
		return aBlind, fmt.Errorf("device %d is not a blind", id)
	}
}

func SetBlind(id int64, level int) (TradfriBlind, error) {
	// Blinds

	uri := fmt.Sprintf("%s/%d", uriDevices, id)
	payload := fmt.Sprintf("{ \"%s\": [{ \"%s\": %d}] }", attrBlindControl, attrBlindPosition, level)
	_, err := PutRequest(uri, payload)

	if err != nil {
		log.Println(err.Error())
		return TradfriBlind{}, err
	}
	return getBlind((id))
}
