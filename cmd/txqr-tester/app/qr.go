package main

import (
	"fmt"
	"log"

	"github.com/divan/txqr/qr"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
)

// QR renders the QR code with accopmanying text.
func (a *App) QR() vecty.ComponentOrHTML {
	state := a.session.State()
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
					vecty.If(state == StateFinished, vecty.Text("Test completed")),
					vecty.If(!a.connected, vecty.Text("Scan QR code to connect")),
					vecty.If(a.connected && !a.session.InProgress() && state != StateAnimating && state != StateFinished, vecty.Text("Scan QR code for next test")),
					vecty.If(a.connected && state == StateAnimating, vecty.Text("Sending data via QR...")),
				),
			),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("card-image", "has-text-centered"),
			),
			vecty.If(!a.connected, a.startQR()),
			vecty.If(a.connected, a.mainQR()),
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
			vecty.If(a.connected && state != StateFinished,
				elem.Paragraph(
					vecty.Markup(
						vecty.Class("card-footer-item", "has-background-success", "has-text-white", "has-text-weight-bold"),
					),
					vecty.Text(
						fmt.Sprintf("Connected"),
					),
				),
			),
			vecty.If(state == StateFinished,
				elem.Paragraph(
					vecty.Markup(
						vecty.Class("card-footer-item", "has-background-primary", "has-text-white", "has-text-weight-bold"),
					),
					vecty.Text(
						fmt.Sprintf("Finished"),
					),
				),
			),
		),
	)
}

func (a *App) startQR() vecty.ComponentOrHTML {
	return renderQR(a.ws.address)
}

func (a *App) mainQR() vecty.ComponentOrHTML {
	state := a.session.State()

	log.Println("MainQR: generatingQR", a.generatingQR)
	if a.generatingQR {
		return loader()
	}

	if state == StateAnimating {
		return renderGIF(a.animatingQR)
	} else if state == StateStarted || state == StateWaitingNext || state == StateNew {
		return renderQR("nextRound")
	}
	return elem.Div()
}

func renderQR(text string) vecty.ComponentOrHTML {
	img, err := qr.Encode(text, 500, qr.Medium)
	if err != nil {
		// TODO(divan): display the error nicely (why this can even happen?)
		return elem.Div(vecty.Text(fmt.Sprintf("qr error: %v", err)))
	}
	return renderPNG(img)
}
