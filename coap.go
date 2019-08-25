package tradfricoap

import (
	"errors"
	"fmt"

	// "os"
	"time"

	// "github.com/moroen/canopus"
	// "github.com/moroen/canopus"
	// "github.com/urfave/cli"
	// "github.com/buger/jsonparser"

	// "github.com/dustin/go-coap"
	"github.com/dustin/go-coap"
	"github.com/eriklupander/dtls"
)

/*
type CoapResult struct {
	msg canopus.MessagePayload
	err error
}
*/
var ErrorTimeout = errors.New("COAP Error: Connection timeout")
var ErrorBadIdent = errors.New("COAP DTLS Error: Wrong credentials?")
var ErrorNoConfig = errors.New("COAP Error: No config")

/*
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
*/

func _request(req coap.Message) (retmsg coap.Message, err error) {
	conf, err := GetConfig()
	if err != nil {
		// result.err = ErrorNoConfig
		// c <- result
		panic("No config")
		// return
	}

	mks := dtls.NewKeystoreInMemory()
	dtls.SetKeyStores([]dtls.Keystore{mks})
	mks.AddKey(conf.Identity, []byte(conf.Passkey))

	listner, err := dtls.NewUdpListener(":0", time.Second*900)
	if err != nil {
		panic(err.Error())
	}

	peerParams := &dtls.PeerParams{
		Addr:             fmt.Sprintf("%s:%d", conf.Gateway, 5684),
		Identity:         conf.Identity,
		HandshakeTimeout: time.Second * 15}

	peer, err := listner.AddPeerWithParams(peerParams)
	if err != nil {
		panic(err.Error())
	}

	peer.UseQueue(true)

	data, err := req.MarshalBinary()
	if err != nil {
		panic(err.Error())
	}

	err = peer.Write(data)
	if err != nil {
		panic(err.Error())
	}

	respData, err := peer.Read(time.Second)
	if err != nil {
		panic(err.Error())
	}

	msg, err := coap.ParseMessage(respData)
	if err != nil {
		panic(err.Error())
	}

	err = listner.Shutdown()
	if err != nil {
		panic(err.Error())
	}

	return msg, nil
}

func _getRequest(URI string, c chan coap.Message) {
	println("Getting request")

	req := coap.Message{
		Type:      coap.Confirmable,
		Code:      coap.GET,
		MessageID: 1,
	}

	req.SetPathString(URI)

	msg, err := _request(req)
	if err != nil {
		panic(err.Error())
	}
	c <- msg
}

func _putRequest(URI, payload string, c chan coap.Message) {
	req := coap.Message{
		Type:      coap.Confirmable,
		Code:      coap.PUT,
		MessageID: 1,
		Payload:   []byte(payload),
	}
	req.SetPathString(URI)

	msg, err := _request(req)
	if err != nil {
		panic(err.Error())
	}
	c <- msg
}

// GetRequest sends a default get
func GetRequest(URI string) (msg coap.Message, err error) {
	c := make(chan coap.Message)

	go _getRequest(URI, c)

	select {
	case res := <-c:
		return res, nil
	case <-time.After(time.Second * 60):
		return coap.Message{}, ErrorTimeout
	}
}

// PutRequest sends a default Put-request
func PutRequest(URI, payload string) (msg coap.Message, err error) {
	c := make(chan coap.Message)

	go _putRequest(URI, payload, c)

	select {
	case res := <-c:
		return res, nil
	case <-time.After(time.Second * 5):
		return coap.Message{}, ErrorTimeout
	}
}
