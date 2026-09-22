package utils

func ValueOrEmpty(record map[string]any, key string) any {
	if value := record[key]; value != nil {
		return value
	}
	return ""
}
