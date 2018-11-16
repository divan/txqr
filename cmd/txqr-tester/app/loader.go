package main

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
)

func loader() *vecty.HTML {
	return elem.Div(
		vecty.Markup(
			vecty.Class("title", "has-text-weight-light"),
		),
		elem.Heading1(
			vecty.Text("Generating..."),
		),
	)
}
