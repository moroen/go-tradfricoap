package tradfricoap

import (
	"fmt"
)

type TradfriGroup struct {
	Id   int64
	Name string
}

func (g TradfriGroup) Describe() string {
	return fmt.Sprint(g)
}

type TradfriGroups []TradfriGroup
