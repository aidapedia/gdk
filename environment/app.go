package environment

import (
	"os"
)

func GetAppEnvironment() string {
	if env := os.Getenv("app_env"); env != "" {
		return env
	}
	return Development
}
