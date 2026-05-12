package ui

import (
	"context"
	"fmt"
	"github.com/NateMartes/swift-tui/pkg/util"
	"github.com/gdamore/tcell/v2"
	swiftSdk "github.com/ncw/swift/v2"
	"github.com/rivo/tview"
	"time"
)

// GetMainTUI builds and returns the fully wired TUI application
func GetMainTUI(client *swiftSdk.Connection) (*tview.Application, *Layout) {
	app := tview.NewApplication()
	layout := BuildLayout(client, app)
	return app, layout
}

// EndMainTUI stops the TUI application
func EndMainTUI(app *tview.Application) {
	app.Stop()
}

// Convert tcell Color type to hex code. Useful for tview dynamic colors
func ColorToTag(c tcell.Color) string {
	return fmt.Sprintf("[#%06x]", c.Hex())
}

// Convert an integer of bytes to a gigabytes
func BytesToGB(a int64) float64 {
	return float64(a) / 1_000_000_000.0
}

// Convert an integer of bytes to a megabytes
func BytesToMB(a int64) float64 {
	return float64(a) / 1_000_000.0
}

// Convert an integer of bytes to a kilobytes
func BytesToKB(a int64) float64 {
	return float64(a) / 1_000.0
}

// Format bytes into either gigabytes, megabytes, kilobytes.
// Returns a string representation of the value that was converted
func FormatBytes(a int64) (float64, string) {
	if a >= 1_000_000_000 {
		return BytesToGB(a), "GB"
	} else if a >= 1_000_000 {
		return BytesToGB(a), "MB"
	} else {
		return BytesToKB(a), "KB"
	}
}

// Gets the header object of the TUI
func GetHeader() *tview.TextView {
	header := tview.NewTextView().
		SetText("OpenStack Swift TUI").
		SetTextColor(TEXT_HEADER_COLOR).
		SetDynamicColors(true)
	header.SetBorder(true).SetBorderColor(BORDER_COLOR)
	header.SetTextAlign(tview.AlignCenter)
	return header
}

// Gets the cluster status pane of the TUI using a Swift connection
func GetClusterStats(client *swiftSdk.Connection) *tview.TextView {
	clusterStats := tview.NewTextView().
		SetDynamicColors(true)
	clusterStats = UpdateClusterStats(client, clusterStats)
	clusterStats.
		SetTitle(" Cluster Stats ").
		SetTitleColor(TEXT_HEADER_COLOR).
		SetBorder(true).
		SetBorderColor(BORDER_COLOR)
	return clusterStats
}

// Updates the cluster status pane using a Swift connection
func UpdateClusterStats(client *swiftSdk.Connection, clusterStats *tview.TextView) *tview.TextView {

	connected := false
	endpointStatus := "Disconnected"
	storageUrl := "Unknow Storage URL"
	containerCount := int64(0)
	objectCount := int64(0)
	accountSize := float64(0.0)
	accountSizeFormat := "GB"

	account, _, err := client.Account(context.Background())
	if err != nil {
		util.LogError(
			fmt.Sprintf("Failed to get account info from %s as %s: %s",
				client.StorageUrl,
				client.UserName,
				err.Error(),
			),
		)
	} else {
		connected = true
		endpointStatus = "Connected"
		storageUrl = client.StorageUrl
		containerCount = account.Containers
		objectCount = account.Objects
		accountSize, accountSizeFormat = FormatBytes(account.BytesUsed)
	}

	headerColorTag := ColorToTag(TEXT_HEADER_COLOR)
	statusColorTag := ColorToTag(TEXT_DISCONNECTED_COLOR)
	if connected {
		statusColorTag = ColorToTag(TEXT_CONNECTED_COLOR)
	}
	textColorTag := ColorToTag(TEXT_COLOR)
	objectSizeTag := ColorToTag(TEXT_ACCENT_COLOR)

	clusterStats = clusterStats.SetText(
		fmt.Sprintf(
			"%sEndpoint %s(● %s)%s:%s %s\n%sAccount:%s %s\n%sContainers:%s %d\n%sObjects:%s %s%d   %sTotal Size:%s %s%.3f %s",
			headerColorTag,
			statusColorTag, endpointStatus, headerColorTag,
			textColorTag, storageUrl,
			headerColorTag, textColorTag,
			client.UserName,
			headerColorTag, textColorTag,
			containerCount,
			headerColorTag, textColorTag,
			textColorTag, objectCount,
			headerColorTag, textColorTag,
			objectSizeTag, accountSize, accountSizeFormat,
		),
	)
	return clusterStats
}

