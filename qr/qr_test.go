package qr

import (
	"image"
	"os"
	"testing"
)

func TestDecoder(t *testing.T) {
	filename := "testdata/helloworld_qr.png"
	fd, err := os.Open(filename)
	if err != nil {
		t.Fatalf("Open test file: %v", err)
	}
	img, _, err := image.Decode(fd)
	if err != nil {
		t.Fatalf("Decode test image: %v", err)
	}
	str, err := Decode(img)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}
	if str != "hello, world" {
		t.Fatalf("Expected 'hello, world', but got '%s'", str)
	}
}
