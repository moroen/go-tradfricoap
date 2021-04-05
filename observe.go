package tradfricoap

import (
	"log"
)

func Observe() {

	log.Println("Observe called")

	/*
		conf, err := GetConfig()
		if err != nil {
			panic(err.Error())
		}

		param := coap.ObserveParams{Host: conf.Gateway, Port: 5684, ID: conf.Identity, Key: conf.Passkey}

		endpoints := `["15001/65554", "15001/65550"]`

		var uris []string

		err = json.Unmarshal([]byte(endpoints), &uris)
		if err != nil {
			panic(err.Error())
		}

		param.URI = uris

		msg := make(chan []byte)
		sign := make(chan bool)
		obserr := make(chan error)

		state := 0

		go coap.Observe(param, msg, sign, obserr)
		for {
			select {
			case message := <-msg:
				fmt.Println(string(message))
			case <-time.After(10 * time.Second):
				SetState(65554, state)
				state = 1 - state
			case <-time.After(120 * time.Second):
				sign <- true
				return
			}
		}
	*/
}
