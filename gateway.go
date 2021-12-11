package tradfricoap

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/buger/jsonparser"
	"github.com/shibukawa/configdir"

	coap "github.com/moroen/gocoap/v5"
	uuid "github.com/satori/go.uuid"
)

var globalGatewayConfig GatewayConfig

type GatewayConfig struct {
	Gateway   string
	Identity  string
	Passkey   string
	KeepAlive int
}

var coapDTLSConnection coap.CoapDTLSConnection
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
	return _loadConfig("gateway.json")
}

func LoadConfigFile(fName string) (config GatewayConfig, err error) {
	return _loadConfig(fName)
}

func _loadConfig(fName string) (config GatewayConfig, err error) {

	folder := configDirs.QueryFolderContainsFile(fName)

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

type DTLSPSKpair struct {
	Ident string
	Key   string
}

func GetNewPSK(gateway, key string) (DTLSPSKpair, error) {

	ident := uuid.NewV4().String()

	payload := fmt.Sprintf("{\"%s\":\"%s\"}", attrIdent, ident)
	URI := uriIdent

	param := coap.RequestParams{Host: gateway, Port: 5684, Uri: URI, Id: "Client_identity", Key: key, Payload: payload}

	coap.CloseDTLSConnection()
	res, err := coap.PostRequest(param)
	if err != nil {
		return DTLSPSKpair{}, err
	}

	psk, err := jsonparser.GetString(res, "9091")
	if err != nil {
		return DTLSPSKpair{}, err
	}

	return DTLSPSKpair{Ident: ident, Key: psk}, nil
}

func GetRequestWithContext(ctx context.Context, URI string) (retmsg []byte, err error) {
	conf, err := GetConfig()
	if err != nil {
		return nil, err
	}

	param := coap.RequestParams{Host: conf.Gateway, Port: 5684, Uri: URI, Id: conf.Identity, Key: conf.Passkey}

	res, err := coap.GetRequestWithContext(ctx, param, 10)
	if err != nil {
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error("GetRequest failed")
		}
	}
	return res, err
}

func GetRequest(URI string) (retmsg []byte, err error) {

	conf, err := GetConfig()
	if err != nil {
		return nil, err
	}

	param := coap.RequestParams{Host: conf.Gateway, Port: 5684, Uri: URI, Id: conf.Identity, Key: conf.Passkey}

	res, err := coap.GetRequest(param)
	if err != nil {
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error("GetRequest failed")
		}
	}

	return res, err
}

func PutRequest(URI string, Payload string) (retmsg []byte, err error) {
	/*
		var res []byte

		conf, err := GetConfig()
		if err != nil {
			panic(err.Error())
		}

		param := coap.RequestParams{Host: conf.Gateway, Port: 5684, Uri: URI, Id: conf.Identity, Key: conf.Passkey, Payload: Payload}

		res, err = coap.PutRequest(param)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error("PutRequest failed")
		}
		return res, err
	*/
	log.Error("PutRequest - disabled")
	return nil, nil
}

func SetCoapRetry(limit uint, delay int) {
	coap.SetRetry(limit, delay)
}

func ConnectGateway(ctx context.Context, cfg GatewayConfig, onConnect func(), onDisconnect func(), onCanceled func(), onConnectionFailed func()) error {
	// SetConfig(cfg)

	coapDTLSConnection = coap.CoapDTLSConnection{}

	coapDTLSConnection.Host = cfg.Gateway
	coapDTLSConnection.Port = 5684
	coapDTLSConnection.Ident = cfg.Identity
	coapDTLSConnection.Key = cfg.Passkey

	coapDTLSConnection.OnConnect = onConnect
	coapDTLSConnection.OnDisconnect = onDisconnect
	coapDTLSConnection.OnCanceled = onCanceled
	coapDTLSConnection.OnConnectionFailed = onConnectionFailed

	coapDTLSConnection.Connect(ctx)
	return nil
}

func CloseConnection() error {
	return coap.CloseDTLSConnection()
}
