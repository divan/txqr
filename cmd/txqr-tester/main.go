package main

import (
	"github.com/gopherjs/vecty"
)

func main() {
	app := NewApp()

	vecty.SetTitle("TXQR Automated Tester")
	vecty.AddStylesheet("css/bulma.css")
	vecty.AddStylesheet("css/bulma-extensions.min.css")
	vecty.AddStylesheet("css/custom.css")
	vecty.RenderBody(app)
}
