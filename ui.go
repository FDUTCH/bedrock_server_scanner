package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/FDUTCH/bedrock_scanner/scanner"
)

func AddMainContents(w fyne.Window) {
	w.SetContent(
		container.NewVBox(
			NewSettingsScreen(w),
		),
	)
}

func NewSettingsScreen(w fyne.Window) fyne.CanvasObject {
	settings := scanner.Settings{}

	return container.NewVBox(
		NewSourceSelector(&settings, w),
		NewPerformanceSetup(&settings),
		NewScanMenu(&settings),
		NewOutPut(&settings),
		widget.NewSeparator(),
		NewSourceFileGenerator(w),
	)
}

func NewLogWindow(name string) *widget.TextGrid {
	logGrid := widget.NewTextGrid()
	logGrid.ShowWhitespace = true
	logGrid.Scroll = fyne.ScrollBoth

	w := fyne.CurrentApp().NewWindow(name)

	w.SetContent(logGrid)
	w.Show()
	return logGrid
}

type name struct {
}
