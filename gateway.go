package tradfricoap

import (
	"encoding/json"
	"errors"

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

func GetConfig() (conf GatewayConfig) {
	return globalGatewayConfig
}

func LoadConfig() (err error) {
	var conf GatewayConfig
	store.Load("gateway.json", &conf)
	if conf == (GatewayConfig{}) {
		err = errors.New("Configuration not found")
	} else {
		SetConfig(conf)
	}
	return err
}

func SaveConfig(conf GatewayConfig) (err error) {
	err = store.Save("gateway.json", &conf)
	return err
}
