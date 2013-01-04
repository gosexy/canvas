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
#cgo CFLAGS: -fopenmp -I./_include

#include <wand/magick_wand.h>

char *MagickGetPropertyName(char **properties, size_t index) {
  return properties[index];
}
*/
import "C"

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

// MagickWand constants
var (
	STROKE_BUTT_CAP   = uint(C.ButtCap)
	STROKE_ROUND_CAP  = uint(C.RoundCap)
	STROKE_SQUARE_CAP = uint(C.SquareCap)

	STROKE_MITER_JOIN = uint(C.MiterJoin)
	STROKE_ROUND_JOIN = uint(C.RoundJoin)
	STROKE_BEVEL_JOIN = uint(C.BevelJoin)

	FILL_EVEN_ODD_RULE = uint(C.EvenOddRule)
	FILL_NON_ZERO_RULE = uint(C.NonZeroRule)

	RAD_TO_DEG = 180 / math.Pi
	DEG_TO_RAD = math.Pi / 180

	UNDEFINED_ORIENTATION    = uint(C.UndefinedOrientation)
	TOP_LEFT_ORIENTATION     = uint(C.TopLeftOrientation)
	TOP_RIGHT_ORIENTATION    = uint(C.TopRightOrientation)
	BOTTOM_RIGHT_ORIENTATION = uint(C.BottomRightOrientation)
	BOTTOM_LEFT_ORIENTATION  = uint(C.BottomLeftOrientation)
	LEFT_TOP_ORIENTATION     = uint(C.LeftTopOrientation)
	RIGHT_TOP_ORIENTATION    = uint(C.RightTopOrientation)
	RIGHT_BOTTOM_ORIENTATION = uint(C.RightBottomOrientation)
	LEFT_BOTTOM_ORIENTATION  = uint(C.LeftBottomOrientation)
)

// Holds a Canvas object
type Canvas struct {
	wand *C.MagickWand

	fg *C.PixelWand
	bg *C.PixelWand

	drawing *C.DrawingWand

	fill   *C.PixelWand
	stroke *C.PixelWand

	filename string
	width    string
	height   string
}

func init() {
	C.MagickWandGenesis()
}

// Private: returns wand's hexadecimal color.
func getPixelHexColor(p *C.PixelWand) string {
	var rgb [3]float64

	rgb[0] = float64(C.PixelGetRed(p))
	rgb[1] = float64(C.PixelGetGreen(p))
	rgb[2] = float64(C.PixelGetBlue(p))

	return fmt.Sprintf("#%02x%02x%02x", int(rgb[0]*255.0), int(rgb[1]*255.0), int(rgb[2]*255.0))
}

// Private: returns MagickTrue or MagickFalse
func magickBoolean(value bool) C.MagickBooleanType {
	if value == true {
		return C.MagickTrue
	}
	return C.MagickFalse
}

// Opens an image file, returns nil on success, error otherwise.
func (self Canvas) Open(filename string) error {
	stat, err := os.Stat(filename)
	if err != nil {
		return err
	}
	if stat.IsDir() == true {
		return fmt.Errorf(`Could not open file "%s": it's a directory!`, filename)
	}
	status := C.MagickReadImage(self.wand, C.CString(filename))
	if status == C.MagickFalse {
		return fmt.Errorf(`Could not open image "%s": %s`, filename, self.Error())
	}
	self.filename = filename
	return nil
}

// Auto-orientates canvas based on its original image's EXIF metadata
func (self Canvas) AutoOrientate() error {

	data := self.Metadata()

	orientation, err := strconv.Atoi(data["exif:Orientation"])

	if err != nil {
		return err
	}

	switch uint(orientation) {
	case TOP_LEFT_ORIENTATION:
		// normal

	case TOP_RIGHT_ORIENTATION:
		self.Flop()

	case BOTTOM_RIGHT_ORIENTATION:
		self.RotateCanvas(math.Pi)

	case BOTTOM_LEFT_ORIENTATION:
		self.Flip()

	case LEFT_TOP_ORIENTATION:
		self.Flip()
		self.RotateCanvas(-math.Pi / 2)

	case RIGHT_TOP_ORIENTATION:
		self.RotateCanvas(-math.Pi / 2)

	case RIGHT_BOTTOM_ORIENTATION:
		self.Flop()
		self.RotateCanvas(-math.Pi / 2)

	case LEFT_BOTTOM_ORIENTATION:
		self.RotateCanvas(math.Pi / 2)

	default:
		return fmt.Errorf("No orientation data found in file.")
	}

	success := C.MagickSetImageOrientation(self.wand, (C.OrientationType)(TOP_LEFT_ORIENTATION))

	if success == C.MagickFalse {
		return fmt.Errorf("Could not orientate photo: %s", self.Error())
	}

	self.SetMetadata("exif:Orientation", (string)(TOP_LEFT_ORIENTATION))

	return nil
}