// Get the container pane using a Swift connection, returing the pane and the currently selected container
func GetContainerList(client *swiftSdk.Connection) (*tview.List, string) {
	containerList := tview.NewList().
		ShowSecondaryText(true).
		SetHighlightFullLine(true).
		SetSelectedBackgroundColor(SELECTED_ITEM_COLOR).
		SetSelectedTextColor(SELECTED_TEXT_COLOR)
	containerList.
		SetTitle(" Containers ").
		SetTitleColor(TEXT_HEADER_COLOR).
		SetBorder(true).
		SetBorderColor(BORDER_COLOR)

	return UpdateContainerList(client, containerList)
}

// Updates the container pane using a Swift connection, returing the pane and the currently selected container
func UpdateContainerList(client *swiftSdk.Connection, containerList *tview.List) (*tview.List, string) {

	containerFormatString := "%d objects · %.3f %s"
	containers := []struct{ name, meta string }{
		{"No Containers Found", ""},
	}
	selectedContainer := NO_CONTAINER_SELECTED

	containersResult, err := client.ContainersAll(context.Background(), nil)
	if err != nil {
		util.LogError(
			fmt.Sprintf("Failed to get account info from %s as %s: %s",
				client.StorageUrl,
				client.UserName,
				err.Error(),
			),
		)
	} else {
		if len(containersResult) > 0 {

			// set the first container as the selected container
			selectedContainer = containersResult[0].Name

			containers = []struct{ name, meta string }{}
			for _, c := range containersResult {
				size, sizeFormat := FormatBytes(c.Bytes)
				containers = append(containers,
					struct{ name, meta string }{
						c.Name,
						fmt.Sprintf(containerFormatString, c.Count, size, sizeFormat),
					},
				)
			}
		}
	}

	for _, c := range containers {
		containerList.AddItem(c.name, c.meta, 0, nil)
	}
	return containerList, selectedContainer
}

// Gets the object table pane using a Swift connection, returing the pane using some currently selected container
func GetObjectTable(client *swiftSdk.Connection, selectedContainer string) (*tview.Table, string) {

	objectTable := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false).
		SetSelectedStyle(tcell.StyleDefault.
			Background(SELECTED_ITEM_COLOR).
			Foreground(SELECTED_TEXT_COLOR))
	objectTable.
		SetTitle(" Objects ").
		SetTitleColor(TEXT_HEADER_COLOR).
		SetBorder(true).
		SetBorderColor(BORDER_COLOR)

	return UpdateObjectTable(client, objectTable, selectedContainer)
}

// Updates the object table pane using a Swift connection, returing the pane using some currently selected container
func UpdateObjectTable(client *swiftSdk.Connection, objectTable *tview.Table, selectedContainer string) (*tview.Table, string) {

	selectedObject := NO_OBJECT_SELECTED
	headers := []string{"Object Name", "Size", "Last Modified", "Content-Type"}
	for col, h := range headers {
		objectTable.SetCell(0, col, tview.NewTableCell(h).
			SetTextColor(TEXT_HEADER_COLOR).
			SetAttributes(tcell.AttrBold).
			SetSelectable(false).
			SetExpansion(1))
	}

	if selectedContainer == NO_CONTAINER_SELECTED {
		return objectTable, selectedObject
	}
	objects, err := client.ObjectsAll(context.Background(), selectedContainer, nil)
	if err != nil {
		util.LogError(
			fmt.Sprintf("Failed to get objects %s as %s: %s",
				client.StorageUrl,
				client.UserName,
				err.Error(),
			),
		)
		return objectTable, selectedObject
	}

	selectedObject = objects[0].Name
	rows := [][]string{}
	for _, o := range objects {
		size, sizeFormat := FormatBytes(o.Bytes)
		rows = append(rows,
			// order matches with top row
			[]string{
				o.Name,
				fmt.Sprintf("%.3f %s", size, sizeFormat),
				o.LastModified.Format(time.RFC1123),
				o.ContentType,
			},
		)
	}

	for r, row := range rows {
		for col, val := range row {
			color := TEXT_COLOR
			if col == 1 {
				color = TEXT_ACCENT_COLOR
			}
			objectTable.SetCell(r+1, col, tview.NewTableCell(val).
				SetTextColor(color).
				SetExpansion(1))
		}
	}

	return objectTable, selectedObject
}

// Gets the metadata pane using a Swift connection, returing the pane using some currently selected object
func GetMetadataView(client *swiftSdk.Connection, selectedObject string) *tview.TextView {
	metadataView := tview.NewTextView().
		SetDynamicColors(true)
	metadataView = UpdateMetadataView(client, metadataView, selectedObject)
	metadataView.
		SetTitle(" Object Metadata ").
		SetTitleColor(TEXT_HEADER_COLOR).
		SetBorder(true).
		SetBorderColor(BORDER_COLOR)
	return metadataView
}

