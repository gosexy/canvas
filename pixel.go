package canvas

/*
#include <wand/MagickWand.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

type Pixel struct {
	wand *C.PixelWand
}

func (self *Pixel) Red() float64 {
	return float64(C.PixelGetRed(self.wand))
}

func (self *Pixel) Green() float64 {
	return float64(C.PixelGetGreen(self.wand))
}

func (self *Pixel) Blue() float64 {
	return float64(C.PixelGetBlue(self.wand))
}

func (self *Pixel) SetColor(color string) error {
	ccolor := C.CString(color)
	defer C.free(unsafe.Pointer(ccolor))

	if C.PixelSetColor(self.wand, ccolor) == C.MagickFalse {
		return fmt.Errorf("Could not set color")
	}

	return nil
}

func (self *Pixel) Destroy() {
	C.DestroyPixelWand(self.wand)
}
