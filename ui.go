package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/FDUTCH/bedrock_scanner/scanner"
	"unicode"
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
	logGrid.Scroll = fyne.ScrollVerticalOnly
	logGrid.ShowWhitespace = true

	w := fyne.CurrentApp().NewWindow(name)
	w.Resize(fyne.NewSize(400, 300))
	w.SetContent(logGrid)
	w.SetFixedSize(true)
	w.Show()
	return logGrid
}

func checkNumber(s string) error {
	for _, char := range s {
		if !unicode.IsNumber(char) {
			return fmt.Errorf("should contain only numbers")
		}
	}
	return nil
}
