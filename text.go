/* text.go contains functions for text annotation
TODO:
getters and setters for
- Stretch (redefine C type)
- Weight (uint)
- Style (redefine C type)
- Resolution (two doubles)
- Decoration (redefine C type)
- Encoding (string)
*/

package canvas

/*
#cgo CFLAGS: -fopenmp -I./_include
#cgo LDFLAGS: -lMagickWand -lMagickCore

#include <wand/magick_wand.h>
*/
import "C"

import (
	"unsafe"
)

type Alignment uint

const (
 	UndefinedAlign Alignment	= Alignment(C.UndefinedAlign)
	LeftAlign					= Alignment(C.LeftAlign)
	CenterAlign					= Alignment(C.CenterAlign)
	RightAlign					= Alignment(C.RightAlign)
)

// structure containing all text properties for an annotation
// except the colors that are defined by FillColor and StrokeColor
type TextProperties struct {
	Font		string
	Family		string
	Size		float64
	// Stretch		C.StretchType
	// Weight		uint
	// Style		C.StyleType
	// Resolution  [2]C.double
	Alignment	Alignment
	Antialias	bool
	// Decoration	C.DecorationType
	// Encoding	string
	UnderColor	*C.PixelWand
}

// Returns a TextProperties structure.
// Parameters:
//   read_default: if false, returns an empty structure.
//				   if true, returns a structure set with current canvas settings
func (self *Canvas) NewTextProperties(read_default bool) *TextProperties {
	if read_default == true {
		cfont := C.DrawGetFont(self.drawing)
		defer C.free(unsafe.Pointer(cfont))
		cfamily := C.DrawGetFontFamily(self.drawing)
		defer C.free(unsafe.Pointer(cfamily))
		csize := C.DrawGetFontSize(self.drawing)
		calignment := C.DrawGetTextAlignment(self.drawing)
		cantialias := C.DrawGetTextAntialias(self.drawing)
		antialias := false
		if cantialias == C.MagickTrue {
			antialias = true
		}

		underColor :=C.NewPixelWand()
		C.DrawGetTextUnderColor(self.drawing, underColor)
		return &TextProperties{
			Font: C.GoString(cfont),
			Family: C.GoString(cfamily),
			Size: float64(csize),
			Alignment: Alignment(calignment),
			Antialias: antialias,
			UnderColor: underColor,
		}
	}
	return &TextProperties{
		UnderColor: C.NewPixelWand(),
	}
}

// Sets canvas' default TextProperties
func (self *Canvas) SetTextProperties(def *TextProperties) {
	if def != nil {
		self.text = def
		self.SetFont(def.Font, def.Size)
		self.SetFontFamily(def.Family)
		self.SetTextAlignment(def.Alignment)
		self.SetTextAntialias(def.Antialias)
	}
}

// Gets a copy of canvas' current TextProperties
func (self *Canvas) TextProperties() *TextProperties {
	if self.text == nil {
		return nil
	}
	cpy := *self.text
	return &cpy
}

// Sets canvas' default font name
func (self *Canvas) SetFontName(font string) {
	self.text.Font = font
	cfont := C.CString(font)
	defer C.free(unsafe.Pointer(cfont))
	C.DrawSetFont(self.drawing, cfont)
}

// Returns canvas' current font name
func (self *Canvas) FontName() string {
	return self.text.Font
}

// Sets canvas' default font family
func (self *Canvas) SetFontFamily(family string) {
	self.text.Family = family
	cfamily := C.CString(family)
	defer C.free(unsafe.Pointer(cfamily))
	C.DrawSetFontFamily(self.drawing, cfamily)	
}

// Returns canvas' current font family
func (self *Canvas) FontFamily() string {
	return self.text.Family
}

// Sets canvas' default font size
func (self *Canvas) SetFontSize(size float64) {
	self.text.Size = size
	C.DrawSetFontSize(self.drawing, C.double(size))		
}

// Returns canvas' current font size
func (self *Canvas) FontSize() float64 {
	return self.text.Size
}


// Sets canvas' font name and size.
// If font is 0-length, the current font family is not changed
// If size is <= 0, the current font size is not changed
func (self *Canvas) SetFont(font string, size float64) {
	if len(font) > 0 {
		self.SetFontName(font)
	}
	if size > 0 {
		self.SetFontSize(size)
	}
}

// Returns canvas' current font name and size
func (self *Canvas) Font() (string, float64) {
	return self.text.Font, self.text.Size
}

// Sets canvas' default text alignment. Available values are:
// UndefinedAlign (?), LeftAlign, CenterAlign, RightAlign
func (self *Canvas) SetTextAlignment(a Alignment) {
	self.text.Alignment = a
	C.DrawSetTextAlignment(self.drawing, C.AlignType(a))
}

// Returns the canvas' current text aligment
func (self *Canvas) TextAlignment() Alignment {
	return self.text.Alignment
}

// Sets canvas' default text antialiasing option.
func (self *Canvas) SetTextAntialias(b bool) {
	self.text.Antialias = b
	C.DrawSetTextAntialias(self.drawing, magickBoolean(b))
}

// Returns the canvas' current text aligment
func (self *Canvas) TextAntialias() bool {
	return self.text.Antialias
}

// Draws a string at the specified coordinates and using the current canvas
// Alignment.
func (self *Canvas) Annotate(text string, x, y float64) {
	c_text := C.CString(text)
	defer C.free(unsafe.Pointer(c_text))
	C.DrawAnnotation(self.drawing, C.double(x), C.double(y), (*C.uchar)(unsafe.Pointer(c_text)))
}

// Draws a string at the specified coordinates and using the specified Text Properties
// Does not modify the canvas' default TextProperties
func (self *Canvas) AnnotateWithProperties(text string, x, y float64, prop *TextProperties) {
	if prop != nil {
		tmp := self.TextProperties()
		self.SetTextProperties(prop)
		self.Annotate(text, x, y)
		self.SetTextProperties(tmp)
	} else {
		self.Annotate(text, x, y)
	}
}
