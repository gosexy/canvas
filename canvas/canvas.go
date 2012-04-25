/*
  Copyright (c) 2012 José Carlos Nieto, http://xiam.menteslibres.org/

  Permission is hereby granted, free of charge, to any person obtaining
  a copy of this software and associated documentation files (the
  "Software"), to deal in the Software without restriction, including
  without limitation the rights to use, copy, modify, merge, publish,
  distribute, sublicense, and/or sell copies of the Software, and to
  permit persons to whom the Software is furnished to do so, subject to
  the following conditions:

  The above copyright notice and this permission notice shall be
  included in all copies or substantial portions of the Software.

  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
  EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
  MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
  NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
  LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
  OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
  WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package canvas

/*
#cgo LDFLAGS: -lMagickWand -lMagickCore 
#cgo CFLAGS: -fopenmp -I/usr/include/ImageMagick  
#include <stdlib.h>
#include <wand/magick_wand.h>
*/
import "C"

import "math"

import "fmt"

var (
  STROKE_BUTT_CAP       = C.ButtCap
  STROKE_ROUND_CAP      = C.RoundCap
  STROKE_SQUARE_CAP     = C.SquareCap
  
  STROKE_MITER_JOIN     = C.MiterJoin
  STROKE_ROUND_JOIN     = C.RoundJoin
  STROKE_BEVEL_JOIN     = C.BevelJoin

  FILL_EVEN_ODD_RULE    = C.EvenOddRule
  FILL_NON_ZERO_RULE    = C.NonZeroRule

  RAD_TO_DEG            = 180/math.Pi
  DEG_TO_RAD            = math.Pi/180
)

type Canvas struct {
  wand *C.MagickWand
  fg *C.PixelWand
  bg *C.PixelWand

  drawing *C.DrawingWand

  fill  *C.PixelWand
  stroke *C.PixelWand

  filename string
  width string
  height string

}

// Private: returns wand's hexadecimal color.
func getPixelHexColor(p *C.PixelWand) (string) {
  var rgb [3]float64

  rgb[0] = float64(C.PixelGetRed(p))
  rgb[1] = float64(C.PixelGetGreen(p))
  rgb[2] = float64(C.PixelGetBlue(p))

  return fmt.Sprintf("#%02x%02x%02x", int(rgb[0]*255.0), int(rgb[1]*255.0), int(rgb[2]*255.0))
}

// Private: returns MagickTrue or MagickFalse 
func magickBoolean(value bool) (C.MagickBooleanType) {
  if value == true {
    return C.MagickTrue
  }
  return C.MagickFalse
}

// Initializes the canvas environment.
func (cv Canvas) Init() {
  C.MagickWandGenesis()
}

// Opens an image file, returns true on success.
func (cv Canvas) Open(filename string) (bool) {
  status := C.MagickReadImage(cv.wand, C.CString(filename))
  if status == C.MagickFalse {
    return false
  }
  return true
}

// Returns canvas' width.
func (cv Canvas) Width() (int) {
  return int(C.MagickGetImageWidth(cv.wand))
}

// Returns canvas' height.
func (cv Canvas) Height() (int) {
  return int(C.MagickGetImageHeight(cv.wand))
}

// Writes canvas to a file, returns true on success.
func (cv Canvas) Write(filename string) (bool) {
  cv.Update()
  status := C.MagickWriteImage(cv.wand, C.CString(filename))
  if status == C.MagickFalse {
    return false
  }
  return true
}

// Changes the size of the canvas, returns true on success.
func (cv Canvas) Resize(width int, height int) (bool) {
  status := C.MagickResizeImage(cv.wand, C.size_t(width), C.size_t(height), C.GaussianFilter, C.double(1.0))
  if status == C.MagickFalse {
    return false
  }
  return true
}

// Adaptively changes the size of the canvas, returns true on success.
func (cv Canvas) AdaptiveResize(width int, height int) (bool) {
  status := C.MagickAdaptiveResizeImage(cv.wand, C.size_t(width), C.size_t(height))
  if status == C.MagickFalse {
    return false
  }
  return true
}

// Changes the compression quality of the canvas. Ranges from 1 (lowest) to 100 (highest).
func (cv Canvas) SetQuality(quality int) (bool) {
  status := C.MagickSetImageCompressionQuality(cv.wand, C.size_t(quality))
  if status == C.MagickFalse {
    return false
  }
  return true
}

// Returns the compression quality of the canvas. Ranges from 1 (lowest) to 100 (highest).
func (cv Canvas) Quality() (int) {
  return int(C.MagickGetImageCompressionQuality(cv.wand))
}

/*
// Sets canvas's foreground color.
func (cv Canvas) SetColor(color string) (bool) {
  status := C.PixelSetColor(cv.fg, C.CString(color))
  if status == C.MagickFalse {
    return false
  }
  return true
}
*/

