package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/FDUTCH/bedrock_scanner/scanner"
	"os"
)

func NewOutPut(settings *scanner.Settings) fyne.CanvasObject {
	outputEntry := widget.NewEntry()
	outputEntry.PlaceHolder = "logs.txt"
	w := writer{}
	outputEntry.OnChanged = func(s string) { w.path = s }

	settings.Out = &w

	return container.NewGridWithColumns(1,
		outputEntry,
		widget.NewLabel("output"),
	)
}

type writer struct {
	path   string
	logger *widget.TextGrid
}

func (w *writer) WriteString(s string) (n int, err error) {
	fyne.Do(func() {
		if w.path != "" {
			os.WriteFile(w.path, []byte(s), 0666)
		}
		if w.logger == nil {
			w.logger = NewLogWindow("Logs")
		}
		w.logger.Append(s)
	})
	return len(s), nil
}
