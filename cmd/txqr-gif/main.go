package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"image"
	"image/gif"
	"io/ioutil"
	"log"
	"os"

	"github.com/divan/txqr"
	"github.com/divan/txqr/qr"
)

func main() {
	splitSize := flag.Int("split", 100, "Chunk size for data split per frame")
	size := flag.Int("size", 300, "QR code size")
	fps := flag.Int("fps", 5, "Animation FPS")
	output := flag.String("o", "out.gif", "Output animated gif file")
	flag.Parse()

	if len(flag.Args()) != 1 {
		log.Fatalf("Usage: %s file", os.Args[0])
	}

	filename := flag.Args()[0]

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("[ERROR] Read input file: %v", err)
	}

	out, err := AnimatedGif(data, *size, *fps, *splitSize, qr.Medium)
	if err != nil {
		log.Fatalf("[ERROR] Creating animated gif: %v", err)
	}

	err = ioutil.WriteFile(*output, out, 0660)
	if err != nil {
		log.Fatalf("[ERROR] Create file: %v", err)
	}
	log.Println("Written output to", *output)
}

func AnimatedGif(data []byte, imgSize int, fps, size int, lvl qr.RecoveryLevel) ([]byte, error) {
	str := base64.StdEncoding.EncodeToString(data)
	chunks, err := txqr.NewEncoder(size).Encode(str)
	if err != nil {
		return nil, fmt.Errorf("encode: %v", err)
	}

	out := &gif.GIF{
		Image: make([]*image.Paletted, len(chunks)),
		Delay: make([]int, len(chunks)),
	}
	for i, chunk := range chunks {
		qr, err := qr.Encode(chunk, imgSize, lvl)
		if err != nil {
			return nil, fmt.Errorf("QR encode: %v", err)
		}
		out.Image[i] = qr.(*image.Paletted)
		out.Delay[i] = fpsToGifDelay(fps)
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
