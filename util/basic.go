package util

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

// ToStr converts any value to string.
func ToStr(v interface{}) string {
	if v == nil {
		return ""
	}
	switch v := v.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case int32:
		return strconv.Itoa(int(v))
	case int64:
		return strconv.FormatInt(v, 10)
	case bool:
		return strconv.FormatBool(v)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case []uint8:
		return string(v)
	default:
		resultJSON, err := json.Marshal(v)
		if err != nil {
			return ""
		}
		return string(resultJSON)
	}
}

// ToBool convert any value to boolean.
func ToBool(v interface{}) bool {
	switch v := v.(type) {
	case string:
		str := strings.TrimSpace(v)
		result, err := strconv.ParseBool(str)
		if err != nil {
			return false
		}
		return result
	case int, int32, int64:
		return v != 0
	default:
		return false
	}
}

// ToInt converts any value to int
func ToInt(v interface{}) int {
	switch v := v.(type) {
	case string:
		str := strings.TrimSpace(v)
		result, err := strconv.Atoi(str)
		if err != nil {
			return 0
		}
		return result
	case int:
		return v
	case int32:
		return int(v)
	case int64:
		return int(v)
	case float32:
		return int(v)
	case float64:
		return int(v)
	case []byte:
		result, err := strconv.Atoi(string(v))
		if err != nil {
			return 0
		}
		return result
	case bool:
		if v {
			return 1
		}
		return 0
	case time.Month:
		return int(v)
	default:
		return 0
	}
}

// ToInt64 converts any value to int64
func ToInt64(v interface{}) int64 {
	switch v := v.(type) {
	case string:
		str := strings.TrimSpace(v)
		result, err := strconv.Atoi(str)
		if err != nil {
			return 0
		}
		return int64(result)
	case int:
		return int64(v)
	case int32:
		return int64(v)
	case int64:
		return v
	case float32:
		return int64(v)
	case float64:
		return int64(v)
	case []byte:
		result, err := strconv.Atoi(string(v))
		if err != nil {
			return 0
		}
		return int64(result)
	case bool:
		if v {
			return 1
		}
		return 0
	case uint32:
		return int64(v)
	default:
		return 0
	}
}

// ToInt8 converts any value to int8
func ToInt8(v interface{}) int8 {
	switch v := v.(type) {
	case string:
		str := strings.TrimSpace(v)
		result, err := strconv.Atoi(str)
		if err != nil {
			return 0
		}
		return int8(result)
	case int:
		return int8(v)
	case int32:
		return int8(v)
	case int64:
		return int8(v)
	case float32:
		return int8(v)
	case float64:
		return int8(v)
	case []byte:
		result, err := strconv.Atoi(string(v))
		if err != nil {
			return 0
		}
		return int8(result)
	case bool:
		if v {
			return 1
		}
		return 0
	default:
		return 0
	}
}

// ToInt32 converts any value to int32
func ToInt32(v interface{}) int32 {
	switch v := v.(type) {
	case string:
		str := strings.TrimSpace(v)
		result, err := strconv.Atoi(str)
		if err != nil {
			return 0
		}
		return int32(result)
	case int:
		return int32(v)
	case int32:
		return v
	case int64:
		return int32(v)
	case float32:
		return int32(v)
	case float64:
		return int32(v)
	case []byte:
		result, err := strconv.Atoi(string(v))
		if err != nil {
			return 0
		}
		return int32(result)
	case bool:
		if v {
			return 1
		}
		return 0
	default:
		return 0
	}
}

// ToFloat64 convert any value to float64.
func ToFloat64(v interface{}) float64 {
	switch v := v.(type) {
	case string:
		str := strings.TrimSpace(v)
		result, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return 0
		}
		return result
	case int:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case float32:
		return float64(v)
	case float64:
		return v
	case []byte:
		result, err := strconv.ParseFloat(string(v), 64)
		if err != nil {
			return 0
		}
		return result
	case json.RawMessage:
		var num float64
		err := json.Unmarshal(v, &num)
		if err != nil {
			return 0
		}
		return num
	default:
		return float64(0)
	}
}

// ToArrayString converts any value to []string
func ToArrayString(v interface{}) []string {
	switch v := v.(type) {
	case []interface{}:
		var result []string
		for _, val := range v {
			result = append(result, ToStr(val))
		}
		return result
	case []int64:
		var result []string
		for _, val := range v {
			result = append(result, strconv.FormatInt(val, 10))
		}
		return result
	case []string:
		return v
	default:
		return []string{}
	}
}

func ToTime(v interface{}) time.Time {
	switch v := v.(type) {
	case string:
		str := strings.TrimSpace(v)
		result, err := time.Parse(time.RFC3339, str)
		if err != nil {
			return time.Time{}
		}
		return result
	case time.Time:
		return v
	default:
		return time.Time{}
	}
}
