package main

import (
	"fmt"
	"strconv"

	"github.com/divan/txqr/qr"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
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
		elem.Heading1(
			vecty.Markup(
				vecty.Class("title", "has-text-weight-light"),
			),
			vecty.Text("Settings"),
		),
		s.chunkSizesRow(),
		s.fpsRow(),
		s.recoveryLevelsRow(),
		s.hint(),
	)
}

// Config returns current configuration.
func (s *Settings) Config() SessionConfig {
	return s.config
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
			s.sizeInput("from", &s.config.StartSize),
			s.sizeInput("to", &s.config.StopSize),
			label("Step"),
			s.numberInput("step", &s.config.SizeStep),
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
			s.numberInput("from", &s.config.StartFPS),
			s.numberInput("to", &s.config.StopFPS),
		),
	)
}

func (s *Settings) sizeInput(name string, val *int) vecty.ComponentOrHTML {
	str := strconv.Itoa(*val)
	stepStr := strconv.Itoa(s.config.SizeStep)
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
					vecty.Attribute("step", stepStr),

					event.Input(func(event *vecty.Event) {
						v := event.Target.Get("value").Int()

						*val = v
						vecty.Rerender(s)
					}),
				),
			),
		),
	)
}

func (s *Settings) numberInput(name string, val *int) vecty.ComponentOrHTML {
	str := strconv.Itoa(*val)
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
					vecty.Attribute("min", "0"),

					event.Input(func(event *vecty.Event) {
						v := event.Target.Get("value").Int()

						*val = v
						vecty.Rerender(s)
					}),
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
			s.checkboxInput("low", levels.has(qr.Low), func(v bool) { levels.set(qr.Low, v) }),
			s.checkboxInput("medium", levels.has(qr.Medium), func(v bool) { levels.set(qr.Medium, v) }),
			s.checkboxInput("high", levels.has(qr.High), func(v bool) { levels.set(qr.High, v) }),
			s.checkboxInput("highest", levels.has(qr.Highest), func(v bool) { levels.set(qr.Highest, v) }),
		),
	)
}

func (s *Settings) checkboxInput(name string, val bool, check func(bool)) vecty.ComponentOrHTML {
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

					event.Change(func(event *vecty.Event) {
						v := event.Target.Get("checked").Bool()
						fmt.Println("check", v)

						check(v)
						vecty.Rerender(s)
					}),
				),
			),
			vecty.Text(name),
		),
	)
}

func (s *Settings) hint() vecty.ComponentOrHTML {
	nChunks := (s.config.StopSize-s.config.StartSize)/s.config.SizeStep + 1
	nFPS := s.config.StopFPS - s.config.StartFPS + 1
	numberOfTests := s.config.Levels.numEnabled() * nFPS * nChunks
	text := fmt.Sprintf("This will run %d tests", numberOfTests)
	return elem.Div(
		vecty.Markup(
			vecty.Class("message", "is-info"),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("message-body"),
			),
			vecty.Text(text),
		),
	)
}
