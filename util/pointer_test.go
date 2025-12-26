package util

import (
	"testing"
)

func TestPtr(t *testing.T) {
	strVal := "test"
	strPtr := Ptr(strVal)
	if *strPtr != strVal {
		t.Errorf("Expected %v, got %v", strVal, *strPtr)
	}

	intVal := 123
	intPtr := Ptr(intVal)
	if *intPtr != intVal {
		t.Errorf("Expected %v, got %v", intVal, *intPtr)
	}
}

func TestVal(t *testing.T) {
	var strPtr *string
	if got := Val(strPtr); got != "" {
		t.Errorf("Expected empty string, got %v", got)
	}

	strVal := "test"
	strPtr = &strVal
	if got := Val(strPtr); got != strVal {
		t.Errorf("Expected %v, got %v", strVal, got)
	}
}

func TestValOr(t *testing.T) {
	var strPtr *string
	def := "default"
	if got := ValOr(strPtr, def); got != def {
		t.Errorf("Expected %v, got %v", def, got)
	}

	strVal := "test"
	strPtr = &strVal
	if got := ValOr(strPtr, def); got != strVal {
		t.Errorf("Expected %v, got %v", strVal, got)
	}
}