// Updates the metadata pane using a Swift connection, returing the pane using some currently selected object
func UpdateMetadataView(client *swiftSdk.Connection, metadataView *tview.TextView, selectedObject string) *tview.TextView {
	metadataView = metadataView.SetText(
		`
		[yellow]Name:[white]         backup-2025-04-01.tar.gz
		[yellow]Size:[white]         1,181,116,006 bytes
		[yellow]ETag:[white]         d41d8cd98f00b204e9800998ecf8427e
		[yellow]Content-Type:[white] application/gzip
		[yellow]Last-Modified:[white] Tue, 01 Apr 2025 02:00:00 GMT
		[yellow]X-Object-Meta-Source:[white] cron-job
		`,
	)
	return metadataView
}

// Gets the log view for the TUI
func GetLogView() *tview.TextView {
	logView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetText("[gray]Waiting for activity...\n")
	logView.
		SetTitle(" Activity Log ").
		SetTitleColor(TEXT_HEADER_COLOR).
		SetBorder(true).
		SetBorderColor(BORDER_COLOR)
	return logView
}

// Gets the status bar for the TUI, displaying what inputs what allowed
func GetStatusBar() *tview.TextView {
	statusBar := tview.NewTextView().
		SetDynamicColors(true)
	statusBar = UpdateStatusBar(statusBar)
	statusBar.SetBackgroundColor(STATUS_BAR_COLOR)
	return statusBar
}

// Updates the status bar for the TUI
func UpdateStatusBar(statusBar *tview.TextView) *tview.TextView {
	statusBar = statusBar.SetText(" [yellow]Tab[white] Switch panel  [yellow]↑↓[white] Navigate [yellow]Q[white] Quit")
	return statusBar
}

// Builds a Layout object for the TUI
func BuildLayout(client *swiftSdk.Connection, app *tview.Application) *Layout {

	l := &Layout{}
	l.ClusterStats = GetClusterStats(client)

	var selectedContainer string = NO_CONTAINER_SELECTED
	l.ContainerList, selectedContainer = GetContainerList(client)

	var selectedObject string = NO_OBJECT_SELECTED
	l.ObjectTable, selectedObject = GetObjectTable(client, selectedContainer)

	l.MetadataView = GetMetadataView(client, selectedObject)
	l.LogView = GetLogView()
	l.StatusBar = GetStatusBar()

	header := GetHeader()
	topBar := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(header, 0, 2, false)

	rightPanel := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(l.ObjectTable, 0, 6, false).
		AddItem(l.MetadataView, 0, 4, false)

	mainContent := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(l.ContainerList, 30, 0, true).
		AddItem(rightPanel, 0, 1, false)

	root := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(topBar, 5, 0, false).
		AddItem(mainContent, 0, 1, true).
		AddItem(
			tview.NewFlex().SetDirection(tview.FlexColumn).
				AddItem(l.LogView, 0, 6, false).
				AddItem(l.ClusterStats, 0, 4, false), 6, 0, false,
		).
		AddItem(l.StatusBar, 1, 0, false)

	l.Pages = tview.NewPages().AddPage("main", root, true, true)

	app = SetupInputHandling(app, l)

	// Update metadata when a container is selected
	l.ContainerList.SetChangedFunc(func(index int, name, secondary string, shortcut rune) {
		l.ObjectTable.Clear()
		UpdateObjectTable(client, l.ObjectTable, name)
		appendLog(l.LogView, fmt.Sprintf("[cyan]Container selected:[white] %s", name))
	})

	// Update metadata when an object row is selected
	l.ObjectTable.SetSelectionChangedFunc(func(row, col int) {
		if row > 0 { // skip header row
			cell := l.ObjectTable.GetCell(row, 0)
			if cell != nil {
				appendLog(l.LogView, fmt.Sprintf("[cyan]Object focused:[white] %s", cell.Text))
				l.MetadataView = UpdateMetadataView(client, l.MetadataView, cell.Text)
			}
		}
	})

	app.SetRoot(l.Pages, true).EnableMouse(true)
	return l
}

func appendLog(v *tview.TextView, msg string) {
	fmt.Fprintf(v, "%s\n", msg)
	v.ScrollToEnd()
}

// Sets up input handling for the TUI
func SetupInputHandling(app *tview.Application, l *Layout) *tview.Application {

	panels := []tview.Primitive{l.ContainerList, l.ObjectTable, l.MetadataView, l.LogView}
	currentPanel := 0

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			currentPanel = (currentPanel + 1) % len(panels)
			app.SetFocus(panels[currentPanel])
			return nil
		case tcell.KeyBacktab:
			currentPanel = (currentPanel - 1 + len(panels)) % len(panels)
			app.SetFocus(panels[currentPanel])
			return nil
		}
		switch event.Rune() {
		case 'q', 'Q':
			app.Stop()
		}
		return event
	})

	return app
}
