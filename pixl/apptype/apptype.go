package apptype

import (
	"image/color"

	"fyne.io/fyne/v2"
)

type BrushType = int

type PxCanvasConfig struct {
	DrawingArea  fyne.Size
	CanvasOffset fyne.Position
	PxRow, PxCol int
	PxSize       int
}

type State struct {
	BrushColor     color.Color
	BrushType      int
	SwatchSelected int
	FilePath       string
}

func (state *State) SetFilePath(path string) {
	state.FilePath = path
}

type Widget interface {
	fyne.CanvasObject
	CreateRenderer() fyne.WidgetRenderer
}

type WidgetRenderer interface {
	Destroy()
	Layout(size fyne.Size)
	MinSize()
	Objects() []fyne.CanvasObject
	Refresh()
}
