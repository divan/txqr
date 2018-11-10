package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/divan/txqr/protocol"
	"github.com/mdp/qrterminal"
	"github.com/pyk/byten"
)

func main() {
	size := flag.Int("size", 100, "Chunk size for data split per frame")
	delay := flag.Duration("delay", 100*time.Millisecond, "Delay between frames")
	flag.Parse()

	if len(flag.Args()) != 1 {
		log.Fatalf("Usage: %s file", os.Args[0])
	}

	filename := flag.Args()[0]

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Read input file failed: %v", err)
	}

	chunks, err := protocol.NewEncoder(*size).Encode(string(data))
	if err != nil {
		log.Fatalf("Encode failed: %v", err)
	}

	var avg time.Duration
	for {
		var total int
		start := time.Now()
		for _, chunk := range chunks {
			config := qrterminal.Config{
				Level:     qrterminal.M,
				Writer:    os.Stdout,
				BlackChar: qrterminal.WHITE,
				WhiteChar: qrterminal.BLACK,
				QuietZone: 1,
			}

			clearScreen()
			total += len(chunk)
			duration := time.Since(start)
			rate := int(time.Second) * total / int(duration)

			fmt.Printf("Speed: %v/s | whole file: %v\n", byten.Size(int64(rate)), avg)

			qrterminal.GenerateWithConfig(chunk, config)

			time.Sleep(*delay)
		}

		avg = time.Since(start)
	}
}

// TODO(divan): replace with crossplatform version
func clearScreen() {
	print("\033[H\033[2J")
}
