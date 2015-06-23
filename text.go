/* text.go contains functions for text annotation
TODO:
getters and setters for
- Stretch (redefine C type)
- Style (redefine C type)
- Resolution (two doubles)
- Decoration (redefine C type)
- Encoding (string)
-  DrawSetTextInterlineSpacing(DrawingWand *,const double),
-  DrawSetTextInterwordSpacing(DrawingWand *,const double),
-  DrawSetGravity(DrawingWand *wand,const GravityType gravity)
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
	Weight		uint
	// Style		C.StyleType
	// Resolution  [2]C.double
	Alignment	Alignment
	Antialias	bool
	// Decoration	C.DecorationType
	// Encoding	string
	Kerning		float64
	// Interline	float64
	// Interword	float64
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
		cweight := C.DrawGetFontWeight(self.drawing)
		calignment := C.DrawGetTextAlignment(self.drawing)
		cantialias := C.DrawGetTextAntialias(self.drawing)
		ckerning := C.DrawGetTextKerning(self.drawing)
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
			Weight: uint(cweight),
			Alignment: Alignment(calignment),
			Antialias: antialias,
			Kerning: float64(ckerning),
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
		self.SetFontFamily(def.Family)
		self.SetFont(def.Font, def.Size)
		self.SetFontWeight(def.Weight)
		self.SetTextAlignment(def.Alignment)
		self.SetTextAntialias(def.Antialias)
		self.SetTextKerning(def.Kerning)
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

// Sets canvas' default font weight
func (self *Canvas) SetFontWeight(weight uint) {
	self.text.Weight = weight
	C.DrawSetFontWeight(self.drawing, C.size_t(weight))		
}

// Returns canvas' current font weight
func (self *Canvas) FontWeight() uint {
	return self.text.Weight
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

// Sets canvas' default text antialiasing option.
func (self *Canvas) SetTextKerning(k float64) {
	self.text.Kerning = k
	C.DrawSetTextKerning(self.drawing, C.double(k))
}

// Returns the canvas' current text aligment
func (self *Canvas) TextKerning() float64 {
	return self.text.Kerning
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
