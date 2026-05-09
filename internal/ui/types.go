package ui

import (
	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"
)

var TEXT_HEADER_COLOR tcell.Color = tcell.ColorLightBlue
var TEXT_COLOR tcell.Color = tcell.ColorWhite
var TEXT_CONNECTED_COLOR tcell.Color = tcell.ColorGreen
var TEXT_DISCONNECTED_COLOR tcell.Color = tcell.ColorRed
var TEXT_ACCENT_COLOR = tcell.ColorGreen
var STATUS_BAR_COLOR tcell.Color = tcell.ColorBlack
var SELECTED_ITEM_COLOR tcell.Color = tcell.ColorLightBlue
var SELECTED_TEXT_COLOR tcell.Color = tcell.ColorBlack
var BORDER_COLOR tcell.Color = tcell.ColorWhite

// Layout holds all TUI panels for reference
type Layout struct {
	App            *tview.Application
	ContainerList  *tview.List
	ObjectTable    *tview.Table
	MetadataView   *tview.TextView
	ClusterStats   *tview.TextView
	LogView        *tview.TextView
	StatusBar      *tview.TextView
	Pages          *tview.Pages
}
