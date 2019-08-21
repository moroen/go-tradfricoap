package tradfricoap

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/shibukawa/configdir"
	"github.com/tucnak/store"
)

var globalGatewayConfig GatewayConfig

type GatewayConfig struct {
	Gateway  string
	Identity string
	Passkey  string
}

func (c GatewayConfig) Describe() string {
	out, _ := json.Marshal(c)
	return string(out)
}

func init() {
	// You must init store with some truly unique path first!
	store.Init("tradfri")
}

func SetConfig(c GatewayConfig) {
	globalGatewayConfig = c
}

func GetConfig() (conf GatewayConfig, err error) {
	if globalGatewayConfig == (GatewayConfig{}) {
		err = ErrorNoConfig
	}
	return globalGatewayConfig, err
}

func LoadConfig() (config GatewayConfig, err error) {

	configDir := configdir.New("", "tradfri")
	folder := configDir.QueryFolderContainsFile("gateway.json")

	if folder == nil {
		return config, errors.New("Config not found")
	}

	data, err := folder.ReadFile("gateway.json")
	if err != nil {
		return config, errors.New("Config not found")
	}

	json.Unmarshal(data, &config)
	SetConfig(config)
	return config, nil
}

func SaveConfig(conf GatewayConfig) (err error) {
	err = store.Save("gateway.json", &conf)
	if err == nil {
		log.Println("Saved new config: ", conf.Describe())
	}
	return err
}

func CreateIdent(gateway, key, ident string) {
	/*
		payload := fmt.Sprintf("{\"%s\":\"%s\"}", attr_Ident, ident)
		URI := uri_Ident

		conn, err := canopus.DialDTLS(fmt.Sprintf("%s:%s", gateway, DefaultPort), "Client_identity", key)
		if err != nil {
			panic(err.Error())
		}

		req := canopus.NewRequest(canopus.MessageConfirmable, canopus.Post)
		req.SetRequestURI(URI)
		req.SetStringPayload(payload)

		response, err := conn.Send(req)
		if err != nil {
			panic(err.Error())
		}

		if response.GetMessage().GetCode() == 65 {
			result := response.GetMessage().GetPayload().GetBytes()
			psk, err := jsonparser.GetString(result, "9091")
			if err != nil {
				panic(err.Error())
			}

			conf := GatewayConfig{Gateway: fmt.Sprintf("%s:%s", gateway, DefaultPort), Identity: ident, Passkey: psk}
			SaveConfig(conf)
			SetConfig(conf)

		} else {
			log.Printf("Unable to get PSK for ident '%s'. Ident already in use?", ident)
		}
		// fmt.Println("Code: ", resp.GetMessage().GetCode())
		// response := resp.GetMessage().GetPayload()
		// fmt.Println("Response: ", response.String())
	*/
}
