package main

import (
	"fmt"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
)

// App is a top-level app component.
type App struct {
	vecty.Core

	session  *Session
	settings *Settings

	ws *WSClient

	connected bool
}

// NewApp creates and inits new app page.
func NewApp() *App {
	wsAddress := js.Global.Get("WSAddress").String()
	fmt.Println("WSaddress:", wsAddress)
	app := &App{
		session:  NewSession(),
		settings: NewSettings(),
	}

	app.ws = NewWSClient(wsAddress, app)

	go app.ws.talkToBackend()

	return app
}

// Render implements the vecty.Component interface.
func (a *App) Render() vecty.ComponentOrHTML {
	return elem.Body(
		a.header(),
		elem.Div(
			vecty.Markup(
				vecty.Class("columns"),
			),
			// Left half
			elem.Div(
				vecty.Markup(
					vecty.Class("column", "is-half"),
				),
				elem.Div(
					vecty.If(a.session.state == StateNew,
						a.StartQR()),
				),
			),
			// Right half
			elem.Div(
				vecty.Markup(
					vecty.Class("column", "is-half"),
				),
				elem.Div(
					a.settings,
				),
			),
		),
		vecty.Markup(
			event.KeyDown(a.KeyListener),
		),
	)
}

func (a *App) header() *vecty.HTML {
	return elem.Section(
		elem.Heading1(
			vecty.Markup(
				vecty.Class("title", "has-text-weight-light"),
			),
			vecty.Text("TXQR Automated Tester"),
		),
		elem.Heading6(
			vecty.Markup(
				vecty.Class("subtitle", "has-text-weight-light"),
			),
			vecty.Text("Run TQXR Reader app on your smartphone and point to the QR code to start testing."),
		),
	)
}

func (a *App) SetConnected(val bool) {
	a.connected = val
	vecty.Rerender(a)
}
