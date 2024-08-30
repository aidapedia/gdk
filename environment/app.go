package environment

import (
	"os"
)

func GetAppEnvironment() string {
	if env := os.Getenv("APP_ENV"); env != "" {
		return env
	}
	return Development
}