// Returns all metadata keys from the currently loaded image.
func (self Canvas) Metadata() map[string]string {
	var n C.ulong
	var i C.ulong

	var value *C.char
	var key *C.char

	data := make(map[string]string)

	properties := C.MagickGetImageProperties(self.wand, C.CString("*"), &n)

	for i = 0; i < n; i++ {
		key = C.MagickGetPropertyName(properties, C.size_t(i))
		value = C.MagickGetImageProperty(self.wand, key)

		data[strings.Trim(C.GoString(key), " ")] = strings.Trim(C.GoString(value), " ")

		C.MagickRelinquishMemory(unsafe.Pointer(value))
		C.MagickRelinquishMemory(unsafe.Pointer(key))
	}

	return data
}

// Returns the latest error reported by the MagickWand API.
func (self Canvas) Error() error {
	var t C.ExceptionType
	message := C.MagickGetException(self.wand, &t)
	C.MagickClearException(self.wand)
	return fmt.Errorf(C.GoString(message))
}

// Associates a metadata key with its value.
func (self Canvas) SetMetadata(key string, value string) error {
	success := C.MagickSetImageProperty(self.wand, C.CString(key), C.CString(value))
	if success == C.MagickFalse {
		return fmt.Errorf("Could not set metadata: %s", self.Error())
	}
	return nil
}

// Creates a horizontal mirror image by reflecting the pixels around the central y-axis.
func (self Canvas) Flop() error {
	success := C.MagickFlopImage(self.wand)
	if success == C.MagickFalse {
		return fmt.Errorf("Could not flop image: %s", self.Error())
	}
	return nil
}

// Creates a vertical mirror image by reflecting the pixels around the central x-axis.
func (self Canvas) Flip() error {
	success := C.MagickFlipImage(self.wand)
	if success == C.MagickFalse {
		return fmt.Errorf("Could not flop image: %s", self.Error())
	}
	return nil
}

// Clones an image to another canvas
func (self Canvas) Clone() *Canvas {
	clone := New()

	clone.SetBackgroundColor("none")

	clone.Blank(self.Width(), self.Height())

	clone.AppendCanvas(self, 0, 0)

	return clone
}

// Creates a centered thumbnail of the canvas.
func (self Canvas) Thumbnail(width uint, height uint) error {

	var ratio float64

	// Normalizing image.

	ratio = math.Min(float64(self.Width())/float64(width), float64(self.Height())/float64(height))

	if ratio < 1.0 {
		// Origin image is smaller than the thumbnail image.
		max := uint(math.Max(float64(width), float64(height)))

		// Empty replacement buffer with transparent background.
		replacement := New()

		replacement.SetBackgroundColor("none")

		replacement.Blank(max, max)

		// Putting original image in the center of the replacement canvas.
		replacement.AppendCanvas(self, int(int(width-self.Width())/2), int(int(height-self.Height())/2))

		// Replacing wand
		C.DestroyMagickWand(self.wand)

		self.wand = C.CloneMagickWand(replacement.wand)

	} else {
		// Is bigger, just resizing.
		err := self.Resize(uint(float64(self.Width())/ratio), uint(float64(self.Height())/ratio))
		if err != nil {
			return err
		}
	}

	// Now we have an image that we can use to crop the thumbnail from.
	err := self.Crop(int(int(self.Width()-width)/2), int(int(self.Height()-height)/2), width, height)

	if err != nil {
		return err
	}

	return nil
}

// Puts a canvas on top of the current one.
func (self Canvas) AppendCanvas(source Canvas, x int, y int) error {
	success := C.MagickCompositeImage(self.wand, source.wand, C.OverCompositeOp, C.long(x), C.long(y))

	if success == C.MagickFalse {
		return fmt.Errorf("Could not append image: %s", self.Error())
	}

	return nil
}

