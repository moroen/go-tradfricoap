package tradfricoap

import ( // "log"
	// "os"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/buger/jsonparser"
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

func GetDevice(id int64) ([]byte, error) {
	msg, err := GetRequest(fmt.Sprintf("%s/%d", uriDevices, id))
	if err != nil {
		return nil, err
	}

	return msg, err
}

func GetDevices() (lights TradfriLights, plugs TradfriPlugs, blinds TradfriBlinds, groups TradfriGroups, err error) {
	result, err := GetRequest(uriDevices)
	if err != nil {
		// fmt.Println(err.Error())
		return nil, nil, nil, nil, err
	}

	msg := result

	_, err = jsonparser.ArrayEach(msg, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if res, err := jsonparser.GetInt(value); err == nil {
			aDevice, err := GetDevice(res)
			if err != nil {
				log.Println(err.Error())
				return
			}

			if _, _, _, err := jsonparser.Get(aDevice, attrLightControl); err == nil {
				if aLight, err := getLightInfo(aDevice); err == nil {
					lights = append(lights, aLight)
				}

			}

			if _, _, _, err := jsonparser.Get(aDevice, attrPlugControl); err == nil {
				if aPlug, err := getPlugInfo(aDevice); err == nil {
					plugs = append(plugs, aPlug)
				}

			}

			if _, _, _, err := jsonparser.Get(aDevice, attrBlindControl); err == nil {
				if aBlind, err := getBlindInfo(aDevice); err == nil {
					blinds = append(blinds, aBlind)
				}

			}

		}
		// time.Sleep(1000 * time.Millisecond)
	})

	groups, err = GetGroups()
	if err != nil {
		panic(err.Error())
	}

	sort.Slice(lights, func(i, j int) bool {
		return lights[i].Id < lights[j].Id
	})

	sort.Slice(plugs, func(i, j int) bool {
		return plugs[i].Id < plugs[j].Id
	})

	sort.Slice(blinds, func(i, j int) bool {
		return blinds[i].Id < blinds[j].Id
	})

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Id < groups[j].Id
	})

	return lights, plugs, blinds, groups, err
}
