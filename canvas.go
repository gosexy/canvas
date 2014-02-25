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
#cgo CFLAGS: -fopenmp -I./_include
#cgo LDFLAGS: -lMagickWand -lMagickCore

#include <wand/magick_wand.h>

char *MagickGetPropertyName(char **properties, size_t index) {
  return properties[index];
}
*/
import "C"

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

// MagickWand constants
const (
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

	POINT_FILTER          = uint(C.PointFilter)
	BOX_FILTER            = uint(C.BoxFilter)
	TRIANGLE_FILTER       = uint(C.TriangleFilter)
	HERMITE_FILTER        = uint(C.HermiteFilter)
	HANNING_FILTER        = uint(C.HanningFilter)
	HAMMING_FILTER        = uint(C.HammingFilter)
	BLACKMAN_FILTER       = uint(C.BlackmanFilter)
	GAUSSIAN_FILTER       = uint(C.GaussianFilter)
	QUADRATIC_FILTER      = uint(C.QuadraticFilter)
	CUBIC_FILTER          = uint(C.CubicFilter)
	CATROM_FILTER         = uint(C.CatromFilter)
	MITCHEL_FILTER        = uint(C.MitchellFilter)
	BESSEL_FILTER         = uint(C.BesselFilter)
	JINC_FILTER           = uint(C.BesselFilter)
	SINC_FAST_FILTER      = uint(C.SincFastFilter)
	SINC_FILTER           = uint(C.SincFilter)
	KAISER_FILTER         = uint(C.KaiserFilter)
	WELSH_FILTER          = uint(C.WelshFilter)
	PARZEN_FILTER         = uint(C.ParzenFilter)
	BOHMAN_FILTER         = uint(C.BohmanFilter)
	BARTLETT_FILTER       = uint(C.BartlettFilter)
	LAGRANGE_FILTER       = uint(C.LagrangeFilter)
	LANCZOS_FILTER        = uint(C.LanczosFilter)
	LANCZOS_SHARP_FILTER  = uint(C.LanczosSharpFilter)
	LANCZOS2_FILTER       = uint(C.Lanczos2Filter)
	LANCZOS2_SHARP_FILTER = uint(C.Lanczos2SharpFilter)
	ROBIDOUX_FILTER       = uint(C.RobidouxFilter)

	UNDEFINED_TYPE              = uint(C.UndefinedType)
	BILEVEL_TYPE                = uint(C.BilevelType)
	GRAYSCALE_TYPE              = uint(C.GrayscaleType)
	GRAYSCALE_MATTER_TYPE       = uint(C.GrayscaleMatteType)
	PALETTE_TYPE                = uint(C.PaletteType)
	PALETTE_MATTE_TYPE          = uint(C.PaletteMatteType)
	TRUE_COLOR_TYPE             = uint(C.TrueColorType)
	TRUE_COLOR_MATTE_TYPE       = uint(C.TrueColorMatteType)
	COLOR_SEPARATION_TYPE       = uint(C.ColorSeparationType)
	COLOR_SEPARATION_MATTE_TYPE = uint(C.ColorSeparationMatteType)
	OPTIMIZE_TYPE               = uint(C.OptimizeType)

	UNDEFINED_INTERLACE = uint(C.UndefinedInterlace)
	NO_INTERLACE        = uint(C.NoInterlace)
	LINE_INTERLACE      = uint(C.LineInterlace)
	PLANE_INTERLACE     = uint(C.PlaneInterlace)
	PARTITION_INTERLACE = uint(C.PartitionInterlace)
	GIF_INTERLACE       = uint(C.GIFInterlace)
	JPEG_INTERLACE      = uint(C.JPEGInterlace)
	PNG_INTERLACE       = uint(C.PNGInterlace)

	UNDEFINED_GRAVITY  = uint(C.UndefinedGravity)
	FORGET_GRAVITY     = uint(C.ForgetGravity)
	NORTH_WEST_GRAVITY = uint(C.NorthWestGravity)
	NORTH_GRAVITY      = uint(C.NorthGravity)
	NORTH_EAST_GRAVITY = uint(C.NorthEastGravity)
	WEST_GRAVITY       = uint(C.WestGravity)
	CENTER_GRAVITY     = uint(C.CenterGravity)
	EAST_GRAVITY       = uint(C.EastGravity)
	SOUTH_WEST_GRAVITY = uint(C.SouthWestGravity)
	SOUTH_GRAVITY      = uint(C.SouthGravity)
	SOUTH_EAST_GRAVITY = uint(C.SouthEastGravity)
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

	quantumRange uint
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
func (self *Canvas) Open(filename string) error {
	stat, err := os.Stat(filename)

	if err != nil {
		return err
	}

	if stat.IsDir() == true {
		return fmt.Errorf(`Could not open file "%s": it's a directory!`, filename)
	}

	cfilename := C.CString(filename)
	status := C.MagickReadImage(self.wand, cfilename)
	C.free(unsafe.Pointer(cfilename))

	if status == C.MagickFalse {
		return fmt.Errorf(`Could not open image "%s": %s`, filename, self.Error())
	}

	self.filename = filename

	return nil
}

func (self *Canvas) SetOption(key, value string) error {
	ckey := C.CString(key)
	cvalue := C.CString(value)

	defer C.free(unsafe.Pointer(ckey))
	defer C.free(unsafe.Pointer(cvalue))

	if C.MagickSetOption(self.wand, ckey, cvalue) == C.MagickFalse {
		return fmt.Errorf(`Could not set option "%s" to "%s": %s`, ckey, cvalue, self.Error())
	}

	return nil
}

func (self *Canvas) SetCaption(content string) error {
	ccontent := C.CString("caption:" + content)
	defer C.free(unsafe.Pointer(ccontent))

	if C.MagickReadImage(self.wand, ccontent) == C.MagickFalse {
		return fmt.Errorf(`Could not open image "%s": %s`, content, self.Error())
	}

	C.MagickDrawImage(self.wand, self.drawing)

	return nil
}

func (self *Canvas) DrawAnnotation(content string, width, height uint) error {
	ccontent := C.CString("caption:" + content)
	defer C.free(unsafe.Pointer(ccontent))

	C.DrawAnnotation(self.drawing, 20, 20, (*C.uchar)(unsafe.Pointer(ccontent)))

	if C.MagickDrawImage(self.wand, self.drawing) == C.MagickFalse {
		return fmt.Errorf(`Could not draw annotation: %s`, self.Error())
	}

	return nil
}

// Reads an image or image sequence from a blob.
func (self *Canvas) OpenBlob(blob []byte, length uint) error {
	status := C.MagickReadImageBlob(self.wand, unsafe.Pointer(&blob[0]), C.size_t(length))

	if status == C.MagickFalse {
		return fmt.Errorf(`Could not open image from blob: %s`, self.Error())
	}

	return nil
}

// Auto-orientates canvas based on its original image's EXIF metadata
func (self *Canvas) AutoOrientate() error {

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
		self.RotateCanvas(math.Pi / 2)

	case RIGHT_BOTTOM_ORIENTATION:
		self.Flop()
		self.RotateCanvas(-math.Pi / 2)

	case LEFT_BOTTOM_ORIENTATION:
		self.RotateCanvas(-math.Pi / 2)

	default:
		return errors.New("No orientation data found in file.")
	}

	success := C.MagickSetImageOrientation(self.wand, (C.OrientationType)(TOP_LEFT_ORIENTATION))

	if success == C.MagickFalse {
		return fmt.Errorf("Could not orientate photo: %s", self.Error())
	}

	self.SetMetadata("exif:Orientation", (string)(TOP_LEFT_ORIENTATION))

	return nil
}