// Rotates the whole canvas.
func (self Canvas) RotateCanvas(rad float64) error {
	success := C.MagickRotateImage(self.wand, self.bg, C.double(RAD_TO_DEG*rad))

	if success == C.MagickFalse {
		return fmt.Errorf("Could not rotate image: %s", self.Error())
	}

	return nil
}

// Returns canvas' width.
func (self Canvas) Width() uint {
	return uint(C.MagickGetImageWidth(self.wand))
}

// Returns canvas' height.
func (self Canvas) Height() uint {
	return uint(C.MagickGetImageHeight(self.wand))
}

// Writes canvas to a file, returns true on success.
func (self Canvas) Write(filename string) error {
	err := self.Update()

	if err != nil {
		return err
	}

	success := C.MagickWriteImage(self.wand, C.CString(filename))

	if success == C.MagickFalse {
		return fmt.Errorf("Could not write: %s", self.Error())
	}

	return nil
}

// Changes the size of the canvas, returns true on success.
func (self Canvas) Resize(width uint, height uint) error {
	success := C.MagickResizeImage(self.wand, C.ulong(width), C.ulong(height), C.GaussianFilter, C.double(1.0))

	if success == C.MagickFalse {
		return fmt.Errorf("Could not resize: %s", self.Error())
	}

	return nil
}

// Adaptively changes the size of the canvas, returns true on success.
func (self Canvas) AdaptiveResize(width uint, height uint) error {
	success := C.MagickAdaptiveResizeImage(self.wand, C.ulong(width), C.ulong(height))

	if success == C.MagickFalse {
		return fmt.Errorf("Could not resize: %s", self.Error())
	}

	return nil
}

// Changes the compression quality of the canvas. Ranges from 1 (lowest) to 100 (highest).
func (self Canvas) SetQuality(quality uint) error {
	success := C.MagickSetImageCompressionQuality(self.wand, C.ulong(quality))

	if success == C.MagickFalse {
		return fmt.Errorf("Could not set compression quality: %s", self.Error())
	}

	return nil
}

// Returns the compression quality of the canvas. Ranges from 1 (lowest) to 100 (highest).
func (self Canvas) Quality() uint {
	return uint(C.MagickGetImageCompressionQuality(self.wand))
}

/*
// Sets canvas's foreground color.
func (self Canvas) SetColor(color string) (bool) {
  status := C.PixelSetColor(self.fg, C.CString(color))
  if status == C.MagickFalse {
    return false
  }
  return true
}
*/

// Sets canvas' background color.
func (self Canvas) SetBackgroundColor(color string) error {
	C.PixelSetColor(self.bg, C.CString(color))
	success := C.MagickSetImageBackgroundColor(self.wand, self.bg)

	if success == C.MagickFalse {
		return fmt.Errorf("Could not set background color: %s", self.Error())
	}

	return nil
}

// Returns canvas' background color.
func (self Canvas) BackgroundColor() string {
	return getPixelHexColor(self.bg)
}

// Sets antialiasing setting for the current drawing stroke.
func (self Canvas) SetStrokeAntialias(value bool) {
	C.DrawSetStrokeAntialias(self.drawing, magickBoolean(value))
}

// Returns antialiasing setting for the current drawing stroke.
func (self Canvas) StrokeAntialias() bool {
	value := C.DrawGetStrokeAntialias(self.drawing)
	if value == C.MagickTrue {
		return true
	}
	return false
}

// Sets the width of the stroke on the current drawing surface.
func (self Canvas) SetStrokeWidth(value float64) {
	C.DrawSetStrokeWidth(self.drawing, C.double(value))
}

// Returns the width of the stroke on the current drawing surface.
func (self Canvas) StrokeWidth() float64 {
	return float64(C.DrawGetStrokeWidth(self.drawing))
}

// Sets the opacity of the stroke on the current drawing surface.
func (self Canvas) SetStrokeOpacity(value float64) {
	C.DrawSetStrokeOpacity(self.drawing, C.double(value))
}

// Returns the opacity of the stroke on the current drawing surface.
func (self Canvas) StrokeOpacity() float64 {
	return float64(C.DrawGetStrokeOpacity(self.drawing))
}

// Sets the type of the line cap on the current drawing surface.
func (self Canvas) SetStrokeLineCap(value uint) {
	C.DrawSetStrokeLineCap(self.drawing, C.LineCap(value))
}

