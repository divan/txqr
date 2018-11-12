package qr

import (
	"errors"
	"fmt"
	"image"

	"github.com/skip2/go-qrcode"
)

// Encode encodes data into the image with QR code.
func Encode(data string, size int) (image.Image, error) {
	qr, err := qrcode.New(data, qrcode.Medium)
	if err != nil {
		return nil, fmt.Errorf("encode QR: %v", err)
	}
	return qr.Image(size), nil
}

// Decode an image with QR code.
func Decode(img image.Image) (string, error) {
	return "", errors.New("not implemented")
}