// Sets canvas' background color.
func (cv Canvas) SetBackgroundColor(color string) (bool) {
  C.PixelSetColor(cv.bg, C.CString(color))
  status := C.MagickSetImageBackgroundColor(cv.wand, cv.bg)
  if status == C.MagickFalse {
    return false
  }
  return true
}

// Returns canvas' background color.
func (cv Canvas) BackgroundColor() (string) {
  return getPixelHexColor(cv.bg)
}

// Sets antialiasing setting for the current drawing stroke.
func (cv Canvas) SetStrokeAntialias(value bool) {
  C.DrawSetStrokeAntialias(cv.drawing, magickBoolean(value))
}

// Returns antialiasing setting for the current drawing stroke.
func (cv Canvas) StrokeAntialias() (bool) {
  value := C.DrawGetStrokeAntialias(cv.drawing)
  if (value == C.MagickTrue) {
    return true
  }
  return false
}

// Sets the width of the stroke on the current drawing surface.
func (cv Canvas) SetStrokeWidth(value float64) {
  C.DrawSetStrokeWidth(cv.drawing, C.double(value))
}

// Returns the width of the stroke on the current drawing surface.
func (cv Canvas) StrokeWidth() (float64) {
  return float64(C.DrawGetStrokeWidth(cv.drawing))
}

// Sets the opacity of the stroke on the current drawing surface.
func (cv Canvas) SetStrokeOpacity(value float64) {
  C.DrawSetStrokeOpacity(cv.drawing, C.double(value))
}

// Returns the opacity of the stroke on the current drawing surface.
func (cv Canvas) StrokeOpacity() (float64) {
  return float64(C.DrawGetStrokeOpacity(cv.drawing))
}

// Sets the type of the line cap on the current drawing surface.
func (cv Canvas) SetStrokeLineCap(value int) {
  C.DrawSetStrokeLineCap(cv.drawing, C.LineCap(value))
}

// Returns the type of the line cap on the current drawing surface.
func (cv Canvas) StrokeLineCap() (int) {
  return int(C.DrawGetStrokeLineCap(cv.drawing))
}

// Sets the type of the line join on the current drawing surface.
func (cv Canvas) SetStrokeLineJoin(value int) {
  C.DrawSetStrokeLineJoin(cv.drawing, C.LineJoin(value))
}

// Returns the type of the line join on the current drawing surface.
func (cv Canvas) StrokeLineJoin() (int) {
  return int(C.DrawGetStrokeLineJoin(cv.drawing))
}

/*
func (cv Canvas) SetFillRule(value int) {
  C.DrawSetFillRule(cv.drawing, C.FillRule(value))
}
*/

// Sets the fill color for enclosed areas on the current drawing surface.
func (cv Canvas) SetFillColor(color string) {
  C.PixelSetColor(cv.fill, C.CString(color))
  C.DrawSetFillColor(cv.drawing, cv.fill)
}

// Returns the fill color for enclosed areas on the current drawing surface.
func (cv Canvas) FillColor() (string) {
  return getPixelHexColor(cv.fill)
}

// Sets the stroke color on the current drawing surface.
func (cv Canvas) SetStrokeColor(color string) {
  C.PixelSetColor(cv.stroke, C.CString(color))
  C.DrawSetStrokeColor(cv.drawing, cv.stroke)
}

// Returns the stroke color on the current drawing surface.
func (cv Canvas) StrokeColor() (string) {
  return getPixelHexColor(cv.stroke)
}

// Draws a circle over the current drawing surface.
func (cv Canvas) Circle(radius float64) {
  C.DrawCircle(cv.drawing, C.double(0), C.double(0), C.double(radius), C.double(0))
}

// Draws a rectangle over the current drawing surface.
func (cv Canvas) Rectangle(x float64, y float64) {
  C.DrawRectangle(cv.drawing, C.double(0), C.double(0), C.double(x), C.double(y))
}

// Moves the current coordinate system origin to the specified coordinate.
func (cv Canvas) Translate(x float64, y float64) {
  C.DrawTranslate(cv.drawing, C.double(x), C.double(y))
}

// Applies a scaling factor to the units of the current coordinate system.
func (cv Canvas) Scale(x float64, y float64) {
  C.DrawScale(cv.drawing, C.double(x), C.double(y))
}

// Draws a line starting on the current coordinate system origin and ending on the specified coordinates.
func (cv Canvas) Line(x float64, y float64) {
  C.DrawLine(cv.drawing, C.double(0), C.double(0), C.double(x), C.double(y))
}

/*
func (cv Canvas) Skew(x float64, y float64) {
  C.DrawSkewX(cv.drawing, C.double(x))
  C.DrawSkewY(cv.drawing, C.double(y))
}
*/

// Applies a rotation of a given angle (in radians) on the current coordinate system.
func (cv Canvas) Rotate(rad float64) {
  deg := RAD_TO_DEG*rad
  C.DrawRotate(cv.drawing, C.double(deg))
}

// Draws an ellipse centered at the current coordinate system's origin.
func (cv Canvas) Ellipse(a float64, b float64) {
  C.DrawEllipse(cv.drawing, C.double(0), C.double(0), C.double(a), C.double(b), 0, 360)
}

