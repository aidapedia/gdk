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

func GetSecretGSMProjectID() string {
	if projectID := os.Getenv("SECRET_GSM_PROJECT_ID"); projectID != "" {
		return projectID
	}
	return ""
}
