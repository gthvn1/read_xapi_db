package theme

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var GruvboxDark = tview.Theme{
	PrimitiveBackgroundColor:    tcell.NewHexColor(0x282828), // bg0_hard
	ContrastBackgroundColor:     tcell.NewHexColor(0x3c3836), // bg1
	MoreContrastBackgroundColor: tcell.NewHexColor(0x504945), // bg2

	BorderColor:   tcell.NewHexColor(0xebdbb2), // fg1
	TitleColor:    tcell.NewHexColor(0x83a598), // blue
	GraphicsColor: tcell.NewHexColor(0xd3869b), // purple

	PrimaryTextColor:   tcell.NewHexColor(0xebdbb2), // fg1
	SecondaryTextColor: tcell.NewHexColor(0xfabd2f), // yellow
	TertiaryTextColor:  tcell.NewHexColor(0xb8bb26), // green

	InverseTextColor:           tcell.NewHexColor(0x282828), // bg
	ContrastSecondaryTextColor: tcell.NewHexColor(0xfe8019), // orange
}
