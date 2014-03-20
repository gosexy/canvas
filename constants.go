package canvas

/*
#cgo CFLAGS: -fopenmp -I./_include
#cgo LDFLAGS: -lMagickWand -lMagickCore

#include <wand/MagickWand.h>
*/
import "C"

import (
	"math"
)

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
