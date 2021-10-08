package tradfricoap

import (
	"fmt"
	"sync"

	"github.com/buger/jsonparser"
	coap "github.com/moroen/gocoap/v4"
	log "github.com/sirupsen/logrus"
)

var _control_channel chan (error)

var _wg *sync.WaitGroup

func ObserveStop() {
	defer _wg.Done()
	log.Info("Stopping tradfri...")
	coap.ObserveStop()
	log.Info("Tradfri stopped")

}

func ObserveRestart(reconnect bool) {
	coap.ObserveRestart(reconnect)
}

func Observe(wg *sync.WaitGroup, callback func([]byte) error) (chan (error), error) {
	var endPoints []string
	_wg = wg
	_wg.Add(1)

	_control_channel = make(chan error)

	conf, err := GetConfig()
	if err != nil {
		fmt.Println("Shite")
		return _control_channel, err
	}

	log.WithFields(log.Fields{
		"Id":  conf.Identity,
		"Key": conf.Passkey,
	}).Debug("COAP observe")

	if result, err := GetRequest(uriDevices); err == nil {
		msg := result

		if _, err = jsonparser.ArrayEach(msg, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			if res, err := jsonparser.GetInt(value); err == nil {
				endPoints = append(endPoints, fmt.Sprintf("15001/%d", res))
			}
		}); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	param := coap.ObserveParams{Host: conf.Gateway, Port: 5684, Id: conf.Identity, Key: conf.Passkey}

	param.Uri = endPoints
	go coap.Observe(param, callback)
	return _control_channel, nil
}
