package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/prop"
)

// renderPNG is a helper for converting image.Image into vecty-compatible
// component displaying this image as PNG.
func renderPNG(img image.Image) vecty.ComponentOrHTML {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		// TODO(divan): display the error nicely (why this can even happen?)
		return elem.Div(vecty.Text(fmt.Sprintf("png error: %v", err)))
	}
	src := base64.StdEncoding.EncodeToString(buf.Bytes())
	src = "data:image/png;base64," + src // prepare to be used as data object in src property
	return elem.Image(
		vecty.Markup(
			prop.Src(src),
		),
	)
}

// renderGIF is a helper for converting animated gif raw bytes into vecty-compatible
// component displaying this image as aGIF.
func renderGIF(data []byte) vecty.ComponentOrHTML {
	src := base64.StdEncoding.EncodeToString(data)
	src = "data:image/gif;base64," + src
	return elem.Image(
		vecty.Markup(
			prop.Src(src),
		),
	)
}
