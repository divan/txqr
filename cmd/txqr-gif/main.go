package main

import (
	"encoding/base64"
	"flag"
	"image"
	"image/gif"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/divan/txqr/protocol"
	"github.com/divan/txqr/qr"
)

func main() {
	splitSize := flag.Int("split", 100, "Chunk size for data split per frame")
	size := flag.Int("size", 300, "QR code size")
	delay := flag.Duration("delay", 100*time.Millisecond, "Delay between frames")
	output := flag.String("o", "out.gif", "Output animated gif file")
	flag.Parse()

	if len(flag.Args()) != 1 {
		log.Fatalf("Usage: %s file", os.Args[0])
	}

	filename := flag.Args()[0]

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Read input file failed: %v", err)
	}

	str := base64.StdEncoding.EncodeToString(data)
	log.Println("Base64 encoded size:", len(str), "bytes")
	chunks, err := protocol.NewEncoder(*splitSize).Encode(str)
	if err != nil {
		log.Fatalf("Encode failed: %v", err)
	}

	out := &gif.GIF{
		Image: make([]*image.Paletted, len(chunks)),
		Delay: make([]int, len(chunks)),
	}
	for i, chunk := range chunks {
		qr, err := qr.Encode(chunk, *size)
		if err != nil {
			log.Fatalf("[ERROR] QR: %v", err)
		}
		out.Image[i] = qr.(*image.Paletted)
		out.Delay[i] = int(*delay / (10 * time.Millisecond)) // yeah, delays are in 100th of a second
	}

	fd, err := os.Create(*output)
	if err != nil {
		log.Fatalf("[ERROR] Create file: %v", err)
	}
	err = gif.EncodeAll(fd, out)
	if err != nil {
		log.Fatalf("[ERROR] Generate gif: %v", err)
	}
	log.Printf("Saved %d frames into %s\n", len(chunks), *output)

}
