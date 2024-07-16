package ui

import (
	"image/color"

	"pixl/swatch"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func BuildSwatches(app *AppInit) *fyne.Container {
	canvasSwatches := make([]fyne.CanvasObject, 0, 4)

	for i := 0; i < cap(app.Swatches); i++ {
		inititialColor := color.NRGBA{255, 255, 255, 255}
		s := swatch.NewSwatch(app.State, inititialColor, i, func(defSwatch *swatch.Swatch) {
			for l := 0; l < cap(app.Swatches); l++ {
				app.Swatches[l].Selected = false
				canvasSwatches[i].Refresh()
			}
			app.State.SwatchSelected = defSwatch.SwatchIndex
			app.State.BrushColor = defSwatch.Color
		})

		if i == 0 {
			s.Selected = true
			app.State.SwatchSelected = 0
			s.Refresh()
		}
		app.Swatches = append(app.Swatches, s)
		canvasSwatches = append(canvasSwatches, s)

	}

	return container.NewGridWrap( canvasSwatches...,fyne.NewSize(20, 20))
}
