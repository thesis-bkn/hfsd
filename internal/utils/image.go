package utils

import (
	"bytes"
	"image"
	"image/png"
	"os"

	"github.com/ztrue/tracerr"
)

// loadPNGFromBytes decodes a PNG image from a byte slice
func loadPNGFromBytes(data []byte) (image.Image, error) {
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return img, nil
}

// savePNG saves a PNG image to the specified file path
func SavePNG(filePath string, data []byte) error {
	// Create a file
	outFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Encode the image to PNG and write it to the file
	img, err := loadPNGFromBytes(data)
	if err != nil {
		return tracerr.Wrap(err)
	}

	err = png.Encode(outFile, img)
	if err != nil {
		return err
	}

	return nil
}
