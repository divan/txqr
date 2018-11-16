package main

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
)

// ReadyScreen renders the widget with call to be ready to scan
// animated QR codes and point smartphone app to it.
func (a *App) ReadyScreen() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			vecty.Class("card"),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("card-header"),
			),
			elem.Div(
				vecty.Markup(
					vecty.Class("card-header-title", "is-centered", "is-success"),
				),
				elem.Heading1(
					vecty.Markup(
						vecty.Class("has-text-weight-bold"),
					),
					vecty.Text("Connected"),
				),
			),
		),
	)
}
