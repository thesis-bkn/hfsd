package utils

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"os"
	"strings"
)

// SaveBase64ImageToFile decodes a Base64 encoded image and saves it to the specified file path
func SaveBase64ImageToFile(base64String, filePath string) error {
	// Split the base64 string to get the format and the actual data
	parts := strings.Split(base64String, ",")
	if len(parts) != 2 {
		return fmt.Errorf("invalid base64 string format")
	}

	// Decode the base64 string to binary data
	imageData, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return fmt.Errorf("failed to decode base64 string: %w", err)
	}

	// Write the binary data to a file
	if err := ioutil.WriteFile(filePath, imageData, 0644); err != nil {
		return fmt.Errorf("failed to write image to file: %w", err)
	}

	return nil
}

// SaveRGBAImageToPNG saves an *image.RGBA to a PNG file
func SaveRGBAImageToJPEG(img *image.RGBA, filePath string) error {
	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Encode the RGBA image to PNG and write to file
	if err := jpeg.Encode(file, img, nil); err != nil {
		return fmt.Errorf("failed to encode image to PNG: %w", err)
	}

	return nil
}
