package tradfricoap

import (
	"fmt"

	"github.com/moroen/canopus"
)

type TradfriLight struct {
	Id     int64
	Name   string
	State  string
	Dimmer int64
}

func (l TradfriLight) Describe() string {
	return fmt.Sprintf("%d: %s - %s (%d)", l.Id, l.Name, l.State, l.Dimmer)
}

func SetState(id string, state int) (msg canopus.MessagePayload, err error) {
	uri := fmt.Sprintf("/15001/%s", id)
	payload := fmt.Sprintf("{ \"%s\": [{ \"%s\": %d }] }", attr_Light_control, attr_light_state, state)
	return PutRequest(uri, payload)
}

func SetLevel(id string, level int) (msg canopus.MessagePayload, err error) {
	uri := fmt.Sprintf("/15001/%s", id)
	payload := fmt.Sprintf("{ \"%s\": [{ \"%s\": %d }] }", attr_Light_control, attr_light_dimmer, level)
	return PutRequest(uri, payload)
}
