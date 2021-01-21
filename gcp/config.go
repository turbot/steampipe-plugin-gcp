package gcp

import (
	"os"
)

// GetDefaultLocation :: return the default location used
func GetDefaultLocation() string {
	location := os.Getenv("GCP_REGION")
	if location == "" {
		panic("GCP_REGION must be set to use the gcp extension")
	}
	return location
}
