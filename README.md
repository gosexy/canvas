## About

gocanvas is an image processing library based on ImageMagick's MagickWand, for the Go programming language.

## License

The [MIT License](https://raw.github.com/xiam/gocanvas/master/LICENSE).

## Installing

In order to compile and install, MagickWand's C header files are required.

### Debian

Debian has an old version of MagickWand, in order to install gocanvas we need to install the old version and then upgrade it.

Getting the old version of MagickWand along all its dependencies.

    $ sudo aptitude install libmagickwand-dev

Installing a newer version of ImageMagick over the old files.

    $ sudo su
    # cd /usr/local/src
    # wget http://www.imagemagick.org/download/ImageMagick.tar.gz 
    # tar xvzf ImageMagick.tar.gz
    # cd ImageMagick-6.x.y
    # ./configure --prefix=/usr
    # make
    # make install

### Arch Linux

Arch Linux already has a recent version of MagickWand.

    $ sudo pacman -S extra/imagemagick

### Windows

Choose your [favorite binary](http://imagemagick.com/script/binary-releases.php#windows).

### OSX

Installing using [MacPorts](http://www.macports.org/)

    $ sudo port install ImageMagick

### Other OS

Please, follow the [install from source](http://imagemagick.com/script/install-source.php?ImageMagick=9uv1bcgofrv21mhftmlk4v1465) tutorial.

### Pulling gocanvas from github

After installing ImageMagick's C header files, pull gocanvas from github:

    $ go get github.com/xiam/gocanvas/canvas

Then, the source will be in:

    $ find $GOPATH/src/github.com/xiam/gocanvas/canvas

Note that if you don't have the required C header files installed, gocanvas will fail to install:

    $ canvas.go:27:30: error: wand/magick_wand.h: No such file or directory

## Updating

After installing, you can use `go get -u github.com/xiam/gocanvas/canvas` to update canvas to the latest version.

## Usage example

Write an example.go file
 
    package main

    import "github.com/xiam/gocanvas/canvas"

    func main() {
      cv := canvas.New()
      defer cv.Destroy()

      // Opening some image from disk.
      opened := cv.Open("examples/input/example.png")

      if opened {

        // Photo auto orientation based on EXIF tags.
        canvas.AutoOrientate()

        // Creating a squared thumbnail
        canvas.Thumbnail(100, 100)

        // Saving the thumbnail to disk.
        canvas.Write("examples/output/example-thumbnail.png")

      }
    }

Then, run it with `go run example.go`

## Documentation

For full documentation run `go doc github.com/xiam/gocanvas/canvas`
