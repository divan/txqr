package main

import (
	"fmt"

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
		vecty.If(len(r.results) > 0,
			elem.Span(
				vecty.Markup(
					vecty.Class("is-pulled-right"),
				),
				r.csvButton(),
			),
		),
		elem.Heading1(
			vecty.Markup(
				vecty.Class("title", "has-text-weight-light"),
			),
			vecty.Text("Results"),
		),
		r.table(),
	)
}

// AddResult adds a new result and refreshes the table.
func (r *ResultsTable) AddResult(res Result) {
	r.results = append(r.results, res)
	vecty.Rerender(r)
}

func (r *ResultsTable) table() vecty.ComponentOrHTML {
	return elem.Table(
		vecty.Markup(
			vecty.Class("table", "is-fullwidth"),
		),
		r.thead(),
		r.tresults(),
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

func (r *ResultsTable) tresults() vecty.ComponentOrHTML {
	rows := make([]vecty.MarkupOrChild, len(r.results))
	for _, res := range r.results {
		rows = append(rows, tableRow(res))
	}
	return elem.TableBody(rows...)
}

func tableRow(r Result) *vecty.HTML {
	return elem.TableRow(
		elem.TableData(vecty.Text(fmt.Sprintf("%s", r.lvl))),
		elem.TableData(vecty.Text(fmt.Sprintf("%d", r.fps))),
		elem.TableData(vecty.Text(fmt.Sprintf("%d", r.size))),
		elem.TableData(vecty.Text(fmt.Sprintf("%s", r.Duration))),
	)
}

func (r *ResultsTable) csvButton() *vecty.HTML {
	return elem.Anchor(
		vecty.Markup(
			vecty.Class("button", "is-success"),
			vecty.Attribute("href", r.csvDataURI()),
			vecty.Attribute("download", "test_results.csv"),
		),
		vecty.Text("Download CSV"),
	)
}
