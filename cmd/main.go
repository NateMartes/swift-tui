package main

import (
	"fmt"
	"github.com/NateMartes/swift-tui/internal/eventloop"
	swift "github.com/NateMartes/swift-tui/internal/swift"
	"github.com/NateMartes/swift-tui/internal/ui"
	"github.com/NateMartes/swift-tui/pkg/errors"
	"github.com/NateMartes/swift-tui/pkg/util"
	swiftSdk "github.com/ncw/swift/v2"
	"os"
)

func SetClient(client *swiftSdk.Connection) {

	// use clouds.yaml file if supplied
	cloudsFile, err := util.CloudsFileVal()
	if util.CloudsFileSupplied() && err != nil {
		util.LogFatal(err.Error(), errors.ARGUMENT_ERROR)
	}

	// default to tempauth if no clouds.yaml file supplied
	if util.CloudsFileSupplied() {
		util.LogDebug("Using clouds file auth")
		client = swift.SetClientFromCloudsFile(cloudsFile)
	} else {
		util.LogDebug("Using Swift tempauth middleware")
		client = swift.SetClientFromTempauth()
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

	var client *swiftSdk.Connection
	SetClient(client)

	app := ui.GetMainTUI()
	eventloop.SetupEventLoop(client, app)
	err = app.Run()
	if err != nil {
		util.LogFatal(
			fmt.Sprintf("failed to start main TUI: %v", err),
			errors.TUI_ERROR,
		)
	}
}
