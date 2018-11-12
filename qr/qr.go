package qr

import (
	"fmt"
	"image"

	"github.com/makiuchi-d/gozxing"
	zqrcode "github.com/makiuchi-d/gozxing/qrcode"
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
	bitmap, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", fmt.Errorf("gozxing bitmap: %v", err)
	}

	hints := make(map[gozxing.DecodeHintType]interface{})
	hints[gozxing.DecodeHintType_PURE_BARCODE] = true
	result, err := zqrcode.NewQRCodeReader().Decode(bitmap, hints)
	if err != nil {
		return "", fmt.Errorf("gozxing: %v", err)
	}
	return result.GetText(), nil
}
