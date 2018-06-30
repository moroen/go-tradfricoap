package tradfricoap

import (
	"errors"

	// "os"
	"time"

	"github.com/moroen/canopus"
	// "github.com/urfave/cli"
	// "github.com/buger/jsonparser"
)

type CoapResult struct {
	msg canopus.MessagePayload
	err error
}

var ErrorTimeout = errors.New("COAP Error: Connection timeout")
var ErrorBadIdent = errors.New("COAP DTLS Error: Wrong credentials?")
var ErrorNoConfig = errors.New("COAP Error: No config")

func _getRequest(URI string, c chan CoapResult) {

	var result CoapResult

	conf, err := GetConfig()
	if err != nil {
		result.err = ErrorNoConfig
		c <- result
		return
	}

	conn, err := canopus.DialDTLS(conf.Gateway, conf.Identity, conf.Passkey)
	if err != nil {
		result.err = err
		c <- result
		return
	}

	req := canopus.NewRequest(canopus.MessageConfirmable, canopus.Get)
	req.SetStringPayload("Hello, canopus")
	req.SetRequestURI(URI)

	resp, err := conn.Send(req)
	if err != nil {
		result.err = ErrorBadIdent
		c <- result
		return
	}

	// response := resp.GetMessage().GetPayload()
	result.err = nil
	result.msg = resp.GetMessage().GetPayload()
	c <- result
}

func _putRequest(URI, payload string, c chan CoapResult) {
	var result CoapResult

	conf, err := GetConfig()
	if err != nil {
		result.err = ErrorNoConfig
		c <- result
		return
	}

	conn, err := canopus.DialDTLS(conf.Gateway, conf.Identity, conf.Passkey)
	if err != nil {
		result.err = err
		c <- result
		return
	}

	req := canopus.NewRequest(canopus.MessageConfirmable, canopus.Put)
	req.SetRequestURI(URI)
	req.SetStringPayload(payload)

	resp, err := conn.Send(req)
	if err != nil {
		result.err = ErrorBadIdent
		c <- result
		return
	}

	result.msg = resp.GetMessage().GetPayload()
	result.err = nil
	c <- result
}

func GetRequest(URI string) (msg canopus.MessagePayload, err error) {
	c := make(chan CoapResult)

	go _getRequest(URI, c)

	select {
	case res := <-c:
		return res.msg, res.err
	case <-time.After(time.Second * 5):
		return nil, ErrorTimeout
	}
}

func PutRequest(URI, payload string) (msg canopus.MessagePayload, err error) {
	c := make(chan CoapResult)

	go _putRequest(URI, payload, c)

	select {
	case res := <-c:
		return res.msg, res.err
	case <-time.After(time.Second * 5):
		return nil, ErrorTimeout
	}
}
