package api

import (
	"errors"
	"net/url"
	"strings"
)

const baseURL = "https://api.tcno.co/sw"

// BRQ fork settings.
const (
	// AppName is the user-facing product name.
	AppName = "BRQ Account Switcher"
	// UpdateRepo is the GitHub "owner/repo" that BRQ releases are published to.
	// Empty disables update checks and the self-updater, so a BRQ build never
	// pulls an upstream TcNo release over itself.
	UpdateRepo = ""
	// TelemetryEnabled gates crash reports, anonymous stats, stability ratings
	// and feedback. They post to upstream TcNo servers, so BRQ keeps them off.
	TelemetryEnabled = false
)

// ErrTelemetryDisabled is returned by senders when TelemetryEnabled is false.
var ErrTelemetryDisabled = errors.New("telemetry disabled in this build")

func UserAgent(version string) string {
	return "BRQ-Acc-Switcher/" + strings.TrimSpace(version)
}

func AnonymousStatsUploadURL() string {
	return baseURL + "/stats/"
}

// Singular. The plural form 404s and the report is dropped.
func CrashURL() string {
	return baseURL + "/crash/"
}

func StabilityRateURL() string {
	return baseURL + "/stability/rate/"
}

func FeedbackURL() string {
	return baseURL + "/feedback/"
}

func VersionCheckURL(version string, debug bool) string {
	v := url.QueryEscape(strings.TrimSpace(version))
	if debug {
		return "https://tcno.co/Projects/AccSwitcher/api?debug&v=" + v
	}
	return "https://tcno.co/Projects/AccSwitcher/api?v=" + v
}
