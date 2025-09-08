package services

import (
	"fmt"
	"strings"
)

// ____________________________________________________INTERNAL____________________________________________________
// убирает пустые значения из updates (для частичных апдейтов)
func FilterUpdates(updates map[string]interface{}) map[string]interface{} {
	filtered := make(map[string]interface{})

	for key, value := range updates {
		if value == nil {
			continue
		}

		switch v := value.(type) {
		case string:
			if strings.TrimSpace(v) == "" {
				continue
			}
			filtered[key] = v
		case float64:
			if v == 0 {
				continue
			}
			filtered[key] = v
		case int, int32, int64:
			if fmt.Sprintf("%v", v) == "0" {
				continue
			}
			filtered[key] = v
		case bool:
			if !v {
				continue
			}
			filtered[key] = v
		default:
			filtered[key] = v
		}
	}

	return filtered
}
