package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/FDUTCH/bedrock_scanner/scanner"
	"runtime"
)

func NewPerformanceSetup(settings *scanner.Settings) fyne.CanvasObject {
	numCpu := float64(runtime.NumCPU())
	maxPPS := numCpu * 10000

	packetsSlider := widget.NewSlider(1000, maxPPS)
	packetsSlider.Step = 1000
	packetsLabel := widget.NewLabel("packets per second: 1000")

	socketsSlider := widget.NewSlider(1, max(float64(runtime.NumCPU())/2, 1))
	var socketsSliderStep float64 = 1
	socketsSlider.Step = socketsSliderStep
	socketsLabel := widget.NewLabel("socket count: 1")

	packetsSlider.OnChanged = func(f float64) {
		settings.PacketsPerSecond = int(f)
		packetsLabel.SetText(fmt.Sprintf("packets per second: %d", settings.PacketsPerSecond))
		settings.NoPPSLimit = f >= maxPPS
	}

	socketsSlider.OnChanged = func(f float64) {
		settings.Sockets = int(f)
		socketsLabel.SetText(fmt.Sprintf("socket count: %d", settings.Sockets))

		packetsSlider.SetValue(min(packetsSlider.Value/socketsSliderStep*f, maxPPS))
		if packetsSlider.Value < f {
			packetsSlider.SetValue(f)
		}
		socketsSliderStep = f
	}

	return container.NewGridWithColumns(2,
		packetsSlider,
		socketsSlider,
		packetsLabel,
		socketsLabel,
	)
}