// Returns the type of the line cap on the current drawing surface.
func (self Canvas) StrokeLineCap() uint {
	return uint(C.DrawGetStrokeLineCap(self.drawing))
}

// Sets the type of the line join on the current drawing surface.
func (self Canvas) SetStrokeLineJoin(value uint) {
	C.DrawSetStrokeLineJoin(self.drawing, C.LineJoin(value))
}

// Returns the type of the line join on the current drawing surface.
func (self Canvas) StrokeLineJoin() uint {
	return uint(C.DrawGetStrokeLineJoin(self.drawing))
}

/*
func (self Canvas) SetFillRule(value int) {
  C.DrawSetFillRule(self.drawing, C.FillRule(value))
}
*/

// Sets the fill color for enclosed areas on the current drawing surface.
func (self Canvas) SetFillColor(color string) {
	C.PixelSetColor(self.fill, C.CString(color))
	C.DrawSetFillColor(self.drawing, self.fill)
}

// Returns the fill color for enclosed areas on the current drawing surface.
func (self Canvas) FillColor() string {
	return getPixelHexColor(self.fill)
}

// Sets the stroke color on the current drawing surface.
func (self Canvas) SetStrokeColor(color string) {
	C.PixelSetColor(self.stroke, C.CString(color))
	C.DrawSetStrokeColor(self.drawing, self.stroke)
}

// Returns the stroke color on the current drawing surface.
func (self Canvas) StrokeColor() string {
	return getPixelHexColor(self.stroke)
}

// Draws a circle over the current drawing surface.
func (self Canvas) Circle(radius float64) {
	C.DrawCircle(self.drawing, C.double(0), C.double(0), C.double(radius), C.double(0))
}

// Draws a rectangle over the current drawing surface.
func (self Canvas) Rectangle(x float64, y float64) {
	C.DrawRectangle(self.drawing, C.double(0), C.double(0), C.double(x), C.double(y))
}

// Moves the current coordinate system origin to the specified coordinate.
func (self Canvas) Translate(x float64, y float64) {
	C.DrawTranslate(self.drawing, C.double(x), C.double(y))
}

// Applies a scaling factor to the units of the current coordinate system.
func (self Canvas) Scale(x float64, y float64) {
	C.DrawScale(self.drawing, C.double(x), C.double(y))
}

// Draws a line starting on the current coordinate system origin and ending on the specified coordinates.
func (self Canvas) Line(x float64, y float64) {
	C.DrawLine(self.drawing, C.double(0), C.double(0), C.double(x), C.double(y))
}

/*
func (self Canvas) Skew(x float64, y float64) {
  C.DrawSkewX(self.drawing, C.double(x))
  C.DrawSkewY(self.drawing, C.double(y))
}
*/

// Applies a rotation of a given angle (in radians) on the current coordinate system.
func (self Canvas) Rotate(rad float64) {
	deg := RAD_TO_DEG * rad
	C.DrawRotate(self.drawing, C.double(deg))
}

// Draws an ellipse centered at the current coordinate system's origin.
func (self Canvas) Ellipse(a float64, b float64) {
	C.DrawEllipse(self.drawing, C.double(0), C.double(0), C.double(a), C.double(b), 0, 360)
}

// Clones the current drawing surface and stores it in a stack.
func (self Canvas) PushDrawing() error {
	success := C.PushDrawingWand(self.drawing)

	if success == C.MagickFalse {
		return fmt.Errorf("Could not push surface: %s", self.Error())
	}

	return nil
}

// Destroys the current drawing surface and returns the latest surface that was pushed to the stack.
func (self Canvas) PopDrawing() error {
	success := C.PopDrawingWand(self.drawing)

	if success == C.MagickFalse {
		return fmt.Errorf("Could not pop surface: %s", self.Error())
	}

	return nil
}

// Copies a drawing surface to the canvas.
func (self Canvas) Update() error {
	success := C.MagickDrawImage(self.wand, self.drawing)

	if success == C.MagickFalse {
		return fmt.Errorf("Could not update image: %s", self.Error())
	}

	return nil
}

// Destroys canvas.
func (self Canvas) Destroy() error {
	if self.wand != nil {
		C.DestroyMagickWand(self.wand)
		self.wand = nil
		return nil
	}
	return fmt.Errorf("Nothing to destroy")
	//C.MagickWandTerminus()
}

