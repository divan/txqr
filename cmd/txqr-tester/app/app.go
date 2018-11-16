package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
)

// App is a top-level app component.
type App struct {
	vecty.Core

	session      *Session
	settings     *Settings
	resultsTable *ResultsTable

	ws *WSClient

	connected bool

	testData    []byte
	animatingQR []byte
}

// NewApp creates and inits new app page.
func NewApp() *App {
	wsAddress := js.Global.Get("WSAddress").String()
	fmt.Println("WSaddress:", wsAddress)
	app := &App{
		session:      NewSession(),
		settings:     NewSettings(),
		resultsTable: NewResultsTable(),
		testData:     newTestData(),
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
				elem.Div(a.QR()),
			),
			// Right half
			elem.Div(
				vecty.Markup(
					vecty.Class("column", "is-half"),
				),
				elem.Div(
					vecty.Markup(
						// TODO(divan): make it really disabled/greyed out
						vecty.MarkupIf(a.session.State() != StateNew,
							vecty.Attribute("disabled", true),
						),
					),
					a.settings,
				),
				elem.HorizontalRule(),
				a.resultsTable,
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

// SetConnected changes the connected status on UI.
func (a *App) SetConnected(val bool) {
	a.connected = val
	vecty.Rerender(a)
}

// ShowNext handles request to show next animated QR.
func (a *App) ShowNext() {
	setup, _ := a.session.StartNext()
	log.Println("Creating animated gif for", setup)
	now := time.Now()
	gif, err := AnimatedGif(a.testData, 500, setup)
	if err != nil {
		log.Println("[ERROR] Can't generate gif: %v", err)
		// TODO: session abort
		return
	}
	log.Println("Took time:", time.Since(now))
	a.animatingQR = gif
	a.session.SetState(StateAnimating)
	vecty.Rerender(a)
}

func newTestData() []byte {
	data := make([]byte, 1024)
	_, err := rand.Read(data)
	if err != nil {
		log.Println("[ERROR] Can't generate rand data: %v", err)
	}
	return data
}

// ProcessResult handles request to process new incoming result.
func (a *App) ProcessResult(res Result) {
	log.Println("Duration was:", res.Duration)

	a.session.SetState(StateWaitingNext)
	vecty.Rerender(a)
}
