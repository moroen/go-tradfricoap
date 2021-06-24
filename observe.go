package tradfricoap

import (
	"fmt"

	coap "github.com/moroen/gocoap/v4"
	log "github.com/sirupsen/logrus"
)

func ObserveStop() {
	coap.ObserveStop()
}

func ObserveRestart(reconnect bool) {
	coap.ObserveRestart(reconnect)
}

func Observe(callback func([]byte) error, control_channel chan (error)) error {
	var endPoints []string

	conf, err := GetConfig()
	if err != nil {
		log.Println(err.Error())
	}

	log.WithFields(log.Fields{
		"Id":  conf.Identity,
		"Key": conf.Passkey,
	}).Debug("COAP observe")

	param := coap.ObserveParams{Host: conf.Gateway, Port: 5684, Id: conf.Identity, Key: conf.Passkey}

	lights, plugs, blinds, _, err := GetDevices()

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Observe - get endpoints failed")
		control_channel <- err
		return err
	}

	for _, light := range lights {
		endPoints = append(endPoints, fmt.Sprintf("15001/%d", light.Id))
	}

	for _, plug := range plugs {
		endPoints = append(endPoints, fmt.Sprintf("15001/%d", plug.Id))
	}

	for _, blind := range blinds {
		endPoints = append(endPoints, fmt.Sprintf("15001/%d", blind.Id))
	}

	param.Uri = endPoints
	go coap.Observe(param, callback)
	return nil
}
