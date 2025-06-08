package error

import (
	"errors"
	"log"
	"testing"
)

func TestCaller(t *testing.T) {
	e := UseErrorV2()
	log.Println(e.(*Error).Caller())
}

func UseError() error {
	return New(errors.New("abc"))
}

func UseErrorV2() error {
	return UseError()
}
