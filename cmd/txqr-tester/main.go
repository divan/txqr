//go:generate gopherjs build ./app -o ./app/app.js
package main

import (
	"flag"
	"log"
)

func main() {
	var noBrowser = flag.Bool("n", false, "Don't start browser automatically")
	flag.Parse()

	if err := StartApp(":1999", *noBrowser); err != nil {
		log.Fatalf("[ERROR] Can't start web server to serve the app")
	}
}
