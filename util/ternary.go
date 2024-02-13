package util

// TernaryEqualInt checks if currentValue is equal to expressionValue, if true, return currentValue, else return defaultVal
func TernaryEqualString(currentValue, expressionValue, defaultVal string) string {
	if currentValue == expressionValue {
		return currentValue
	}
	return defaultVal
}
