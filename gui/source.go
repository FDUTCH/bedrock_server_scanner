package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/FDUTCH/bedrock_scanner/internal/util"
	"github.com/FDUTCH/bedrock_scanner/scanner"
	"net/netip"
)

func NewSourceSelector(settings *scanner.Settings, w fyne.Window) fyne.CanvasObject {
	sourceEntry := widget.NewEntry()
	sourceEntry.Validator = prefixValidator
	sourceEntry.OnChanged = func(s string) { settings.Source = util.NewPrefixReader(s) }
	sourceEntry.PlaceHolder = "0.0.0.0/0"

	box := container.NewGridWithColumns(2,
		sourceEntry, newFileSelector("select a source file", w, settings, sourceEntry),
		widget.NewLabel("source"), widget.NewLabel("file selector"),
	)

	return box
}

func prefixValidator(str string) error {
	_, err := netip.ParsePrefix(str)
	return err
}

func newFileSelector(label string, w fyne.Window, settings *scanner.Settings, entry *widget.Entry) fyne.CanvasObject {
	return widget.NewButton(label, func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if settings.Source != nil {
				settings.Source.Close()
			}
			settings.Source = reader
			entry.Text = reader.URI().Name()
			entry.Refresh()
		}, w)
	})
}
