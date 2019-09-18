package tradfricoap

import (
	"fmt"
	"sort"

	// "sort"

	"github.com/buger/jsonparser"
)

type TradfriGroup struct {
	Id     int64
	Name   string
	State  string
	Dimmer int64
}

func (g TradfriGroup) Describe() string {
	return fmt.Sprintf("%d: %s - %s (%d)", g.Id, g.Name, g.State, g.Dimmer)
}

type TradfriGroups []TradfriGroup

func GetGroup(id int64) (aGroup TradfriGroup, err error) {
	msg, err := GetRequest(fmt.Sprintf("%s/%d", uriGroups, id))
	if err != nil {
		return aGroup, err
	}

	// fmt.Println(msg.String())

	currentGroup := msg
	aGroup.Id = id

	if value, err := jsonparser.GetString(currentGroup, attrGroupName); err == nil {
		aGroup.Name = value
	}

	if value, err := jsonparser.GetInt(currentGroup, attrLightState); err == nil {
		if value == 1 {
			aGroup.State = "On"
		} else {
			aGroup.State = "Off"
		}
	}

	if value, err := jsonparser.GetInt(currentGroup, attrLightDimmer); err == nil {
		aGroup.Dimmer = value
	}

	return aGroup, nil
}

func GetGroups() (TradfriGroups, error) {
	groups := []TradfriGroup{}

	result, err := GetRequest(uriGroups)
	if err != nil {
		return nil, err
	}

	msg := result

	_, err = jsonparser.ArrayEach(msg, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if res, err := jsonparser.GetInt(value); err == nil {
			aGroup, err := GetGroup(res)
			if err == nil {
				groups = append(groups, aGroup)
			}
		}
	})
	if err != nil {
		panic(err.Error())
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Id < groups[j].Id
	})

	return groups, err
}

func GroupSetState(id int64, state int) (group TradfriGroup, err error) {
	uri := fmt.Sprintf("%s/%d", uriGroups, id)
	payload := fmt.Sprintf("{\"%s\":%d}", attrLightState, state)
	// fmt.Println(uri, payload)
	_, err = PutRequest(uri, payload)
	if err != nil {
		return group, err
	}
	return GetGroup(id)
}

func GroupSetLevel(id int64, level int) (group TradfriGroup, err error) {
	uri := fmt.Sprintf("%s/%d", uriGroups, id)
	payload := fmt.Sprintf("{\"%s\":%d}", attrLightDimmer, level)
	_, err = PutRequest(uri, payload)
	if err != nil {
		return group, err
	}
	return GetGroup(id)
}
