package ui

import (
	"context"
	"fmt"
	"github.com/NateMartes/swift-tui/pkg/util"
	"github.com/gdamore/tcell/v2"
	swiftSdk "github.com/ncw/swift/v2"
	"github.com/rivo/tview"
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

func ColorToTag(c tcell.Color) string {
	return fmt.Sprintf("[#%06x]", c.Hex())
}

func GetHeader() *tview.TextView {
	header := tview.NewTextView().
		SetText("OpenStack Swift TUI").
		SetTextColor(TEXT_HEADER_COLOR).
		SetDynamicColors(true)
	header.SetBorder(true).SetBorderColor(BORDER_COLOR)
	header.SetTextAlign(tview.AlignCenter)
	return header
}

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

func UpdateClusterStats(client *swiftSdk.Connection, clusterStats *tview.TextView) *tview.TextView {

	connected := false
	endpointStatus := "Disconnected"
	containerCount := int64(0)
	objectCount := int64(0)
	accountSizeGB := float64(0.0)

	account, _, err := client.Account(context.Background())
	if err != nil {
		util.LogError(
			fmt.Sprintf("Failed to get account info from %s as %s: %s",
				client.AuthUrl,
				client.UserName,
				err.Error(),
			),
		)
	} else {
		connected = true
		endpointStatus = "Connected"
		containerCount = account.Containers
		objectCount = account.Objects
		accountSizeGB = float64(account.BytesUsed) / 1_000_000.0
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
			"%sEndpoint %s(● %s)%s:%s %s\n%sAccount:%s %s\n%sContainers:%s %d\n%sObjects:%s %s%d   %sTotal Size:%s %s%.1f GB",
			headerColorTag,
			statusColorTag, endpointStatus, headerColorTag,
			textColorTag, client.AuthUrl,
			headerColorTag, textColorTag,
			client.UserName,
			headerColorTag, textColorTag,
			containerCount,
			headerColorTag, textColorTag,
			textColorTag, objectCount,
			headerColorTag, textColorTag,
			objectSizeTag, accountSizeGB,
		),
	)
	return clusterStats
}

func GetContainerList() *tview.List {
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

	containerList = UpdateContainerList(nil, containerList)
	return containerList
}

func UpdateContainerList(client *swiftSdk.Connection, containerList *tview.List) *tview.List {
	containers := []struct{ name, meta string }{
		{"backup-data", "42 objects · 1.2 GB"},
		{"static-assets", "318 objects · 4.8 GB"},
		{"logs-2024", "9 objects · 230 MB"},
		{"user-uploads", "1024 objects · 18 GB"},
		{"temp-cache", "7 objects · 50 MB"},
	}
	for _, c := range containers {
		containerList.AddItem(c.name, c.meta, 0, nil)
	}
	return containerList
}

func GetObjectTable() *tview.Table {

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

	objectTable = UpdateObjectTable(nil, objectTable)
	return objectTable
}

func UpdateObjectTable(client *swiftSdk.Connection, objectTable *tview.Table) *tview.Table {

	headers := []string{"Object Name", "Size", "Last Modified", "Content-Type"}
	for col, h := range headers {
		objectTable.SetCell(0, col, tview.NewTableCell(h).
			SetTextColor(TEXT_HEADER_COLOR).
			SetAttributes(tcell.AttrBold).
			SetSelectable(false).
			SetExpansion(1))
	}

	rows := [][]string{
		{"README.md", "4.2 KB", "2025-04-12 09:14", "text/markdown"},
		{"backup-2025-04-01.tar.gz", "1.1 GB", "2025-04-01 02:00", "application/gzip"},
		{"config.json", "812 B", "2025-03-28 17:45", "application/json"},
		{"logo.png", "58 KB", "2025-02-10 11:22", "image/png"},
		{"report-q1.pdf", "2.3 MB", "2025-04-05 13:00", "application/pdf"},
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

	return objectTable
}

func GetMetadataView() *tview.TextView {
	metadataView := tview.NewTextView().
		SetDynamicColors(true)
	metadataView = UpdateMetadataView(nil, metadataView)
	metadataView.
		SetTitle(" Object Metadata ").
		SetTitleColor(TEXT_HEADER_COLOR).
		SetBorder(true).
		SetBorderColor(BORDER_COLOR)
	return metadataView
}

func UpdateMetadataView(client *swiftSdk.Connection, metadataView *tview.TextView) *tview.TextView {
	metadataView = metadataView.SetText(
		`[yellow]Name:[white]         backup-2025-04-01.tar.gz
		[yellow]Size:[white]         1,181,116,006 bytes
		[yellow]ETag:[white]         d41d8cd98f00b204e9800998ecf8427e
		[yellow]Content-Type:[white] application/gzip
		[yellow]Last-Modified:[white] Tue, 01 Apr 2025 02:00:00 GMT
		[yellow]X-Object-Meta-Source:[white] cron-job`,
	)
	return metadataView
}

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

func GetStatusBar() *tview.TextView {
	statusBar := tview.NewTextView().
		SetDynamicColors(true)
	statusBar = UpdateStatusBar(nil, statusBar)
	statusBar.SetBackgroundColor(STATUS_BAR_COLOR)
	return statusBar
}

func UpdateStatusBar(client *swiftSdk.Connection, statusBar *tview.TextView) *tview.TextView {
	statusBar = statusBar.SetText(" [yellow]Tab[white] Switch panel  [yellow]↑↓[white] Navigate [yellow]Q[white] Quit")
	return statusBar
}

func BuildLayout(client *swiftSdk.Connection, app *tview.Application) *Layout {

	l := &Layout{}
	l.ClusterStats = GetClusterStats(client)
	l.ContainerList = GetContainerList()
	l.ObjectTable = GetObjectTable()
	l.MetadataView = GetMetadataView()
	l.LogView = GetLogView()
	l.StatusBar = GetStatusBar()

	header := GetHeader()
	topBar := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(header, 0, 2, false)

	rightPanel := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(l.ObjectTable, 8, 3, false).
		AddItem(l.MetadataView, 8, 0, false)

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
		UpdateObjectTable(nil, l.ObjectTable)
		appendLog(l.LogView, fmt.Sprintf("[cyan]Container selected:[white] %s", name))
	})

	// Update metadata when an object row is selected
	l.ObjectTable.SetSelectionChangedFunc(func(row, col int) {
		if row > 0 { // skip header row
			cell := l.ObjectTable.GetCell(row, 0)
			if cell != nil {
				appendLog(l.LogView, fmt.Sprintf("[cyan]Object focused:[white] %s", cell.Text))
				l.MetadataView = UpdateMetadataView(nil, l.MetadataView)
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
