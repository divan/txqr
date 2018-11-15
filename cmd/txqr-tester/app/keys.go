package main

import (
	"fmt"

	"github.com/gopherjs/vecty"
)

// KeyListener implements listener for keydown events.
func (a *App) KeyListener(e *vecty.Event) {
	key := e.Get("key").String()
	switch key {
	case " ":
		fmt.Println("Space")
	}
}
