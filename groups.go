package tradfricoap

import (
	"fmt"

	"github.com/buger/jsonparser"
)

type TradfriGroup struct {
	Id   int64
	Name string
}

func (g TradfriGroup) Describe() string {
	return fmt.Sprint(g)
}

type TradfriGroups []TradfriGroup

func GetGroup(id string) (TradfriGroup, error) {
	var aGroup TradfriGroup

	return aGroup, nil
}

func GetGroups() (TradfriGroups, error) {
	payload, err := GetRequest(uri_Groups)
	if err != nil {
		return nil, err
	}

	msg := payload.GetBytes()

	//msg = strings.Trim(msg, "[")
	//msg = strings.Trim(msg, "]")
	//result := strings.Split(msg, ",")

	groups := []TradfriGroup{}

	jsonparser.ArrayEach(msg, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if res, err := jsonparser.GetInt(value); err == nil {
			fmt.Println(res)
		}
	})

	// for i := range result {
	// 	aGroup, err := GetGroup(result[i])
	// 	if err == nil {
	// 		groups = append(groups, aGroup)
	// 	}
	// }

	// sort.Slice(groups, func(i, j int) bool {
	// 	return groups[i].Id < groups[j].Id
	// })

	return groups, err
}
