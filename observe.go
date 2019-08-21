package tradfricoap

func Observe(uri string, status chan bool) {
	/*
		// fmt.Printf("Calling Observe for %s", uri)

		// conn, err := canopus.DialDTLS(globalGatewayConfig.Gateway, globalGatewayConfig.Identity, globalGatewayConfig.Passkey)
		conn, err := canopus.DialDTLS(globalGatewayConfig.Gateway, "IKEA03", "v8CCsb5bbw2qmfwU")
		fmt.Println("Using IDENT: ", "IKEA03")
		if err != nil {
			panic(err.Error())
		}

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
						if notifyCount == 2 {
							fmt.Println("[CLIENT >> ] Canceling observe after 2 notifications..")
							go conn.CancelObserveResource(uri, tok)
							// go conn.StopObserve(obsChannel)
							done <- true
							return
						}
					} else {
						done <- true
						return
					}
				}
			}
		}()
		<-done
		fmt.Println("Done")
		status <- true
		// stat := conn.Close()
		// if stat != nil {
		// 	panic(err.Error())
		// }
	*/
}
