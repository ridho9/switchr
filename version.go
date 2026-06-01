package main

import (
	"runtime/debug"
)

// Build-time variables, injected via -ldflags (make build).
// When not set (go install), versionString() falls back to
// runtime/debug.ReadBuildInfo to read the module version
// that Go automatically embeds in every binary.
var (
	Version = ""
	Commit  = ""
)

func versionString() string {
	if Version != "" {
		if Commit != "" {
			return Version + " (" + Commit + ")"
		}
		return Version
	}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}

	v := info.Main.Version
	for _, s := range info.Settings {
		if s.Key == "vcs.revision" {
			if len(s.Value) >= 7 {
				v += " (" + s.Value[:7] + ")"
			}
			break
		}
	}
	return v
}
