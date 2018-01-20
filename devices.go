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
