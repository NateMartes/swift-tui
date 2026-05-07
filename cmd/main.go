package main

import (
	"context"
	"fmt"
	"os"
	"github.com/NateMartes/go-swift-tui/internal/swift"
	"github.com/NateMartes/go-swift-tui/pkg/errors"
	"github.com/NateMartes/go-swift-tui/pkg/util"
)

// Get arguments for tempauth login as
// hostname, port, username, password, use_https
func GetTempAuthArgs() (string, int, string, string, bool) {
	username, err := util.UsernameVal()
	if util.UsernameSupplied() && err != nil {
		util.LogFatal(err.Error(), errors.ARGUMENT_ERROR)
	}
	password, err := util.PasswordVal()
	if util.PasswordSupplied() && err != nil {
		util.LogFatal(err.Error(), errors.ARGUMENT_ERROR)
	}
	swiftHostname, err := util.SwiftHostnameVal()
	if util.SwiftHostnameSupplied() && err != nil {
		util.LogFatal(err.Error(), errors.ARGUMENT_ERROR)
	}
	swiftPort, err := util.SwiftPortVal()
	if util.SwiftPortSupplied() && err != nil {
		util.LogFatal(err.Error(), errors.ARGUMENT_ERROR)
	}
	useHTTPS, err := util.NoHTTPSVal()
	if util.NoHTTPSSupplied() && err != nil {
		util.LogFatal(err.Error(), errors.ARGUMENT_ERROR)
	}

	return swiftHostname, swiftPort, username, password, useHTTPS
}

func main() {
	
	util.SetupLogger()
	util.ParseArgs()

	// Get debug and help arguments
	val, err := util.HelpVal()
	if util.HelpSupplied() && err != nil {
		util.LogFatal(err.Error(), errors.ARGUMENT_ERROR)
	}
	if val {
		fmt.Println(util.UsageMsg())
		os.Exit(errors.SUCCESS)
	}

	val, err = util.DebugVal()
	if util.DebugSupplied() && err != nil {
		util.LogFatal(err.Error(), errors.ARGUMENT_ERROR)
	}
	if val {
		util.SetDebugLogging(true)
	}

	// use tempauth if supplied
	usingTempAuth, err := util.TempAuthVal()
	if util.TempAuthSupplied() && err != nil {
		util.LogFatal(err.Error(), errors.ARGUMENT_ERROR)
	}
	if usingTempAuth {
		hostname, port, username, password, useHttps := GetTempAuthArgs()
		swift.GetTempauthClient(context.Background(), hostname, port, username, password, useHttps)
	}
}