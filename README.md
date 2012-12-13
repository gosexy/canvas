## gosexy/canvas

``gosexy/canvas`` is an image processing library based on ImageMagick's MagickWand, for the Go programming language.

## Requeriments

### Mac OSX

The ImageMagick's header files are required. If you're using ``brew`` the installation is straightforward.

```sh
$ brew install imagemagick
```

### Debian

Debian has an old version of MagickWand (6.6.x), this binding was built against 6.7.x. Please check out the
[squeeze branch](https://github.com/gosexy/canvas/tree/squeeze) to get a version that works on Debian Squeeze and
probably other debian-based distros. This may not be required for Ubuntu.

### Arch Linux

Arch Linux already has a recent version of MagickWand.

```sh
$ sudo pacman -S extra/imagemagick
```

### Other Operative Systems

Please, follow the [install from source](http://imagemagick.com/script/install-source.php?ImageMagick=9uv1bcgofrv21mhftmlk4v1465) tutorial.

## Installation

After installing ImageMagick's header files, pull gocanvas from github:

```sh
$ go get github.com/gosexy/canvas
```

## Updating

After installing, you can use `go get -u github.com/gosexy/canvas` to keep up to date.

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

You can read ``gosexy/canvas`` documentation from a terminal

```go
$ go doc github.com/gosexy/canvas
```

Or you can [browse it](http://go.pkgdoc.org/github.com/gosexy/canvas) online.
