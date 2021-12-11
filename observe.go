package tradfricoap

import (
	"context"
	"fmt"

	"github.com/buger/jsonparser"
	coap "github.com/moroen/gocoap/v5"
	log "github.com/sirupsen/logrus"
)

var _control_channel chan (error)

var ctxObserve context.Context

func ObserveStop() {
	coap.ObserveStop()
}

func ObserveRestart(reconnect bool) {
	coap.ObserveRestart(reconnect)
}

func Observe(callback func([]byte) error, keepAlive int) (chan (error), error) {
	var endPoints []string

	_control_channel = make(chan error)

	conf, err := GetConfig()
	if err != nil {
		return _control_channel, err
	}

	log.WithFields(log.Fields{
		"Id":  conf.Identity,
		"Key": conf.Passkey,
	}).Debug("COAP observe")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if result, err := GetRequestWithContext(ctx, uriDevices); err == nil {
		msg := result

		if _, err = jsonparser.ArrayEach(msg, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			if res, err := jsonparser.GetInt(value); err == nil {
				endPoints = append(endPoints, fmt.Sprintf("15001/%d", res))
			}
		}); err != nil {
			return nil, err
		}
	} else {
		fmt.Println("Fails here")
		return nil, err
	}

	param := coap.ObserveParams{
		Host:      conf.Gateway,
		Port:      5684,
		Uri:       []string{},
		Id:        conf.Identity,
		Key:       conf.Passkey,
		KeepAlive: keepAlive,
	}

	param.Uri = endPoints
	go coap.Observe(param, callback)
	return _control_channel, nil
}
