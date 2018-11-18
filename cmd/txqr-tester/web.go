//go:generate go-bindata-assetfs app/index.tmpl app/app.js app/css/...
package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type PageInfo struct {
	WSAddress string
}

// indexTmpl is a html template for index page.
var (
	indexTmpl *template.Template
)

func init() {
	data, err := Asset("app/index.tmpl")
	if err != nil {
		panic(err)
	}
	indexTmpl = template.Must(template.New("index.tmpl").Parse(string(data)))
}

// StartApp generates app page, serves it via http
// and tries to open it using default browser.
func StartApp(bind string, noBrowser bool) error {
	ip := GetLocalIP()
	info := PageInfo{
		WSAddress: fmt.Sprintf("ws://%s%s/ws", ip, bind),
	}
	http.HandleFunc("/index.html", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, info)
	})
	http.Handle("/", redirectToIndex(http.FileServer(assetFS())))
	http.HandleFunc("/ws", NewWSBridge().Handle)

	if !noBrowser {
		go StartBrowser("http://localhost" + bind)
	}

	return http.ListenAndServe(bind, nil)
}

// handler handles index page.
func handler(w http.ResponseWriter, r *http.Request, info PageInfo) {
	err := indexTmpl.Execute(w, info)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(os.Stderr, "[ERROR] Failed to render template:", err)
		return
	}
}

// StartBrowser tries to open the URL in a browser
// and reports whether it succeeds.
//
// Orig. code: golang.org/x/tools/cmd/cover/html.go
func StartBrowser(url string) bool {
	// try to start the browser
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	fmt.Println("If browser window didn't appear, please go to this url:", url)
	return cmd.Start() == nil
}

func redirectToIndex(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(w, r, "/index.html", http.StatusMovedPermanently)
		}
		h.ServeHTTP(w, r)
	})
}
