// Package version holds build metadata that is injected at link time via
// -ldflags, e.g.:
//
//	go build -ldflags "-X github.com/day0ops/preissuance-ext-authz/pkg/version.Version=1.2.3"
package version

import "fmt"

// Version, Commit and BuildDate are overridden at build time via -ldflags.
// Their zero values are used for unreleased/local builds.
var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

// String returns a human-readable summary of the build metadata, e.g.
// "preissuance-ext-authz version 0.1.0 (commit abc123, built 2026-01-01)".
func String() string {
	return fmt.Sprintf("preissuance-ext-authz version %s (commit %s, built %s)", Version, Commit, BuildDate)
}
