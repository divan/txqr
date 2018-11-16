package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"time"

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

// csv returns two-dimensional slice of prepared strings,
// to be used with encoding/csv package.
func (r *ResultsTable) csv() [][]string {
	ret := [][]string{}
	row := []string{"QR Level", "FPS", "Chunk Size", "Duration (ms)"}
	ret = append(ret, row)
	for _, result := range r.results {
		lvl := fmt.Sprintf("%s", result.lvl)
		fps := fmt.Sprintf("%d", result.fps)
		sz := fmt.Sprintf("%d", result.size)
		dur := fmt.Sprintf("%d", result.Duration*time.Millisecond)
		row := []string{lvl, fps, sz, dur}
		ret = append(ret, row)
	}
	return ret
}

// WriteAsCSV writes results into io.Writer in CSV format.
func (r *ResultsTable) WriteAsCSV(w io.Writer) error {
	cw := csv.NewWriter(w)
	for _, row := range r.csv() {
		err := cw.Write(row)
		if err != nil {
			return fmt.Errorf("write csv: %v", err)
		}
	}
	cw.Flush()
	if err := cw.Error(); err != nil {
		return fmt.Errorf("write csv: %v", err)
	}
	return nil
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
	return elem.Button(
		vecty.Markup(
			vecty.Class("button", "is-success"),
		),
		vecty.Text("Download CSV"),
	)
}
