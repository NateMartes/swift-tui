package main

import (
	"fmt"
	swift "github.com/NateMartes/swift-tui/internal/swift"
	"github.com/NateMartes/swift-tui/internal/ui"
	"github.com/NateMartes/swift-tui/pkg/errors"
	"github.com/NateMartes/swift-tui/pkg/util"
	swiftSdk "github.com/ncw/swift/v2"
	"os"
)

// Gets a Swift connection depending on the arguments provided
func SetClient() *swiftSdk.Connection {

	// use clouds.yaml file if supplied
	cloudsFile, err := util.CloudsFileVal()
	if util.CloudsFileSupplied() && err != nil {
		util.LogFatal(err.Error(), errors.ARGUMENT_ERROR)
	}

	// default to tempauth if no clouds.yaml file supplied
	if util.CloudsFileSupplied() {
		util.LogDebug("Using clouds file auth")
		return swift.SetClientFromCloudsFile(cloudsFile)
	} else {
		util.LogDebug("Using Swift tempauth middleware")
		return swift.SetClientFromTempauth()
	}
}

func main() {

	util.SetupLogger()
	util.ParseArgs()

	// Get help argument if any
	val, err := util.HelpVal()
	if util.HelpSupplied() && err != nil {
		util.LogFatal(err.Error(), errors.ARGUMENT_ERROR)
	}
	if val {
		fmt.Println(util.UsageMsg())
		os.Exit(errors.SUCCESS)
	}

	// Get debug argument if any
	val, err = util.DebugVal()
	if util.DebugSupplied() && err != nil {
		util.LogFatal(err.Error(), errors.ARGUMENT_ERROR)
	}
	if val {
		util.SetDebugLogging(true)
	}

	client := SetClient()
	app, _ := ui.GetMainTUI(client)
	err = app.Run()
	if err != nil {
		util.LogFatal(
			fmt.Sprintf("failed to start main TUI: %v", err),
			errors.TUI_ERROR,
		)
	}
}
