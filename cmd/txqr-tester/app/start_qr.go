package main

import (
	"fmt"

	"github.com/divan/txqr/qr"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
)

// StartQR renders the QR code with information needed to start
// a new testing from smartphone.
func (a *App) StartQR() vecty.ComponentOrHTML {
	img, err := qr.Encode(a.ws.address, 500, qr.Medium)
	if err != nil {
		// TODO(divan): display the error nicely (why this can even happen?)
		return elem.Div(vecty.Text(fmt.Sprintf("qr error: %v", err)))
	}
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
					vecty.Class("card-header-title", "is-centered"),
				),
				elem.Heading1(
					vecty.Markup(
						vecty.Class("has-text-weight-bold"),
					),
					vecty.If(!a.connected, vecty.Text("Scan QR code to connect")),
					vecty.If(a.connected, vecty.Text("Scan QR code to start testing")),
				),
			),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("card-image", "has-text-centered"),
			),
			renderImage(img),
		),

		elem.Footer(
			vecty.Markup(
				vecty.Class("card-footer"),
			),
			vecty.If(!a.connected,
				elem.Paragraph(
					vecty.Markup(
						vecty.Class("card-footer-item"),
					),
					vecty.Text(
						fmt.Sprintf("Started WS server on: %s", a.ws.address),
					),
				),
			),
			vecty.If(a.connected,
				elem.Paragraph(
					vecty.Markup(
						vecty.Class("card-footer-item", "has-background-success", "has-text-white", "has-text-weight-bold"),
					),
					vecty.Text(
						fmt.Sprintf("Connected"),
					),
				),
			),
		),
	)
}
