package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/FDUTCH/bedrock_scanner/scanner"
	"strconv"
	"unicode"
)

func NewScanMenu(settings *scanner.Settings) fyne.CanvasObject {
	settings.Port = 19132

	portEntry := widget.NewEntry()
	portEntry.Text = "19132"
	portEntry.Validator = checkNumber
	portEntry.OnChanged = func(s string) {
		val, err := strconv.Atoi(s)
		if err == nil {
			settings.Port = val
		}
	}

	return container.NewGridWithColumns(2,
		portEntry, widget.NewButton("scan", func() {
			set := *settings
			go set.Scan()
		}),
	)
}

func checkNumber(s string) error {
	for _, char := range s {
		if !unicode.IsNumber(char) {
			return fmt.Errorf("should contain only numbers")
		}
	}
	return nil
}
