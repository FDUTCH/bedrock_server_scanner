package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/FDUTCH/bedrock_scanner/scanner"
	"runtime"
)

func NewPerformanceSetup(settings *scanner.Settings) fyne.CanvasObject {
	var (
		numCpu               = float64(runtime.NumCPU())
		maxPPS               = numCpu * 10000
		socketsCount float64 = 1
		packetsMin   float64 = 1000
	)

	settings.Sockets = int(socketsCount)
	settings.PacketsPerSecond = int(packetsMin)

	packetsSlider := widget.NewSlider(packetsMin, maxPPS)
	packetsSlider.Step = packetsMin
	packetsLabel := widget.NewLabel("packets per second: 1000")

	socketsSlider := widget.NewSlider(socketsCount, max(float64(runtime.NumCPU())/2, 1))
	socketsSlider.Step = socketsCount
	socketsLabel := widget.NewLabel(fmt.Sprintf("socket count: %d", int(socketsCount)))

	packetsSlider.OnChanged = func(f float64) {
		settings.PacketsPerSecond = int(f)
		if f >= maxPPS {
			packetsLabel.SetText("unlimited")
		} else {
			packetsLabel.SetText(fmt.Sprintf("packets per second: %d", settings.PacketsPerSecond))
		}
		settings.NoPPSLimit = f >= maxPPS
	}

	socketsSlider.OnChanged = func(f float64) {
		settings.Sockets = int(f)
		socketsLabel.SetText(fmt.Sprintf("socket count: %d", settings.Sockets))

		packetsSlider.SetValue(min(packetsSlider.Value/socketsCount*f, maxPPS))
		if packetsSlider.Value < f {
			packetsSlider.SetValue(f)
		}
		socketsCount = f
	}

	return container.NewGridWithColumns(2,
		packetsSlider,
		socketsSlider,
		packetsLabel,
		socketsLabel,
	)
}
