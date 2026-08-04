// Package version holds build metadata for Inspired Trek73.
//
// Version and BuildDate are intended to be overridden at build time via
// linker flags, e.g.:
//
//	go build -ldflags "-X github.com/appliedinspiration/inspired_trek73/internal/version.Version=1.0.0 \
//	    -X github.com/appliedinspiration/inspired_trek73/internal/version.BuildDate=2026-01-01" ./cmd/inspired_trek73
package version

import "runtime"

// These variables are set at build time via -ldflags where possible.
// Sensible defaults are provided for local `go run`/`go build` usage.
var (
	// Version is the Inspired Trek73 release version.
	Version = "dev"
	// BuildDate is the date the binary was built, in YYYY-MM-DD form.
	BuildDate = "unknown"
	// GitCommit is the short commit hash the binary was built from.
	GitCommit = "unknown"
)

// GoVersion returns the Go compiler version used to build the binary.
func GoVersion() string {
	return runtime.Version()
}

// String returns a multi-line human-readable version report.
func String() string {
	return "Inspired Trek73 " + Version + "\n" +
		"Build date:  " + BuildDate + "\n" +
		"Git commit:  " + GitCommit + "\n" +
		"Go compiler: " + GoVersion() + "\n" +
		"Platform:    " + runtime.GOOS + "/" + runtime.GOARCH
}
