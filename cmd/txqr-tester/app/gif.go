package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/gif"

	"github.com/divan/txqr"
	"github.com/divan/txqr/qr"
)

func AnimatedGif(data []byte, imgSize int, setup *testSetup) ([]byte, error) {
	str := base64.StdEncoding.EncodeToString(data)
	fmt.Printf("Generating Animated QR Gif for %d bytes\n", len(str))
	chunks, err := txqr.NewEncoder(setup.size).Encode(str)
	if err != nil {
		return nil, fmt.Errorf("encode: %v", err)
	}

	out := &gif.GIF{
		Image: make([]*image.Paletted, len(chunks)),
		Delay: make([]int, len(chunks)),
	}
	for i, chunk := range chunks {
		qr, err := qr.Encode(chunk, imgSize, setup.lvl)
		if err != nil {
			return nil, fmt.Errorf("QR encode: %v", err)
		}
		out.Image[i] = qr.(*image.Paletted)
		out.Delay[i] = fpsToGifDelay(setup.fps)
	}

	var buf bytes.Buffer
	err = gif.EncodeAll(&buf, out)
	if err != nil {
		return nil, fmt.Errorf("gif create: %v", err)
	}
	return buf.Bytes(), nil
}

// fpsToGifDelay converts fps value into animated GIF
// delay value, which is in 100th of second
func fpsToGifDelay(fps int) int {
	if fps == 0 {
		return 100 // default value, 1 sec
	}
	return 100 / fps
}
