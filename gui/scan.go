package gui

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

	scanButton := widget.NewButton("scan", nil)

	scanActivity := widget.NewActivity()
	scanActivity.Hide()
	scanActivity.Stop()

	content := container.NewGridWithColumns(2,
		portEntry, scanButton,
	)

	scanButton.OnTapped = func() {
		set := *settings
		content.Objects[1] = scanActivity
		scanActivity.Start()
		scanActivity.Show()
		content.Refresh()
		go func() {
			set.Scan()

			fyne.Do(func() {
				scanActivity.Stop()
				scanActivity.Hide()
				content.Objects[1] = scanButton
				content.Refresh()
			})

		}()
	}

	return content
}

func checkNumber(s string) error {
	for _, char := range s {
		if !unicode.IsNumber(char) {
			return fmt.Errorf("should contain only numbers")
		}
	}
	return nil
}
