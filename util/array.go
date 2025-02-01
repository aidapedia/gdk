package util

import "strconv"

// ArrayStringToString function used to convert array string to string
func ArrayStringToString(arr []string, delimiter string) string {
	var str string
	for _, v := range arr {
		str += v + delimiter
	}
	return str[:len(str)-1]
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

// ToArrayInt64 converts any value to []int64
func ToArrayInt64(v interface{}) []int64 {
	switch v := v.(type) {
	case []interface{}:
		var result []int64
		for _, val := range v {
			result = append(result, ToInt64(val))
		}
		return result
	case []int64:
		return v
	default:
		return []int64{}
	}
}

// ToArrayInt32 converts any value to []int32
func ToArrayInt32(v interface{}) []int32 {
	switch v := v.(type) {
	case []interface{}:
		var result []int32
		for _, val := range v {
			result = append(result, ToInt32(val))
		}
		return result
	case []int32:
		return v
	default:
		return []int32{}
	}
}

// ToArrayInt8 converts any value to []int8
func ToArrayInt8(v interface{}) []int8 {
	switch v := v.(type) {
	case []interface{}:
		var result []int8
		for _, val := range v {
			result = append(result, ToInt8(val))
		}
		return result
	case []int8:
		return v
	default:
		return []int8{}
	}
}
