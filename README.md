## gosexy/canvas

`gosexy/canvas` is an image processing library based on ImageMagick's MagickWand, for the Go programming language. It uses cgo.

## Requeriments

ImageMagick's development files are required.

```sh
# OSX
$ brew install imagemagick

# Arch Linux
$ sudo pacman -S extra/imagemagick

# Debian
$ sudo aptitude install libmagickwand-dev
```

## Installation

Just pull `gosexy/canvas` from github using `go get`:

```sh
$ go get github.com/gosexy/canvas
```

## Usage

```go
package main

import "github.com/gosexy/canvas"

func main() {
  img := canvas.New()

  // Opening some image from disk.
  err := img.Open("examples/input/example.png")

  if err == nil {

    // Photo auto orientation based on EXIF tags.
    img.AutoOrientate()

    // Creating a squared thumbnail
    img.Thumbnail(100, 100)

    // Saving the thumbnail to disk.
    img.Write("examples/output/example-thumbnail.png")

  }
}
```

## Documentation

Read the `gosexy/canvas` documentation from a terminal

```go
$ go doc github.com/gosexy/canvas
```

Alternatively, you can [browse it](http://go.pkgdoc.org/github.com/gosexy/canvas) online.
