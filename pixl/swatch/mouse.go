package swatch

import (
	"fmt"

	"fyne.io/fyne/v2/driver/desktop"
)

func (swatch *Swatch) MouseDown(md *desktop.MouseEvent) {
	swatch.clickHandler(swatch)
	swatch.Selected = true
	swatch.Refresh()
	fmt.Println("selected", swatch)
}

func (swatch *Swatch) MouseUp(md *desktop.MouseEvent) {}
