package singleton

import (
	"log"
	"os"
)

var (
	CF_ACCOUNT_ID     = os.Getenv("CF_ACCOUNT_ID")
	CF_D1_DATABASE_ID = os.Getenv("CF_D1_DATABASE_ID")
	D1_API_TOKEN      = os.Getenv("CF_D1_API_TOKEN")
	COUNTERS_API_KEY  = os.Getenv("COUNTERS_API_KEY")
	PORT              = os.Getenv("PORT")
)

func ValidateRequiredEnv() {
	required := map[string]string{
		"CF_ACCOUNT_ID":     CF_ACCOUNT_ID,
		"CF_D1_DATABASE_ID": CF_D1_DATABASE_ID,
		"CF_D1_API_TOKEN":   D1_API_TOKEN,
	}
	missing := false
	for name, val := range required {
		if val == "" {
			log.Printf("FATAL: required env var %s is not set", name)
			missing = true
		}
	}
	if missing {
		log.Fatal("Aborting: one or more required environment variables are missing")
	}
}
