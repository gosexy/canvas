## About

gocanvas is an image processing library based on ImageMagick's
MagickWand (http://www.imagemagick.org) for the Go programming language.

## License

The MIT license.

## Installing

Pull the last stable version from github:

    $ go get github.com/xiam/gocanvas/canvas

After pulling, the source will be in:

    $ find $GOPATH/src/github.com/xiam/gocanvas/canvas

## Updating

You can use `go get -u github.com/xiam/gocanvas/canvas` to update canvas to the latest version.

## Usage example

Write an example.go file
 
    package main

    import "github.com/xiam/gocanvas/canvas"

    func main() {
      cv := canvas.New()
      defer cv.Destroy()

      opened := cv.Open("examples/input/example.png")

      if opened {
        cv.SetQuality(90)
        cv.Write("examples/output/example.jpg")
      }
    }

And then run it with `go run example.go`

## Documentation

For full documentation run `go doc github.com/xiam/gocanvas/canvas`
