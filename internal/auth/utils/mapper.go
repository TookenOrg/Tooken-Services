package utils

import "encoding/json"

func MapOneToOne(input interface{}, output interface{}) {
	marshal, err := json.Marshal(input)
	if err != nil {
		return
	}

	json.Unmarshal(marshal, &output)
}
