package util

import "encoding/json"

// TOJSON convert an interface as JSON string
func TOJSON(data interface{}) (result string) {
	str, err := json.Marshal(data)
	if err == nil {
		result = string(json.RawMessage(str))
	}
	return
}
