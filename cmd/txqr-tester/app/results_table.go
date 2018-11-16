package main

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
)

// ResultsTable is a widget for configuring testing session results.
type ResultsTable struct {
	vecty.Core

	results []Result
}

// NewResultsTable creates and inits new results table widget.
func NewResultsTable() *ResultsTable {
	return &ResultsTable{}
}

// Render implements the vecty.Component interface for ResultsTable.
func (r *ResultsTable) Render() vecty.ComponentOrHTML {
	return elem.Div(
		elem.Heading1(
			vecty.Markup(
				vecty.Class("title", "has-text-weight-light"),
			),
			vecty.Text("Results"),
		),
		r.table(),
	)
}

func (r *ResultsTable) table() vecty.ComponentOrHTML {
	return elem.Table(
		vecty.Markup(
			vecty.Class("table", "is-fullwidth"),
		),
		r.thead(),
	)
}

func (r *ResultsTable) thead() vecty.ComponentOrHTML {
	return elem.TableHead(
		elem.TableRow(
			elem.TableHeader(vecty.Text("QR Lvl")),
			elem.TableHeader(vecty.Text("FPS")),
			elem.TableHeader(vecty.Text("Chunk size")),
			elem.TableHeader(vecty.Text("Result")),
		),
	)
}
