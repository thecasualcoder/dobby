package cmd

import "fmt"

var majorVersion string
var minorVersion string

// BuildVersion returns build version specified at compile time
// by default it will be 1.0.0-dev
func BuildVersion() string {
	major := majorVersion
	minor := minorVersion
	if major == "" {
		major = "1.0"
	}
	if minor == "" {
		minor = "0-dev"
	}
	return fmt.Sprintf("%s.%s", major, minor)
}

// SetBuildVersion use this to override the major and minor version
// of the application
func SetBuildVersion(major string, minor string) {
	majorVersion = major
	minorVersion = minor
}
