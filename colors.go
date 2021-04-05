package tradfricoap

import (
	"fmt"
	"log"
	"sort"

	colorful "github.com/lucasb-eyer/go-colorful"
)

type ColorMap map[int]map[string]string

func cwMap() ColorMap {
	var whiteBalance = ColorMap{
		0:  {"Name": "Off", "Hex": "000000"},
		10: {"Name": "Cold", "Hex": "f5faf6"},
		20: {"Name": "Normal", "Hex": "f1e0b5"},
		30: {"Name": "Warm", "Hex": "efd275"},
	}

	return whiteBalance
}

func cwsMap() ColorMap {
	var cws = ColorMap{
		0:   {"Name": "Off", "Hex": "000000"},
		10:  {"Name": "Blue", "Hex": "4a418a"},
		20:  {"Name": "Candlelight", "Hex": "ebb63e"},
		30:  {"Name": "Cold sky", "Hex": "dcf0f8"},
		40:  {"Name": "Cool daylight", "Hex": "eaf6fb"},
		50:  {"Name": "Cool white", "Hex": "f5faf6"},
		60:  {"Name": "Dark Peach", "Hex": "da5d41"},
		70:  {"Name": "Light Blue", "Hex": "6c83ba"},
		80:  {"Name": "Light Pink", "Hex": "e8bedd"},
		90:  {"Name": "Light Purple", "Hex": "c984bb"},
		100: {"Name": "Lime", "Hex": "a9d62b"},
		110: {"Name": "Peach", "Hex": "e57345"},
		120: {"Name": "Pink", "Hex": "e491af"},
		130: {"Name": "Saturated Red", "Hex": "dc4b31"},
		140: {"Name": "Saturated Pink", "Hex": "d9337c"},
		150: {"Name": "Saturated Purple", "Hex": "8f2686"},
		160: {"Name": "Sunrise", "Hex": "f2eccf"},
		170: {"Name": "Yellow", "Hex": "d6e44b"},
		180: {"Name": "Warm Amber", "Hex": "e78834"},
		190: {"Name": "Warm glow", "Hex": "efd275"},
		200: {"Name": "Warm white", "Hex": "f1e0b5"},
	}
	return cws
}

func hexForLevel(colorMap ColorMap, level int) string {
	return colorMap[level]["Hex"]
}

func ListColorsInMap(colorMap ColorMap) {
	var keys []int

	for k := range colorMap {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, val := range keys {
		fmt.Println(fmt.Sprintf("%d: %s", val, colorMap[val]["Name"]))
	}
}

func SetRGB(id int64, rgb string) {
	fmt.Printf("Device: %d - RGB: %s\n", id, rgb)
	c, err := colorful.Hex("#517AB8")
	if err != nil {
		log.Fatal(err)
	}

	h, s, v := c.Hsv()
	x, y, z := c.Xyz()

	fmt.Println("HSV: ", h, s, v)

	fmt.Println("xyZ:", x, y, z)

	//uri := fmt.Sprintf("%s/%d", uri_Devices, id)

	//payload := fmt.Sprintf("{ \"%s\": [{ \"%s\": %d }] }", attr_Light_control, attr_light_state, state)
	//fmt.Println(payload)
	//_, err = PutRequest(uri, payload)

	// if err != nil {
	//	return device, err
	//}
}
