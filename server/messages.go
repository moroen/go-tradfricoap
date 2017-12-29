package main

import (
	coap "github.com/moroen/go-tradfricoap"
)

type returnMessageSimple struct {
	Action string
	Status string
	Result string
}

type returnMessageDevices struct {
	Action string             `json:"action"`
	Status string             `json:"status"`
	Result coap.TradfriLights `json:"result"`
}
