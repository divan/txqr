package main

import (
	"strconv"

	"github.com/divan/txqr/qr"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/prop"
)

// Settings is a widget for configuring testing session settings.
type Settings struct {
	vecty.Core

	config SessionConfig
}

// NewSettings creates and inits new settings widget.
func NewSettings() *Settings {
	settings := &Settings{
		config: DefaultSessionConfig(),
	}

	return settings
}

// Render implements the vecty.Component interface for Settings.
func (s *Settings) Render() vecty.ComponentOrHTML {
	return elem.Div(
		s.chunkSizesRow(),
		s.fpsRow(),
		s.recoveryLevelsRow(),
	)
}

func (s *Settings) chunkSizesRow() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			vecty.Class("field", "is-horizontal"),
		),
		label("Chunk sizes"),
		elem.Div(
			vecty.Markup(
				vecty.Class("field-body"),
			),
			numberInput("from", s.config.StartSize),
			numberInput("to", s.config.StopSize),
			label("Step"),
			numberInput("step", s.config.SizeStep),
		),
	)
}

func (s *Settings) fpsRow() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			vecty.Class("field", "is-horizontal"),
		),
		label("FPS"),
		elem.Div(
			vecty.Markup(
				vecty.Class("field-body"),
			),
			numberInput("from", s.config.StartFPS),
			numberInput("to", s.config.StopFPS),
		),
	)
}

func numberInput(name string, val int) vecty.ComponentOrHTML {
	str := strconv.Itoa(val)
	return elem.Div(
		vecty.Markup(
			vecty.Class("field"),
		),
		elem.Paragraph(
			vecty.Markup(
				vecty.Class("control", "is-expanded"),
			),
			elem.Input(
				vecty.Markup(
					vecty.Class("input"),
					prop.Type(prop.TypeNumber),
					vecty.Attribute("placeholder", name),
					vecty.Attribute("value", str),
				),
			),
		),
	)
}

func label(name string) vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			vecty.Class("field-label", "is-normal"),
		),
		elem.Label(
			vecty.Markup(
				vecty.Class("label"),
			),
			vecty.Text(name+":"),
		),
	)
}

func (s *Settings) recoveryLevelsRow() vecty.ComponentOrHTML {
	levels := s.config.Levels
	return elem.Div(
		vecty.Markup(
			vecty.Class("field", "is-horizontal"),
		),
		label("Recovery levels"),
		elem.Div(
			vecty.Markup(
				vecty.Class("field-body"),
			),
			checkboxInput("low", levels.has(qr.Low)),
			checkboxInput("medium", levels.has(qr.Medium)),
			checkboxInput("high", levels.has(qr.High)),
			checkboxInput("highest", levels.has(qr.Highest)),
		),
	)
}

func checkboxInput(name string, val bool) vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			vecty.Class("control"),
		),
		elem.Label(
			vecty.Markup(
				vecty.Class("checkbox"),
				vecty.Style("padding", "7px"),
			),
			elem.Input(
				vecty.Markup(
					prop.Type(prop.TypeCheckbox),
					vecty.Style("margin", "7px"),
					vecty.Attribute("name", name),
					vecty.MarkupIf(val,
						vecty.Attribute("checked", "true"),
					),
				),
			),
			vecty.Text(name),
		),
	)
}