// Returns all metadata keys from the currently loaded image.
func (self *Canvas) Metadata() map[string]string {
	var n C.size_t
	var i C.size_t

	var value *C.char
	var key *C.char

	data := make(map[string]string)

	cplist := C.CString("*")

	properties := C.MagickGetImageProperties(self.wand, cplist, &n)

	C.free(unsafe.Pointer(cplist))

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
func (self *Canvas) Error() error {
	var t C.ExceptionType
	ptr := C.MagickGetException(self.wand, &t)
	message := C.GoString(ptr)
	C.MagickClearException(self.wand)
	C.MagickRelinquishMemory(unsafe.Pointer(ptr))
	return errors.New(message)
}

// Associates a metadata key with its value.
func (self *Canvas) SetMetadata(key string, value string) error {
	ckey := C.CString(key)
	cval := C.CString(value)

	success := C.MagickSetImageProperty(self.wand, ckey, cval)

	C.free(unsafe.Pointer(ckey))
	C.free(unsafe.Pointer(cval))

	if success == C.MagickFalse {
		return fmt.Errorf("Could not set metadata: %s", self.Error())
	}

	return nil
}

// Creates a horizontal mirror image by reflecting the pixels around the central y-axis.
func (self *Canvas) Flop() error {
	success := C.MagickFlopImage(self.wand)

	if success == C.MagickFalse {
		return fmt.Errorf("Could not flop image: %s", self.Error())
	}

	return nil
}

// Creates a vertical mirror image by reflecting the pixels around the central x-axis.
func (self *Canvas) Flip() error {
	success := C.MagickFlipImage(self.wand)

	if success == C.MagickFalse {
		return fmt.Errorf("Could not flop image: %s", self.Error())
	}

	return nil
}

// Adjusts the contrast of an image with a non-linear sigmoidal contrast algorithm. Increase the contrast of the image using a sigmoidal transfer function without saturating highlights or shadows. Contrast indicates how much to increase the contrast (0 is none; 3 is typical; 20 is pushing it); mid-point indicates where midtones fall in the resultant image (0 is white; 50 is middle-gray; 100 is black). Set sharpen to true to increase the image contrast otherwise the contrast is reduced.
func (self *Canvas) SigmoidalContrast(sharpen bool, alpha float64, beta float64) error {
	status := C.MagickSigmoidalContrastImage(self.wand, magickBoolean(sharpen), C.double(alpha), C.double(beta))

	if status == C.MagickFalse {
		return fmt.Errorf("Could not contrast image: %s", self.Error())
	}

	return nil
}

// Enhances the intensity differences between the lighter and darker elements of the image. Set sharpen to a value other than 0 to increase the image contrast otherwise the contrast is reduced.
func (self *Canvas) Contrast(sharpen bool) error {
	status := C.MagickContrastImage(self.wand, magickBoolean(sharpen))

	if status == C.MagickFalse {
		return fmt.Errorf("Could not contrast image: %s", self.Error())
	}

	return nil
}

// Clones an image to another canvas
func (self *Canvas) Clone() *Canvas {
	clone := New()

	clone.SetBackgroundColor("none")

	clone.Blank(self.Width(), self.Height())

	clone.AppendCanvas(self, 0, 0)

	return clone
}

// Converts the current image into a thumbnail of the specified
// width and height preserving ratio. It uses Crop() to clip the
// image to the specified area.
//
// If width or height are bigger than the current image, a centered
// thumbnail will be produced.
//
// Is width and height are smaller than the current image, the image
// will be resized and cropped, if needed.
func (self *Canvas) Thumbnail(width uint, height uint) error {

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
func (self *Canvas) AppendCanvas(source *Canvas, x int, y int) error {
	success := C.MagickCompositeImage(self.wand, source.wand, C.OverCompositeOp, C.ssize_t(x), C.ssize_t(y))

	if success == C.MagickFalse {
		return fmt.Errorf("Could not append image: %s", self.Error())
	}

	return nil
}

// Rotates the whole canvas.
func (self *Canvas) RotateCanvas(rad float64) error {
	success := C.MagickRotateImage(self.wand, self.bg, C.double(RAD_TO_DEG*rad))

	if success == C.MagickFalse {
		return fmt.Errorf("Could not rotate image: %s", self.Error())
	}

	return nil
}

// Returns canvas' width.
func (self *Canvas) Width() uint {
	return uint(C.MagickGetImageWidth(self.wand))
}

// Returns canvas' height.
func (self *Canvas) Height() uint {
	return uint(C.MagickGetImageHeight(self.wand))
}

// Writes canvas to a file, returns true on success.
func (self *Canvas) Write(filename string) error {
	err := self.Update()

	if err != nil {
		return err
	}

	cfilename := C.CString(filename)
	success := C.MagickWriteImage(self.wand, cfilename)
	C.free(unsafe.Pointer(cfilename))

	if success == C.MagickFalse {
		return fmt.Errorf("Could not write: %s", self.Error())
	}

	return nil
}

// Changes the size of the canvas, returns true on success.
func (self *Canvas) Resize(width uint, height uint) error {
	success := C.MagickResizeImage(self.wand, C.size_t(width), C.size_t(height), C.GaussianFilter, C.double(1.0))

	if success == C.MagickFalse {
		return fmt.Errorf("Could not resize: %s", self.Error())
	}

	return nil
}

// Changes the size of the canvas using specified filter and blur, returns true on success.
func (self *Canvas) ResizeWithFilter(width uint, height uint, filter uint, blur float32) error {
	if width == 0 && height == 0 {
		return errors.New("Please specify at least one of dimensions")
	}

	if width == 0 || height == 0 {
		origHeight := uint(C.MagickGetImageHeight(self.wand))
		origWidth := uint(C.MagickGetImageWidth(self.wand))

		if width == 0 {
			ratio := float32(origHeight) / float32(height)
			width = uint(float32(origWidth) / ratio)
		}
		if height == 0 {
			ratio := float32(origWidth) / float32(width)
			height = uint(float32(origHeight) / ratio)
		}
	}

	success := C.MagickResizeImage(self.wand, C.size_t(width), C.size_t(height), C.FilterTypes(filter), C.double(blur))

	if success == C.MagickFalse {
		return fmt.Errorf("Could not resize: %s", self.Error())
	}

	return nil
}

// Sharpens an image. We convolve the image with a Gaussian operator of the
// given radius and standard deviation (sigma). For reasonable results, the
// radius should be larger than sigma.
// Use a radius of 0 and selects a suitable radius for you.
// You can pass 0 as channel number - to use default channels
func (self *Canvas) SharpenImage(radius float32, sigma float32, channel int) error {
	if channel == 0 {
		channel = C.DefaultChannels
	}

	success := C.MagickSharpenImageChannel(self.wand, C.ChannelType(channel), C.double(radius), C.double(sigma))
	if success == C.MagickFalse {
		return fmt.Errorf("Could not resize: %s", self.Error())
	}

	return nil
}

func (self *Canvas) GetImageBlob() ([]byte, error) {
	return self.Blob()
}

func (self *Canvas) Blob() ([]byte, error) {
	var size C.size_t = 0

	p := unsafe.Pointer(C.MagickGetImageBlob(self.wand, &size))

	if size == 0 {
		return nil, errors.New("Could not get image blob.")
	}

	blob := C.GoBytes(p, C.int(size))

	C.MagickRelinquishMemory(p)

	return blob, nil
}

// Adaptively changes the size of the canvas, returns true on success.
func (self *Canvas) AdaptiveResize(width uint, height uint) error {
	success := C.MagickAdaptiveResizeImage(self.wand, C.size_t(width), C.size_t(height))

	if success == C.MagickFalse {
		return fmt.Errorf("Could not resize: %s", self.Error())
	}

	return nil
}

// Changes the compression quality of the canvas. Ranges from 1 (lowest) to 100 (highest).
func (self *Canvas) SetQuality(quality uint) error {
	success := C.MagickSetImageCompressionQuality(self.wand, C.size_t(quality))

	if success == C.MagickFalse {
		return fmt.Errorf("Could not set compression quality: %s", self.Error())
	}

	return nil
}

// Returns the compression quality of the canvas. Ranges from 1 (lowest) to 100 (highest).
func (self *Canvas) Quality() uint {
	return uint(C.MagickGetImageCompressionQuality(self.wand))
}

/*
// Sets canvas's foreground color.
func (self *Canvas) SetColor(color string) (bool) {
  status := C.PixelSetColor(self.fg, C.CString(color))
  if status == C.MagickFalse {
    return false
  }
  return true
}
*/

// Sets canvas' background color.
func (self *Canvas) SetBackgroundColor(color string) error {
	var status C.MagickBooleanType

	ccolor := C.CString(color)
	status = C.PixelSetColor(self.bg, ccolor)
	C.free(unsafe.Pointer(ccolor))

	if status == C.MagickFalse {
		return fmt.Errorf("Could not set pixel color: %s", self.Error())
	}

	status = C.MagickSetImageBackgroundColor(self.wand, self.bg)

	if status == C.MagickFalse {
		return fmt.Errorf("Could not set background color: %s", self.Error())
	}

	return nil
}

// Returns canvas' background color.
func (self *Canvas) BackgroundColor() string {
	return getPixelHexColor(self.bg)
}

// Sets antialiasing setting for the current drawing stroke.
func (self *Canvas) SetStrokeAntialias(value bool) {
	C.DrawSetStrokeAntialias(self.drawing, magickBoolean(value))
}

// Returns antialiasing setting for the current drawing stroke.
func (self *Canvas) StrokeAntialias() bool {
	value := C.DrawGetStrokeAntialias(self.drawing)
	if value == C.MagickTrue {
		return true
	}
	return false
}

// Sets the width of the stroke on the current drawing surface.
func (self *Canvas) SetStrokeWidth(value float64) {
	C.DrawSetStrokeWidth(self.drawing, C.double(value))
}

// Returns the width of the stroke on the current drawing surface.
func (self *Canvas) StrokeWidth() float64 {
	return float64(C.DrawGetStrokeWidth(self.drawing))
}

// Sets the opacity of the stroke on the current drawing surface.
func (self *Canvas) SetStrokeOpacity(value float64) {
	C.DrawSetStrokeOpacity(self.drawing, C.double(value))
}

// Returns the opacity of the stroke on the current drawing surface.
func (self *Canvas) StrokeOpacity() float64 {
	return float64(C.DrawGetStrokeOpacity(self.drawing))
}

// Sets the type of the line cap on the current drawing surface.
func (self *Canvas) SetStrokeLineCap(value uint) {
	C.DrawSetStrokeLineCap(self.drawing, C.LineCap(value))
}

// Returns the type of the line cap on the current drawing surface.
func (self *Canvas) StrokeLineCap() uint {
	return uint(C.DrawGetStrokeLineCap(self.drawing))
}

// Sets the type of the line join on the current drawing surface.
func (self *Canvas) SetStrokeLineJoin(value uint) {
	C.DrawSetStrokeLineJoin(self.drawing, C.LineJoin(value))
}

// Returns the type of the line join on the current drawing surface.
func (self *Canvas) StrokeLineJoin() uint {
	return uint(C.DrawGetStrokeLineJoin(self.drawing))
}

/*
func (self *Canvas) SetFillRule(value int) {
  C.DrawSetFillRule(self.drawing, C.FillRule(value))
}
*/

// Sets the fill color for enclosed areas on the current drawing surface.
func (self *Canvas) SetFillColor(color string) {
	ccolor := C.CString(color)
	C.PixelSetColor(self.fill, ccolor)
	C.free(unsafe.Pointer(ccolor))
	C.DrawSetFillColor(self.drawing, self.fill)
}

// Returns the fill color for enclosed areas on the current drawing surface.
func (self *Canvas) FillColor() string {
	return getPixelHexColor(self.fill)
}

// Sets the stroke color on the current drawing surface.
func (self *Canvas) SetStrokeColor(color string) {
	ccolor := C.CString(color)
	C.PixelSetColor(self.stroke, ccolor)
	C.free(unsafe.Pointer(ccolor))
	C.DrawSetStrokeColor(self.drawing, self.stroke)
}

// Returns the stroke color on the current drawing surface.
func (self *Canvas) StrokeColor() string {
	return getPixelHexColor(self.stroke)
}

// Draws a circle over the current drawing surface.
func (self *Canvas) Circle(radius float64) {
	C.DrawCircle(self.drawing, C.double(0), C.double(0), C.double(radius), C.double(0))
}

// Draws a rectangle over the current drawing surface.
func (self *Canvas) Rectangle(x float64, y float64) {
	C.DrawRectangle(self.drawing, C.double(0), C.double(0), C.double(x), C.double(y))
}

// Moves the current coordinate system origin to the specified coordinate.
func (self *Canvas) Translate(x float64, y float64) {
	C.DrawTranslate(self.drawing, C.double(x), C.double(y))
}

// Applies a scaling factor to the units of the current coordinate system.
func (self *Canvas) Scale(x float64, y float64) {
	C.DrawScale(self.drawing, C.double(x), C.double(y))
}

// Draws a line starting on the current coordinate system origin and ending on the specified coordinates.
func (self *Canvas) Line(x float64, y float64) {
	C.DrawLine(self.drawing, C.double(0), C.double(0), C.double(x), C.double(y))
}

/*
func (self *Canvas) Skew(x float64, y float64) {
  C.DrawSkewX(self.drawing, C.double(x))
  C.DrawSkewY(self.drawing, C.double(y))
}
*/

// Applies a rotation of a given angle (in radians) on the current coordinate system.
func (self *Canvas) Rotate(rad float64) {
	deg := RAD_TO_DEG * rad
	C.DrawRotate(self.drawing, C.double(deg))
}

// Draws an ellipse centered at the current coordinate system's origin.
func (self *Canvas) Ellipse(a float64, b float64) {
	C.DrawEllipse(self.drawing, C.double(0), C.double(0), C.double(a), C.double(b), 0, 360)
}

// Clones the current drawing surface and stores it in a stack.
func (self *Canvas) PushDrawing() error {
	success := C.PushDrawingWand(self.drawing)

	if success == C.MagickFalse {
		return fmt.Errorf("Could not push surface: %s", self.Error())
	}

	return nil
}

// Destroys the current drawing surface and returns the latest surface that was pushed to the stack.
func (self *Canvas) PopDrawing() error {
	success := C.PopDrawingWand(self.drawing)

	if success == C.MagickFalse {
		return fmt.Errorf("Could not pop surface: %s", self.Error())
	}

	return nil
}

// Copies a drawing surface to the canvas.
func (self *Canvas) Update() error {
	success := C.MagickDrawImage(self.wand, self.drawing)

	if success == C.MagickFalse {
		return fmt.Errorf("Could not update image: %s", self.Error())
	}

	return nil
}

// Destroys canvas.
func (self *Canvas) Destroy() error {

	if self.bg != nil {
		C.DestroyPixelWand(self.bg)
		self.bg = nil
	}

	if self.fg != nil {
		C.DestroyPixelWand(self.fg)
		self.fg = nil
	}

	if self.stroke != nil {
		C.DestroyPixelWand(self.stroke)
		self.stroke = nil
	}

	if self.fill != nil {
		C.DestroyPixelWand(self.fill)
		self.fill = nil
	}

	if self.drawing != nil {
		C.DestroyDrawingWand(self.drawing)
		self.drawing = nil
	}

	if self.wand == nil {
		return errors.New("Nothing to destroy")
	} else {
		C.DestroyMagickWand(self.wand)
		self.wand = nil
	}

	return nil
}

func Finalize() {
	C.MagickWandTerminus()
}

// Creates an empty canvas of the given dimensions.
func (self *Canvas) Blank(width uint, height uint) error {
	success := C.MagickNewImage(self.wand, C.size_t(width), C.size_t(height), self.bg)

	if success == C.MagickFalse {
		return fmt.Errorf("Could not create image: %s", self.Error())
	}

	return nil
}

// Convolves the canvas with a Gaussian function given its standard deviation.
func (self *Canvas) Blur(sigma float64) error {
	success := C.MagickBlurImage(self.wand, C.double(0), C.double(sigma))

	if success == C.MagickFalse {
		return fmt.Errorf("Could not blur image: %s", self.Error())
	}

	return nil
}

// Adaptively blurs the image by blurring less intensely near the edges and more intensely far from edges.
func (self *Canvas) AdaptiveBlur(sigma float64) error {
	success := C.MagickAdaptiveBlurImage(self.wand, C.double(0), C.double(sigma))

	if success == C.MagickFalse {
		return fmt.Errorf("Could not blur image: %s", self.Error())
	}

	return nil
}

// Adds random noise to the canvas.
func (self *Canvas) AddNoise() error {
	success := C.MagickAddNoiseImage(self.wand, C.GaussianNoise)

	if success == C.MagickFalse {
		return fmt.Errorf("Could not add noise: %s", self.Error())
	}

	return nil
}

// Removes a region of a canvas and collapses the canvas to occupy the removed portion.
func (self *Canvas) Chop(x int, y int, width uint, height uint) error {
	success := C.MagickChopImage(self.wand, C.size_t(width), C.size_t(height), C.ssize_t(x), C.ssize_t(y))

	if success == C.MagickFalse {
		return fmt.Errorf("Could not chop: %s", self.Error())
	}

	return nil
}

// Extracts a region from the canvas.
func (self *Canvas) Crop(x int, y int, width uint, height uint) error {
	success := C.MagickCropImage(self.wand, C.size_t(width), C.size_t(height), C.ssize_t(x), C.ssize_t(y))

	if success == C.MagickFalse {
		return fmt.Errorf("Could not crop: %s", self.Error())
	}

	return nil
}

func (self *Canvas) SetSize(width, height uint) error {
	if C.MagickSetSize(self.wand, C.size_t(width), C.size_t(height)) == C.MagickFalse {
		return fmt.Errorf("Could not set size: %s", self.Error())
	}

	return nil
}

func (self *Canvas) SetContrast(factor float64) error {
	factor = math.Max(-100, factor)
	factor = math.Min(100, factor)

	success := C.MagickBrightnessContrastImage(self.wand, 0, C.double(factor))

	if success == C.MagickFalse {
		return fmt.Errorf("Could not set contrast: %s", self.Error())
	}

	return nil
}

// Adjusts the canvas's brightness given a factor (-1.0 thru 1.0)
func (self *Canvas) SetBrightness(factor float64) error {
	factor = math.Max(-1, factor)
	factor = math.Min(1, factor)

	success := C.MagickModulateImage(self.wand, C.double(100+factor*100.0), C.double(100), C.double(100))

	if success == C.MagickFalse {
		return fmt.Errorf("Could not set brightness: %s", self.Error())
	}

	return nil
}

// Adjusts the canvas's saturation given a factor (-1.0 thru 1.0)
func (self *Canvas) SetSaturation(factor float64) error {

	factor = math.Max(-1, factor)
	factor = math.Min(1, factor)

	success := C.MagickModulateImage(self.wand, C.double(100), C.double(100+factor*100.0), C.double(100))

	if success == C.MagickFalse {
		return fmt.Errorf("Could not set saturation: %s", self.Error())
	}

	return nil
}

// Adjusts the canvas's hue given a factor (-1.0 thru 1.0)
func (self *Canvas) SetHue(factor float64) error {

	factor = math.Max(-1, factor)
	factor = math.Min(1, factor)

	success := C.MagickModulateImage(self.wand, C.double(100), C.double(100), C.double(100+factor*100.0))

	if success == C.MagickFalse {
		return fmt.Errorf("Could not set hue: %s", self.Error())
	}

	return nil
}

func (self *Canvas) InterlaceScheme() uint {
	return uint(C.MagickGetImageInterlaceScheme(self.wand))
}

func (self *Canvas) SetInterlaceScheme(scheme uint) error {
	if C.MagickSetImageInterlaceScheme(self.wand, C.InterlaceType(scheme)) == C.MagickFalse {
		return fmt.Errorf("Could not set interlace scheme: %s", self.Error())
	}

	return nil
}

// Sets the format of a particular image
func (self *Canvas) SetFormat(format string) error {
	cformat := C.CString(format)
	defer C.free(unsafe.Pointer(cformat))

	if C.MagickSetImageFormat(self.wand, cformat) == C.MagickFalse {
		return fmt.Errorf("Could not set format: %s", self.Error())
	}

	return nil
}

func (self *Canvas) GetFormat() string {
	return self.Format()
}

func (self *Canvas) Format() string {
	ptr := C.MagickGetImageFormat(self.wand)
	defer C.free(unsafe.Pointer(ptr))

	return C.GoString(ptr)
}

func (self *Canvas) Strip() error {
	var status C.MagickBooleanType

	status = C.MagickStripImage(self.wand)

	if status == C.MagickFalse {
		return fmt.Errorf("Could not strip: %s", self.Error())
	}

	return nil
}

func (self *Canvas) SetType(imageType uint) error {
	var status C.MagickBooleanType

	status = C.MagickSetImageType(self.wand, C.ImageType(imageType))

	if status == C.MagickFalse {
		return fmt.Errorf("Could not set type: %s", self.Error())
	}

	return nil
}

func (self *Canvas) Type() uint {
	return uint(C.MagickGetImageType(self.wand))
}

func (self *Canvas) SetSepiaTone(threshold float64) error {
	threshold = math.Max(0.0, threshold)
	threshold = math.Min(100.0, threshold)
	threshold = (float64(self.QuantumRange()) * threshold) / 100.0

	if C.MagickSepiaToneImage(self.wand, C.double(threshold)) == C.MagickFalse {
		return fmt.Errorf("Could not apply sepia effect: %s", self.Error())
	}

	return nil
}

func (self *Canvas) QuantumRange() uint {
	if self.quantumRange == 0 {
		var quantumRange C.size_t = 0

		C.GetMagickQuantumRange(&quantumRange)

		self.quantumRange = uint(quantumRange)
	}

	return self.quantumRange
}

func (self *Canvas) SetFontFamily(name string) error {
	family := C.CString(name)
	defer C.free(unsafe.Pointer(family))

	if C.MagickSetFont(self.wand, family) == C.MagickFalse {
		return fmt.Errorf("Could not set font family: %s", self.Error())
	}

	return nil
}

func (self *Canvas) SetFontSize(size float64) {
	C.MagickSetPointsize(self.wand, C.double(size))
}

func (self *Canvas) SetGravity(gravity uint) {
	C.MagickSetGravity(self.wand, C.GravityType(gravity))
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
