package theme

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var Monokai = tview.Theme{
	PrimitiveBackgroundColor:    tcell.NewHexColor(0x272822), // background
	ContrastBackgroundColor:     tcell.NewHexColor(0x3E3D32), // darker background
	MoreContrastBackgroundColor: tcell.NewHexColor(0x49483E), // even darker
	BorderColor:                 tcell.NewHexColor(0xF8F8F2), // foreground
	TitleColor:                  tcell.NewHexColor(0x66D9EF), // blue
	GraphicsColor:               tcell.NewHexColor(0xAE81FF), // purple
	PrimaryTextColor:            tcell.NewHexColor(0xF8F8F2), // foreground
	SecondaryTextColor:          tcell.NewHexColor(0xE6DB74), // yellow
	TertiaryTextColor:           tcell.NewHexColor(0xA6E22E), // green
	InverseTextColor:            tcell.NewHexColor(0x272822), // background
	ContrastSecondaryTextColor:  tcell.NewHexColor(0xFD971F), // orange
}
