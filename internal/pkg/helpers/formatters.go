package helpers

import (
	"fmt"
	"log"
	"math"
	"strings"
	"time"
)

func DatetimeToString(field time.Time) string {
	return field.Format("2006-01-02 15:04:05")
}

func DBByteToString(v interface{}) string {
	if b, ok := v.([]uint8); ok {
		return string(b)
	}
	return ""
}

func NumberFormat(number interface{}, decimals int, decPoint, thousandsSep string) string {
	log.Println(number)
	var floatValue float64

	switch v := number.(type) {
	case float64:
		floatValue = v
	case float32:
		floatValue = float64(v)
	case int:
		floatValue = float64(v)
	case int32:
		floatValue = float64(v)
	case int64:
		floatValue = float64(v)
	default:
		return "Unsupported type"
	}

	// Handle special cases
	if math.IsNaN(floatValue) {
		return "NaN"
	}
	if math.IsInf(floatValue, 0) {
		if floatValue > 0 {
			return "Inf"
		}
		return "-Inf"
	}

	// Round the number to the specified number of decimal places
	str := fmt.Sprintf("%."+fmt.Sprintf("%d", decimals)+"f", math.Abs(floatValue))

	// Split the string into integer and decimal parts
	parts := strings.Split(str, ".")

	// Add thousand separators to the integer part
	integerPart := parts[0]
	var result []string
	for i := len(integerPart); i > 0; i -= 3 {
		if i-3 > 0 {
			result = append([]string{thousandsSep + integerPart[i-3:i]}, result...)
		} else {
			result = append([]string{integerPart[0:i]}, result...)
		}
	}

	// Join the parts back together
	formattedNumber := strings.Join(result, "")
	if decimals > 0 {
		formattedNumber += decPoint + parts[1]
	}

	// Add the negative sign if necessary
	if floatValue < 0 {
		formattedNumber = "-" + formattedNumber
	}

	return formattedNumber
}
