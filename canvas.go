package main

/*
#cgo LDFLAGS: -lMagickWand -lMagickCore 
#cgo CFLAGS: -fopenmp -I/usr/include/ImageMagick  
#include <stdlib.h>
#include <wand/magick_wand.h>
*/
import "C"

/*
import (
  "errors"
  "fmt"
)

type Errno int

func (e Errno) Error() string {
  s := errText[e]
  if s == "" {
    return fmt.Sprintf("errno %d", int(e))
  }
  return s
}
*/

type Canvas struct {
  wand *C.MagickWand
  filename string
  width string
  height string
}

func (c Canvas) init() {
  C.MagickWandGenesis()
}

func (c Canvas) Open(filename string) (bool) {
  status := C.MagickReadImage(c.wand, C.CString(filename))
  if status == C.MagickFalse {
    return false
  }
  return true
}

func (c Canvas) Write(filename string) (bool) {
  status := C.MagickWriteImage(c.wand, C.CString(filename))
  if status == C.MagickFalse {
    return false
  }
  return true
}

func (c Canvas) Destroy() {
  if c.wand != nil {
    C.DestroyMagickWand(c.wand)
  }
  C.MagickWandTerminus()
}

func NewCanvas() *Canvas {
  c := &Canvas{}
  c.init()
  c.wand = C.NewMagickWand()
  return c
}

func main() {

  canvas := NewCanvas()

  opened := canvas.Open("example.jpg")

  if opened {
    canvas.Write("example-go.png")
  }

  canvas.Destroy()

}
