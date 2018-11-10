package txqr

import (
	"fmt"
	"image"

	"github.com/divan/txqr/protocol"
	"github.com/divan/txqr/qr"
)

func Encode(str string) (<-chan image.Image, error) {
	size := 100
	chunks, err := protocol.NewEncoder(size).Encode(str)
	if err != nil {
		return nil, fmt.Errorf("encode: %v", err)
	}
	ch := make(chan image.Image)

	go func(ch chan<- image.Image) {
		defer close(ch)
		for _, chunk := range chunks {
			img, err := qr.Encode(chunk, 512)
			if err != nil {
				// TODO: handle error better
				fmt.Errorf("[ERROR] encode: %v", err)
				continue
			}
			ch <- img
		}
	}(ch)

	return ch, nil
}

func Decode(ch <-chan image.Image) (string, error) {
	dec := protocol.NewDecoder()
	for img := range ch {
		chunk, err := qr.Decode(img)
		if err != nil {
			return "", fmt.Errorf("decode: %v", err)
		}
		err = dec.DecodeChunk(chunk)
		if err != nil {
			return "", fmt.Errorf("decode chunk: %v", err)
		}
	}
	return dec.Data(), nil
}
