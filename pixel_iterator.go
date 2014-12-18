package canvas

/*
#include <wand/MagickWand.h>

static PixelWand* get_pixel_wands_at(PixelWand** pixel_wands, size_t position) {
    return pixel_wands[position];
}
*/
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

type PixelIterator struct {
	iterator *C.PixelIterator
}

func (self *PixelIterator) NextRow() (row []*Pixel) {
	pixelCount := C.size_t(0)
	pixels := C.PixelGetNextIteratorRow(self.iterator, &pixelCount)

	if pixels == nil {
		return nil
	}

	for i := 0; i < int(pixelCount); i++ {
		wand := C.get_pixel_wands_at(pixels, C.size_t(i))
		row = append(row, &Pixel{wand: wand})
	}

	return
}

func (self *PixelIterator) Sync() error {
	if C.PixelSyncIterator(self.iterator) == C.MagickFalse {
		return fmt.Errorf("cannot sync iterator: %s", self.Error())
	}

	return nil
}

func (self *PixelIterator) Error() error {
	var exception C.ExceptionType

	ptr := C.PixelGetIteratorException(self.iterator, &exception)
	message := C.GoString(ptr)

	C.PixelClearIteratorException(self.iterator)
	C.MagickRelinquishMemory(unsafe.Pointer(ptr))

	return errors.New(message)
}

func (self *PixelIterator) Destroy() {
	C.DestroyPixelIterator(self.iterator)
}