// Creates an empty canvas of the given dimensions.
func (self Canvas) Blank(width uint, height uint) error {
	success := C.MagickNewImage(self.wand, C.ulong(width), C.ulong(height), self.bg)

	if success == C.MagickFalse {
		return fmt.Errorf("Could not create image: %s", self.Error())
	}

	return nil
}

// Convolves the canvas with a Gaussian function given its standard deviation.
func (self Canvas) Blur(sigma float64) error {
	success := C.MagickBlurImage(self.wand, C.double(0), C.double(sigma))

	if success == C.MagickFalse {
		return fmt.Errorf("Could not blur image: %s", self.Error())
	}

	return nil
}

// Adaptively blurs the image by blurring less intensely near the edges and more intensely far from edges.
func (self Canvas) AdaptiveBlur(sigma float64) error {
	success := C.MagickAdaptiveBlurImage(self.wand, C.double(0), C.double(sigma))

	if success == C.MagickFalse {
		return fmt.Errorf("Could not blur image: %s", self.Error())
	}

	return nil
}

// Adds random noise to the canvas.
func (self Canvas) AddNoise() error {
	success := C.MagickAddNoiseImage(self.wand, C.GaussianNoise)

	if success == C.MagickFalse {
		return fmt.Errorf("Could not add noise: %s", self.Error())
	}

	return nil
}

// Removes a region of a canvas and collapses the canvas to occupy the removed portion.
func (self Canvas) Chop(x int, y int, width uint, height uint) error {
	success := C.MagickChopImage(self.wand, C.ulong(width), C.ulong(height), C.long(x), C.long(y))

	if success == C.MagickFalse {
		return fmt.Errorf("Could not chop: %s", self.Error())
	}

	return nil
}

// Extracts a region from the canvas.
func (self Canvas) Crop(x int, y int, width uint, height uint) error {
	success := C.MagickCropImage(self.wand, C.ulong(width), C.ulong(height), C.long(x), C.long(y))

	if success == C.MagickFalse {
		return fmt.Errorf("Could not crop: %s", self.Error())
	}

	return nil
}

// Adjusts the canvas's brightness given a factor (-1.0 thru 1.0)
func (self Canvas) SetBrightness(factor float64) error {

	factor = math.Max(-1, factor)
	factor = math.Min(1, factor)

	success := C.MagickModulateImage(self.wand, C.double(100+factor*100.0), C.double(100), C.double(100))

	if success == C.MagickFalse {
		return fmt.Errorf("Could not set brightness: %s", self.Error())
	}

	return nil
}

// Adjusts the canvas's saturation given a factor (-1.0 thru 1.0)
func (self Canvas) SetSaturation(factor float64) error {

	factor = math.Max(-1, factor)
	factor = math.Min(1, factor)

	success := C.MagickModulateImage(self.wand, C.double(100), C.double(100+factor*100.0), C.double(100))

	if success == C.MagickFalse {
		return fmt.Errorf("Could not set saturation: %s", self.Error())
	}

	return nil
}

// Adjusts the canvas's hue given a factor (-1.0 thru 1.0)
func (self Canvas) SetHue(factor float64) error {

	factor = math.Max(-1, factor)
	factor = math.Min(1, factor)

	success := C.MagickModulateImage(self.wand, C.double(100), C.double(100), C.double(100+factor*100.0))

	if success == C.MagickFalse {
		return fmt.Errorf("Could not set hue: %s", self.Error())
	}

	return nil
}

// Returns a new canvas object.
func New() *Canvas {
	self := &Canvas{}

	self.wand = C.NewMagickWand()

	self.fg = C.NewPixelWand()
	self.bg = C.NewPixelWand()

	self.fill = C.NewPixelWand()
	self.stroke = C.NewPixelWand()

	self.drawing = C.NewDrawingWand()

	//self.SetColor("#ffffff")
	self.SetBackgroundColor("none")

	self.SetStrokeColor("#ffffff")
	self.SetStrokeAntialias(true)
	self.SetStrokeWidth(1.0)
	self.SetStrokeOpacity(1.0)
	self.SetStrokeLineCap(STROKE_ROUND_CAP)
	self.SetStrokeLineJoin(STROKE_ROUND_JOIN)

	//self.SetFillRule(FILL_EVEN_ODD_RULE)
	self.SetFillColor("#888888")

	return self
}
