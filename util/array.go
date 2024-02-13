package util

// ArrayStringToString function used to convert array string to string
func ArrayStringToString(arr []string, delimiter string) string {
	var str string
	for _, v := range arr {
		str += v + delimiter
	}
	return str[:len(str)-1]
}
