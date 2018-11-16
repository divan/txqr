package main

import (
	"fmt"

	"github.com/divan/txqr/qr"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
)

// QR renders the QR code with accopmanying text.
func (a *App) QR() vecty.ComponentOrHTML {
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
			vecty.If(!a.connected, a.startQR()),
			vecty.If(a.connected, a.syncQR()),
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

func (a *App) startQR() vecty.ComponentOrHTML {
	return renderQR(a.ws.address)
}

func (a *App) syncQR() vecty.ComponentOrHTML {
	return renderQR("nextRound")
}

func renderQR(text string) vecty.ComponentOrHTML {
	img, err := qr.Encode(text, 500, qr.Medium)
	if err != nil {
		// TODO(divan): display the error nicely (why this can even happen?)
		return elem.Div(vecty.Text(fmt.Sprintf("qr error: %v", err)))
	}
	return renderImage(img)
}
