## About

gocanvas is an image processing library based on ImageMagick's
MagickWand (http://www.imagemagick.org) for the Go programming language.

## License

The MIT license.

## Installing

Pull the last stable version from github:

    $ go get github.com/xiam/gocanvas/canvas

After pulling, the source will be in:

    $ cd $GOPATH/src/github.com/gocanvas/canvas

## Updating

You can use `go get -u -a` to update all installed packages.

## Usage example
 
    import "github.com/xiam/gocanvas/canvas"

    canvas := NewCanvas()

    opened := canvas.Open("examples/input/example.png")

    if opened {
      canvas.SetQuality(90)
      canvas.Write("examples/output/example.jpg")
    }

    canvas.Destroy()

## Documentation

For full documentation run `go doc canvas`
