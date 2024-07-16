package ui

import "fyne.io/fyne/v2/container"

func Setup(app *AppInit) {
	swatches := BuildSwatches(app)

	colorPicker := SetupColorPicker(app)

	appLayout := container.NewBorder(nil, swatches, nil, colorPicker)

	app.PixlWindow.SetContent(appLayout)
}
