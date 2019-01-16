package utils

import (
	"encoding/json"
	"fmt"
)

func ToJSON(i interface{}) string {
	b, err := json.Marshal(i)
	if err != nil {
		return fmt.Sprintf(`{"msg": "utils: failed to marshal JSON: %v"}`, err)
	}
	return string(b)
}
