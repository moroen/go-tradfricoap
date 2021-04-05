package tradfricoap

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/buger/jsonparser"
	"github.com/shibukawa/configdir"

	coap "github.com/moroen/gocoap/v4"
)

var globalGatewayConfig GatewayConfig

type GatewayConfig struct {
	Gateway  string
	Identity string
	Passkey  string
}

var configDirs = configdir.New("", "tradfri")

// ErrorNoConfig error
var ErrorNoConfig = errors.New("Tradfri Error: No config")

func (c GatewayConfig) Describe() string {
	out, _ := json.Marshal(c)
	return string(out)
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

	folder := configDirs.QueryFolderContainsFile("gateway.json")

	if folder == nil {
		return config, errors.New("Config not found")
	}

	data, err := folder.ReadFile("gateway.json")
	if err != nil {
		return config, errors.New("Config not found")
	}

	if err := json.Unmarshal(data, &config); err == nil {
		SetConfig(config)
	}

	return config, nil
}

func SaveConfig(conf GatewayConfig) (err error) {
	data, _ := json.Marshal(&conf)
	folders := configDirs.QueryFolders(configdir.Global)

	err = folders[0].WriteFile("gateway.json", data)
	if err == nil {
		log.Println("Saved new config: ", conf.Describe())
	} else {
		log.Println(err.Error())
	}
	return err
}

func CreateIdent(gateway, key, ident string) error {

	payload := fmt.Sprintf("{\"%s\":\"%s\"}", attrIdent, ident)
	URI := uriIdent

	param := coap.RequestParams{Host: gateway, Port: 5684, Uri: URI, Id: "Client_identity", Key: key, Payload: payload}

	res, err := coap.PostRequest(param)
	if err != nil {
		return err
	}

	psk, err := jsonparser.GetString(res, "9091")
	if err != nil {
		return err
	}

	conf := GatewayConfig{Gateway: fmt.Sprintf("%s", gateway), Identity: ident, Passkey: psk}
	SaveConfig(conf)
	SetConfig(conf)
	return nil
}

func GetRequest(URI string) (retmsg []byte, err error) {

	conf, err := GetConfig()
	if err != nil {
		return nil, err
	}

	param := coap.RequestParams{Host: conf.Gateway, Port: 5684, Uri: URI, Id: conf.Identity, Key: conf.Passkey}

	res, err := coap.GetRequest(param)
	if err != nil {
		if err == coap.ErrorHandshake {
			log.Fatalln("Connection timed out")
		} else {
			log.Println(err.Error())
			// panic(err.Error())
		}
	}

	return res, err
}

func PutRequest(URI string, Payload string) (retmsg []byte, err error) {

	conf, err := GetConfig()
	if err != nil {
		panic(err.Error())
	}

	param := coap.RequestParams{Host: conf.Gateway, Port: 5684, Uri: URI, Id: conf.Identity, Key: conf.Passkey, Payload: Payload}

	res, err := coap.PutRequest(param)
	if err != nil {
		panic(err.Error())
	}

	return res, nil
}
