package utils

import "encoding/json"

// ToJSONString converts any payload to a JSON string, ignoring errors.
func ToJSONString(payload interface{}) string {
	if payload == nil {
		return ""
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return ""
	}
	return string(data)
}
