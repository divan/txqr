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

	testData     []byte
	animatingQR  []byte
	generatingQR bool
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
				vecty.If(!a.session.InProgress() && a.session.State() != StateFinished, elem.Div(
					a.settings,
				)),
				vecty.If(a.session.InProgress() || a.session.State() == StateFinished, elem.Div(
					a.resultsTable,
				)),
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
	if a.session.State() == StateNew {
		a.session.UpdateConfig(a.settings.Config())
	}

	setup, ok := a.session.StartNext()
	if !ok {
		vecty.Rerender(a)
		return
	}

	a.generatingQR = true
	vecty.Rerender(a)
	time.Sleep(100 * time.Millisecond) // wait till JS thread pickup rerender before running heavy computaiton. without this delay it'll never rerender

	log.Println("Creating animated gif for", setup)
	now := time.Now()
	gif, err := AnimatedGif(a.testData, 500, setup)
	if err != nil {
		log.Println("[ERROR] Can't generate gif: %v", err)
		// TODO: session abort
		a.generatingQR = false
		vecty.Rerender(a)
		return
	}
	log.Println("Took time:", time.Since(now))
	a.animatingQR = gif

	a.session.SetState(StateAnimating)
	a.generatingQR = false
	vecty.Rerender(a)
}

func newTestData() []byte {
	data := make([]byte, 10*1024)
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
	a.resultsTable.AddResult(res)
	vecty.Rerender(a)
}
