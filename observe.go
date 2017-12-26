package tradfricoap

import (
	"fmt"

	"github.com/moroen/canopus"
)

func Observe(uri string) {
	conn, err := canopus.DialDTLS(globalGatewayConfig.Gateway, globalGatewayConfig.Identity, globalGatewayConfig.Passkey)
	tok, err := conn.ObserveResource(uri)
	if err != nil {
		panic(err.Error())
	}

	obsChannel := make(chan canopus.ObserveMessage)
	done := make(chan bool)
	go conn.Observe(obsChannel)

	notifyCount := 0
	go func() {
		for {
			select {
			case obsMsg, open := <-obsChannel:
				if open {
					notifyCount++
					// msg := obsMsg.Msg\
					resource := obsMsg.GetResource()
					val := obsMsg.GetValue()

					fmt.Println("[CLIENT >> ] Got Change Notification for resource and value: ", notifyCount, resource, val)
					go conn.CancelObserveResource(uri, tok)
					go conn.StopObserve(obsChannel)

					obsChannel = make(chan canopus.ObserveMessage)
					go conn.Observe(obsChannel)
				} else {
					// done <- true
					// return
				}
			}
		}
	}()
	<-done
	fmt.Println("Done")
}
