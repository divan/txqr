//go:generate gopherjs build ./app -o ./app/app.js
package main

import (
	"flag"
	"log"
)

func main() {
	flag.Parse()

	if err := StartApp(":1999"); err != nil {
		log.Fatalf("[ERROR] Can't start web server to serve the app")
	}
}
