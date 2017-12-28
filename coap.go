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

/*
func SetLighControll(device, attribute string, value int) {

	conn, err := canopus.DialDTLS(globalGatewayConfig.Gateway, globalGatewayConfig.Identity, globalGatewayConfig.Passkey)
	if err != nil {
		panic(err.Error())
	}

	payload := fmt.Sprintf("{ \"%s\": [{ \"%s\": %d }] }", attr_Light_control, attr_light_state, stateValue)

	fmt.Println(payload)

	req := canopus.NewRequest(canopus.MessageConfirmable, canopus.Put)
	req.SetRequestURI(fmt.Sprintf("/15001/%s", endpoint))
	req.SetStringPayload(payload)

	resp, err := conn.Send(req)

	if err != nil {
		log.Println(err)
	} else {
		log.Println("Got Response:")
		log.Println(resp.GetMessage().GetPayload().String())
	}
}
*/
