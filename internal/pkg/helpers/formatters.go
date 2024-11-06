package helpers

import "time"

func DatetimeToString(field time.Time) string {
	return field.Format("2006-01-02 15:04:05")
}

func DBByteToString(v interface{}) string {
	if b, ok := v.([]uint8); ok {
		return string(b)
	}
	return ""
}
