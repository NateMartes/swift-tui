package util

import (
	"github.com/spf13/pflag"
	"os"
	"fmt"
	"bytes"
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
    fmt.Fprintf(&buf, "usage: %s\n%s\n", os.Args[0], flagSummary)
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
	pflag.BoolP("help", "h", false, "Display help message")
	pflag.BoolP("temp-auth", "t", false, "Signal to use the OpenStack Swift tempauth middleware for login")
	pflag.StringP("username", "u", "", "The username to log in with to OpenStack Swift")
	pflag.StringP("password", "x", "", "The password to log in with to OpenStack Swift")
	pflag.IntP("swift-port", "p", 8080, "The port to use to connect to OpenStack Swift, default is 8080")
	pflag.StringP("swift-hostname", "s", "localhost", "The hostname to use to connect to OpenStack Swift, default is localhost")
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

// Get the value of the temp-auth argument
func TempAuthVal() (bool, error) {
    v, err := pflag.CommandLine.GetBool("temp-auth")
    return v, err
}

// Get the value of the username argument
func UsernameVal() (string, error) {
    v, err := pflag.CommandLine.GetString("username")
    return v, err
}

// Get the value of the password argument
func PasswordVal() (string, error) {
    v, err := pflag.CommandLine.GetString("password")
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

// Check if the debug argument was supplied
func DebugSupplied() bool {
    return pflag.CommandLine.Changed("debug")
}

// Check if the help argument was supplied
func HelpSupplied() bool {
    return pflag.CommandLine.Changed("help")
}

// Check if the temp-auth argument was supplied
func TempAuthSupplied() bool {
    return pflag.CommandLine.Changed("temp-auth")
}

// Check if the username argument was supplied
func UsernameSupplied() bool {
    return pflag.CommandLine.Changed("username")
}

// Check if the password argument was supplied
func PasswordSupplied() bool {
    return pflag.CommandLine.Changed("password")
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
func NoHTTPSSupplied() (bool) {
    return pflag.CommandLine.Changed("no-https")
}