// Clones the current drawing surface and stores it in a stack.
func (cv Canvas) PushDrawing() (bool) {
  status := C.PushDrawingWand(cv.drawing)
  if (status == C.MagickFalse) {
    return false
  }
  return true
}

// Destroys the current drawing surface and returns the latest surface that was pushed to the stack.
func (cv Canvas) PopDrawing() (bool) {
  status := C.PopDrawingWand(cv.drawing)
  if (status == C.MagickFalse) {
    return false
  }
  return true
}

// Copies a drawing surface to the canvas.
func (cv Canvas) Update() {
  C.MagickDrawImage(cv.wand, cv.drawing)
}

// Destroys canvas.
func (cv Canvas) Destroy() {
  if cv.wand != nil {
    C.DestroyMagickWand(cv.wand)
  }
  C.MagickWandTerminus()
}

// Creates an empty canvas of the given dimensions.
func (cv Canvas) Blank(width int, height int) (bool) {
  status := C.MagickNewImage(cv.wand, C.size_t(width), C.size_t(height), cv.bg)
  if status == C.MagickFalse {
    return false
  }
  return true
}

// Convolves the canvas with a Gaussian function given its standard deviation.
func (cv Canvas) Blur(sigma float64) (bool) {
  status := C.MagickBlurImage(cv.wand, C.double(0), C.double(sigma))
  if status == C.MagickFalse {
    return false
  }
  return true
}

// Adaptively blurs the image by blurring less intensely near the edges and more intensely far from edges.
func (cv Canvas) AdaptiveBlur(sigma float64) (bool) {
  status := C.MagickAdaptiveBlurImage(cv.wand, C.double(0), C.double(sigma))
  if status == C.MagickFalse {
    return false
  }
  return true
}

// Adds random noise to the canvas.
func (cv Canvas) AddNoise() (bool) {
  status := C.MagickAddNoiseImage(cv.wand, C.GaussianNoise)
  if status == C.MagickFalse {
    return false
  }
  return true
}

// Removes a region of a canvas and collapses the canvas to occupy the removed portion.
func (cv Canvas) Chop(x int, y int, width int, height int) (bool) {
  status := C.MagickChopImage(cv.wand, C.size_t(width), C.size_t(height), C.ssize_t(x), C.ssize_t(y))
  if status == C.MagickFalse {
    return false
  }
  return true
}

// Extracts a region from the canvas.
func (cv Canvas) Crop(x int, y int, width int, height int) (bool) {
  status := C.MagickCropImage(cv.wand, C.size_t(width), C.size_t(height), C.ssize_t(x), C.ssize_t(y))
  if status == C.MagickFalse {
    return false
  }
  return true
}

// Adjusts the canvas's brightness given a factor (-1.0 thru 1.0)
func (cv Canvas) SetBrightness(factor float64) (bool) {

  factor = math.Max(-1, factor)
  factor = math.Min(1, factor)

  status := C.MagickModulateImage(cv.wand, C.double(100 + factor*100.0), C.double(100), C.double(100))
  
  if status == C.MagickFalse {
    return false
  }

  return true
}

// Adjusts the canvas's saturation given a factor (-1.0 thru 1.0)
func (cv Canvas) SetSaturation(factor float64) (bool) {

  factor = math.Max(-1, factor)
  factor = math.Min(1, factor)

  status := C.MagickModulateImage(cv.wand, C.double(100), C.double(100 + factor*100.0), C.double(100))
  
  if status == C.MagickFalse {
    return false
  }

  return true
}

// Adjusts the canvas's hue given a factor (-1.0 thru 1.0)
func (cv Canvas) SetHue(factor float64) (bool) {

  factor = math.Max(-1, factor)
  factor = math.Min(1, factor)

  status := C.MagickModulateImage(cv.wand, C.double(100), C.double(100), C.double(100 + factor*100.0))
  
  if status == C.MagickFalse {
    return false
  }

  return true
}

// Returns a new canvas object.
func New() *Canvas {
  cv := &Canvas{}
  
  cv.Init()
  
  cv.wand = C.NewMagickWand()
  
  cv.fg = C.NewPixelWand()
  cv.bg = C.NewPixelWand()
  
  cv.fill   = C.NewPixelWand()
  cv.stroke = C.NewPixelWand()

  cv.drawing = C.NewDrawingWand()

  //cv.SetColor("#ffffff")
  cv.SetBackgroundColor("#000000")

  cv.SetStrokeColor("#ffffff")
  cv.SetStrokeAntialias(true)
  cv.SetStrokeWidth(1.0)
  cv.SetStrokeOpacity(1.0)
  cv.SetStrokeLineCap(STROKE_ROUND_CAP)
  cv.SetStrokeLineJoin(STROKE_ROUND_JOIN)

  //cv.SetFillRule(FILL_EVEN_ODD_RULE)
  cv.SetFillColor("#888888")
  
  return cv
}

