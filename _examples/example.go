package main

import (
	"github.com/gosexy/canvas"
	"fmt"
)

func main() {
	image := canvas.New()

	// Opening some image from disk.
	err := image.Open("input/example.png")

	if err == nil {

		defer image.Destroy()

		// Photo auto orientation based on EXIF tags.
		image.AutoOrientate()

		// Creating a squared thumbnail
		image.Thumbnail(100, 100)

		// Saving the thumbnail to disk.
		image.Write("output/example-thumbnail.png")
	} else {
		fmt.Printf("Could not open image: %s\n", err)
	}
}
