package tradfricoap

import (

	// "os"

	"github.com/moroen/canopus"
	// "github.com/urfave/cli"
	// "github.com/buger/jsonparser"
)

func GetRequest(URI string) (msg canopus.MessagePayload, err error) {

	conn, err := canopus.DialDTLS(globalGatewayConfig.Gateway, globalGatewayConfig.Identity, globalGatewayConfig.Passkey)
	if err != nil {
		panic(err.Error())
	}

	req := canopus.NewRequest(canopus.MessageConfirmable, canopus.Get)
	req.SetStringPayload("Hello, canopus")
	req.SetRequestURI(URI)

	resp, err := conn.Send(req)
	if err != nil {
		panic(err.Error())
	}

	response := resp.GetMessage().GetPayload()
	return response, err
}

func PutRequest(URI, payload string) (msg canopus.MessagePayload, err error) {
	conn, err := canopus.DialDTLS(globalGatewayConfig.Gateway, globalGatewayConfig.Identity, globalGatewayConfig.Passkey)
	if err != nil {
		panic(err.Error())
	}

	// fmt.Println(URI, payload)

	req := canopus.NewRequest(canopus.MessageConfirmable, canopus.Put)
	req.SetRequestURI(URI)
	req.SetStringPayload(payload)

	resp, err := conn.Send(req)
	if err != nil {
		panic(err.Error())
	}
	response := resp.GetMessage().GetPayload()
	// println("Response: ", response.String())
	return response, err
}
