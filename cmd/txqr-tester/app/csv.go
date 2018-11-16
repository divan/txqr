package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/url"
	"time"
)

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
		dur := fmt.Sprintf("%d", result.Duration/time.Millisecond)
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

func (r *ResultsTable) csvDataURI() string {
	var buf bytes.Buffer
	if err := r.WriteAsCSV(&buf); err != nil {
		log.Println("[ERROR] Can't create CSV:", err)
		return ""
	}
	prefix := "data:text/csv;charset=utf-8,"
	return prefix + url.QueryEscape(buf.String())
}
