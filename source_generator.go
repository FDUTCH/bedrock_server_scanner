package main

import (
	"bufio"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/FDUTCH/bedrock_scanner/scanner"
	"os"
	"strings"
)

func NewSourceFileGenerator(w fyne.Window) fyne.CanvasObject {
	asNumberEntry := widget.NewEntry()
	asNumberEntry.PlaceHolder = "AS12345"
	asNumberEntry.Validator = verifyASNumber

	outputFileEntry := widget.NewEntry()
	outputFileEntry.PlaceHolder = "out.txt"

	generateActivity := widget.NewActivity()
	generateActivity.Hide()
	generateButton := widget.NewButton("generate source file", func() {

		if err := verifyASNumber(asNumberEntry.Text); err != nil {
			dialog.ShowError(err, w)
			return
		}

		if outputFileEntry.Text == "" {
			dialog.ShowError(fmt.Errorf("output can not be empty"), w)
		}

		generateActivity.Show()
		generateActivity.Start()
		go func() {
			generateSourceFile(asNumberEntry.Text, outputFileEntry.Text, w)
			fyne.Do(func() {
				generateActivity.Stop()
				generateActivity.Hide()
			})
		}()
	})

	return container.NewVBox(
		container.NewGridWithColumns(2,
			asNumberEntry, outputFileEntry,
			widget.NewLabel("AS number"), widget.NewLabel("output file"),
			generateButton, generateActivity,
		),
	)
}

func verifyASNumber(s string) error {
	if len(s) < 3 {
		return fmt.Errorf("AS is to short")
	}
	if !strings.HasPrefix(s, "AS") {
		return fmt.Errorf("AS should start with AS")
	}
	if err := checkNumber(s[2:]); err != nil {
		return fmt.Errorf("AS should contain only numbers")
	}
	return nil
}

func generateSourceFile(as string, out string, w fyne.Window) {
	file, err := os.OpenFile(out, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		dialog.ShowError(err, w)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, prefix := range scanner.GetRangesScraping(as) {
		writer.WriteString(prefix.String() + "\n")
	}
	_ = writer.Flush()
	fyne.Do(
		func() {
			dir, err := os.Getwd()
			if err != nil {
				panic(err)
			}
			fyne.CurrentApp().SendNotification(fyne.NewNotification("source file is ready", dir+"/"+out))
		},
	)
}
