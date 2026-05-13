package util

import (
	"bytes"
	"fmt"
	"github.com/spf13/pflag"
	"os"
)

// Get the usage message for action-to-gitlab-ci
func UsageMsg() string {
	var buf bytes.Buffer
	pflag.CommandLine.SetOutput(&buf)
	var flagSummary string
	pflag.VisitAll(func(f *pflag.Flag) {
		if f.Value.Type() == "bool" {
			flagSummary += fmt.Sprintf(" [-%s|--%s]\n", f.Shorthand, f.Name)
		} else {
			flagSummary += fmt.Sprintf(" [-%s|--%s <%s>]\n", f.Shorthand, f.Name, f.Name)
		}
	})
	fmt.Fprintf(&buf, "\nusage: %s\n%s\n", os.Args[0], flagSummary)
	fmt.Fprintf(&buf, "A Terminal User Interface (TUI) for OpenStack Swift clusters.\n\n")
	fmt.Fprintf(&buf, "Arguments:\n")
	pflag.PrintDefaults()
	pflag.CommandLine.SetOutput(os.Stderr)
	return buf.String()
}

// Parse arguments from the command line
func ParseArgs() {

	// Args
	pflag.BoolP("debug", "d", false, "Turn on debug messaging")
	pflag.StringP("debug-file", "f", "", "A filename to place debug logging into if debug logging is enabled (default is a random filename in the current directory)")
	pflag.BoolP("help", "h", false, "Display help message")
	pflag.StringP("clouds-file-path", "c", "", "Use an OpenStackClient (aka OSC) clouds.yaml file to login to OpenStack Swift")
	pflag.StringP("cloud-name", "n", "", "Cloud to use in the clouds.yaml file to connect to OpenStack Swift")
	pflag.StringP("username", "u", "", "The username to log in with to OpenStack Swift's tempauth middleware")
	pflag.StringP("api-key", "a", "", "The api-key/password to log in with to OpenStack Swift's tempauth middleware")
	pflag.IntP("swift-port", "p", 8080, "The port to use to connect to OpenStack Swift")
	pflag.StringP("swift-hostname", "s", "localhost", "The hostname to use to connect to OpenStack Swift")
	pflag.BoolP("no-https", "l", false, "Signal to not use HTTPS for connecting to OpenStack Swift")
	pflag.Parse()
}

// Get the value of the help argument
func HelpVal() (bool, error) {
	v, err := pflag.CommandLine.GetBool("help")
	return v, err
}

// Get the value of the debug argument
func DebugVal() (bool, error) {
	v, err := pflag.CommandLine.GetBool("debug")
	return v, err
}

// Get the value of the debug-file argument
func DebugFileVal() (string, error) {
	v, err := pflag.CommandLine.GetString("debug-file")
	return v, err
}

// Get the value of the username argument
func UsernameVal() (string, error) {
	v, err := pflag.CommandLine.GetString("username")
	return v, err
}

// Get the value of the api-key argument
func ApiKeyVal() (string, error) {
	v, err := pflag.CommandLine.GetString("api-key")
	return v, err
}

// Get the value of the swift-port argument
func SwiftPortVal() (int, error) {
	v, err := pflag.CommandLine.GetInt("swift-port")
	return v, err
}

// Get the value of the swift-hostname argument
func SwiftHostnameVal() (string, error) {
	v, err := pflag.CommandLine.GetString("swift-hostname")
	return v, err
}

// Get the value of the no-https argument
func NoHTTPSVal() (bool, error) {
	v, err := pflag.CommandLine.GetBool("no-https")
	return v, err
}

// Get the value of the clouds-file-path argument
func CloudsFileVal() (string, error) {
	v, err := pflag.CommandLine.GetString("clouds-file-path")
	return v, err
}

// Get the value of the cloud-name argument
func CloudNameVal() (string, error) {
	v, err := pflag.CommandLine.GetString("cloud-name")
	return v, err
}

// Check if the cloud-name argument was supplied
func CloudNameSupplied() bool {
	return pflag.CommandLine.Changed("cloud-name")
}

// Check if the clouds-file-path argument was supplied
func CloudsFileSupplied() bool {
	return pflag.CommandLine.Changed("clouds-file-path")
}

// Check if the debug argument was supplied
func DebugSupplied() bool {
	return pflag.CommandLine.Changed("debug")
}

// Check if the debug-file argument was supplied
func DebugFileSupplied() bool {
	return pflag.CommandLine.Changed("debug-file")
}

// Check if the help argument was supplied
func HelpSupplied() bool {
	return pflag.CommandLine.Changed("help")
}

// Check if the username argument was supplied
func UsernameSupplied() bool {
	return pflag.CommandLine.Changed("username")
}

// Check if the api-key argument was supplied
func ApiKeySupplied() bool {
	return pflag.CommandLine.Changed("api-key")
}

// Check if the swift-port argument was supplied
func SwiftPortSupplied() bool {
	return pflag.CommandLine.Changed("swift-port")
}

// Check if the swift-hostname argument was supplied
func SwiftHostnameSupplied() bool {
	return pflag.CommandLine.Changed("swift-hostname")
}

// Check if the no-https argument was supplied
func NoHTTPSSupplied() bool {
	return pflag.CommandLine.Changed("no-https")
}
