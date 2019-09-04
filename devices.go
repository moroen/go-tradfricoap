package tradfricoap

import ( // "log"
	// "os"
	"fmt"
	"strconv"
	"strings"
	// "github.com/moroen/canopus"
	// "github.com/urfave/cli"
)

func trimJSON(json string) string {
	json = strings.Trim(json, "[")
	json = strings.Trim(json, "]")
	return json
}

func ValidateDeviceID(id string) error {
	if _, err := strconv.Atoi(id); err != nil {
		return fmt.Errorf("%s doesn't appear to be a valid tradfri device", id)
	}
	return nil
}

func ValidateOnOff(arg string) error {
	if strings.ToLower(arg) == "on" || strings.ToLower(arg) == "off" || strings.ToLower(arg) == "1" || strings.ToLower(arg) == "0" {
		return nil
	} else {
		return fmt.Errorf("%s isn't an allowed setting, use 'on', 'off', '1' or '0'", arg)
	}
}

/*
func GetDevices() (TradfriLights, TradfriPlugs, TradfriGroups, error) {

}
*/
