// Package config loads runtime configuration for the preissuance-ext-authz
// service from environment variables.
package config

import (
	"os"
	"strings"
)

// Environment variable names recognized by Load.
const (
	envPort              = "PORT"
	envAllowedPrincipals = "ALLOWED_PRINCIPALS"
)

// Default values used when the corresponding environment variable is unset.
const defaultPort = "9002"

// Config holds the runtime configuration for the preissuance-ext-authz
// service.
type Config struct {
	// Port is the gRPC listen port for the ext_authz service.
	Port string
	// AllowedPrincipals is the set of CheckRequest.Attributes.Source.Principal
	// values allowed through the pre-issuance gate.
	AllowedPrincipals map[string]bool
}

// Load reads configuration from environment variables, applying defaults
// where documented.
func Load() *Config {
	return &Config{
		Port:              getEnvOrDefault(envPort, defaultPort),
		AllowedPrincipals: parsePrincipals(os.Getenv(envAllowedPrincipals)),
	}
}

// parsePrincipals splits a comma-separated list of principals into a set,
// trimming whitespace and dropping empty entries.
func parsePrincipals(raw string) map[string]bool {
	principals := make(map[string]bool)
	for _, p := range strings.Split(raw, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			principals[p] = true
		}
	}
	return principals
}

// getEnvOrDefault returns the value of the environment variable named key,
// or fallback if it is unset or empty.
func getEnvOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
