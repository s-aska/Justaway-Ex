package handlers

import (
	"encoding/json"
)

func encodeJson(d interface{}) (j string) {
	b, _ := json.Marshal(d)
	return string(b)
}
