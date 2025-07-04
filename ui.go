package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/FDUTCH/bedrock_scanner/scanner"
	"io"
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
		// source selector.
		widget.NewLabel("Source Settings"),
		NewSourceSelector(&settings, w),
		widget.NewLabel("Performance Settings"),
		NewPerformanceSetup(&settings),
	)
}

func NewLogWindow(name string) io.StringWriter {
	logGrid := widget.NewTextGrid()
	logGrid.Scroll = fyne.ScrollVerticalOnly
	logGrid.ShowWhitespace = true

	w := fyne.CurrentApp().NewWindow(name)
	w.Resize(fyne.NewSize(400, 300))
	w.SetContent(logGrid)
	w.SetFixedSize(true)
	w.Show()
	return &writer{grid: logGrid}
}

type writer struct {
	grid *widget.TextGrid
}

func (w *writer) WriteString(s string) (n int, err error) {
	fyne.Do(func() {
		w.grid.Append(s)
	})
	return len(s), nil
}
