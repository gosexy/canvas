package canvas

import (
	"math"
	"testing"
    "os"
)

/*
  Example image is form Yuko Honda
  http://www.flickr.com/photos/yukop/6779040884/
*/

func TestFail(t *testing.T) {
	canvas := New()

	err := canvas.Open("_examples/input/example.dat")

	if err == nil {
		canvas.Destroy()
		t.Errorf("Test should have failed.")
	}
}

func TestMalformed(t *testing.T) {
	canvas := New()

	err := canvas.Open("_examples/input/malformed.png")

	if err == nil {
		canvas.Destroy()
		t.Errorf("Test should have failed.")
	}

}

func TestOpenWrite(t *testing.T) {
	canvas := New()

	err := canvas.Open("_examples/input/example.png")

	if err == nil {
		canvas.AutoOrientate()

		canvas.SetQuality(90)

		canvas.Write("_examples/output/example.jpg")
	} else {
		t.Errorf("Error: %s\n", err)
	}

	canvas.Destroy()
}

func TestOpenBlobWrite(t *testing.T) {
	canvas := New()

    file, err := os.Open("_examples/input/example.png")
    if err != nil {
		t.Errorf("Error: %s\n", err)
    }

    defer file.Close()

    buf := make([]byte, 1191346, 1191346)
    num, err := file.Read(buf)
    if err != nil {
		t.Errorf("Error: %s\n", err)
    }

	err = canvas.OpenBlob(buf, uint(num))

	if err == nil {
		canvas.AutoOrientate()

		canvas.SetQuality(90)

		canvas.Write("_examples/output/example.jpg")
	} else {
		t.Errorf("Error: %s\n", err)
	}

	canvas.Destroy()
}

func TestThumbnail(t *testing.T) {
	canvas := New()

	err := canvas.Open("_examples/input/example.png")

	if err == nil {
		canvas.AutoOrientate()

		canvas.Thumbnail(100, 100)

		canvas.Write("_examples/output/example-thumbnail.png")
	} else {
		t.Errorf("Error: %s\n", err)
	}

	canvas.Destroy()
}

func TestClone(t *testing.T) {
	canvas := New()

	err := canvas.Open("_examples/input/example.png")

	if err == nil {

		clone := canvas.Clone()
		clone.Resize(100, 100)
		clone.Write("_examples/output/cloned-100x100.png")

		clone.Destroy()

		canvas.Write("_examples/output/not-cloned.png")
	} else {
		t.Errorf("Error: %s\n", err)
	}

	canvas.Destroy()
}

func TestResize(t *testing.T) {
	canvas := New()

	err := canvas.Open("_examples/input/example.png")

	if err == nil {
		canvas.Resize(100, 100)
		canvas.Write("_examples/output/example-100x100.png")
	} else {
		t.Errorf("Error: %s\n", err)
	}

	canvas.Destroy()
}

func TestBlank(t *testing.T) {
	canvas := New()

	canvas.SetBackgroundColor("#00ff00")

	err := canvas.Blank(400, 400)

	if err == nil {
		canvas.Write("_examples/output/example-blank.png")
	} else {
		t.Errorf("Failed to create blank image.")
	}

	canvas.Destroy()
}

func TestSettersAndGetters(t *testing.T) {

	canvas := New()

	err := canvas.Blank(400, 400)

	if err != nil {
		t.Errorf("Could not create blank image.")
	}

	const backgroundColor = "#112233"

	canvas.SetBackgroundColor(backgroundColor)

	if gotBackgroundColor := canvas.BackgroundColor(); gotBackgroundColor != backgroundColor {
		t.Errorf("Got %s, expecting %s", gotBackgroundColor, backgroundColor)
	}

	const strokeAntialias = true

	canvas.SetStrokeAntialias(strokeAntialias)

	if gotStrokeAntialias := canvas.StrokeAntialias(); gotStrokeAntialias != strokeAntialias {
		t.Errorf("Got %t, expecting %t.", gotStrokeAntialias, strokeAntialias)
	}

	const strokeWidth = 2.0

	canvas.SetStrokeWidth(strokeWidth)

	if gotStrokeWidth := canvas.StrokeWidth(); gotStrokeWidth != strokeWidth {
		t.Errorf("Got %f, expecting %f.", gotStrokeWidth, strokeWidth)
	}

	const strokeOpacity = 1.0

	canvas.SetStrokeOpacity(strokeOpacity)

	if gotStrokeOpacity := canvas.StrokeOpacity(); gotStrokeOpacity != strokeOpacity {
		t.Errorf("Got %f, expecting %f.", gotStrokeOpacity, strokeOpacity)
	}

	strokeLineCap := STROKE_SQUARE_CAP

	canvas.SetStrokeLineCap(strokeLineCap)

	if gotStrokeLineCap := canvas.StrokeLineCap(); gotStrokeLineCap != strokeLineCap {
		t.Errorf("Got %d, expecting %d.", gotStrokeLineCap, strokeLineCap)
	}

	strokeLineJoin := STROKE_ROUND_JOIN

	canvas.SetStrokeLineJoin(strokeLineJoin)

	if gotStrokeLineJoin := canvas.StrokeLineJoin(); gotStrokeLineJoin != strokeLineJoin {
		t.Errorf("Got %d, expecting %d.", gotStrokeLineJoin, strokeLineJoin)
	}

	const fillColor = "#112233"

	canvas.SetFillColor(fillColor)

	if gotFillColor := canvas.FillColor(); gotFillColor != fillColor {
		t.Errorf("Got %s, expecting %s", gotFillColor, fillColor)
	}

	const strokeColor = "#112233"

	canvas.SetStrokeColor(strokeColor)

	if gotStrokeColor := canvas.StrokeColor(); gotStrokeColor != strokeColor {
		t.Errorf("Got %s, expecting %s", gotStrokeColor, strokeColor)
	}

	const quality = 76

	canvas.SetQuality(quality)

	if gotQuality := canvas.Quality(); gotQuality != quality {
		t.Errorf("Got %d, expecting %d", gotQuality, quality)
	}

	canvas.Destroy()
}

func TestDrawLine(t *testing.T) {

	canvas := New()

	canvas.SetBackgroundColor("#000000")

	err := canvas.Blank(400, 400)

	if err == nil {

		canvas.Translate(200, 200)
		canvas.SetStrokeWidth(10)
		canvas.SetStrokeColor("#ffffff")
		canvas.Line(100, 100)

		canvas.Write("_examples/output/example-line.png")
	} else {
		t.Errorf("Failed to create blank image.")
	}

	canvas.Destroy()
}

func TestDrawCircle(t *testing.T) {
	canvas := New()

	canvas.SetBackgroundColor("#000000")

	err := canvas.Blank(400, 400)

	if err == nil {

		canvas.SetFillColor("#ff0000")

		canvas.PushDrawing()
		canvas.Translate(200, 200)
		canvas.SetStrokeWidth(5)
		canvas.SetStrokeColor("#ffffff")
		canvas.Circle(100)
		canvas.PopDrawing()

		canvas.PushDrawing()
		canvas.Translate(100, 100)
		canvas.SetStrokeWidth(3)
		canvas.SetStrokeColor("#ffffff")
		canvas.Circle(20)
		canvas.PopDrawing()

		canvas.Write("_examples/output/example-circle.png")
	} else {
		t.Errorf("Failed to create blank image.")
	}

	canvas.Destroy()
}

func TestDrawRectangle(t *testing.T) {
	canvas := New()

	canvas.SetBackgroundColor("#000000")

	err := canvas.Blank(400, 400)

	if err == nil {

		canvas.SetFillColor("#ff0000")

		canvas.Translate(200-50, 200+75)
		canvas.SetStrokeWidth(5)
		canvas.SetStrokeColor("#ffffff")
		canvas.Rectangle(100, -150)

		canvas.Write("_examples/output/example-rectangle.png")
	} else {
		t.Errorf("Failed to create blank image.")
	}

	canvas.Destroy()
}

func TestDrawEllipse(t *testing.T) {
	canvas := New()

	err := canvas.Blank(400, 400)

	if err == nil {

		canvas.SetFillColor("#ff0000")

		canvas.PushDrawing()
		canvas.Translate(200, 200)
		canvas.Rotate(math.Pi / 3)
		canvas.Ellipse(50, 180)
		canvas.PopDrawing()

		canvas.SetFillColor("#ff00ff")

		canvas.PushDrawing()
		canvas.Translate(200, 200)
		canvas.Rotate(-math.Pi / 3)
		canvas.Ellipse(25, 90)
		canvas.PopDrawing()

		canvas.Write("_examples/output/example-ellipse.png")
	} else {
		t.Errorf("Failed to create blank image.")
	}

	canvas.Destroy()
}

func TestBlur(t *testing.T) {
	canvas := New()
	defer canvas.Destroy()

	err := canvas.Open("_examples/input/example.png")

	if err == nil {
		canvas.Blur(3)
		canvas.Write("_examples/output/example-blur.png")
	} else {
		t.Errorf("Failed to create blank image.")
	}
}

func TestModulate(t *testing.T) {
	canvas := New()
	defer canvas.Destroy()

	err := canvas.Open("_examples/input/example.png")

	if err == nil {
		canvas.SetBrightness(-0.5)
		canvas.SetHue(0.2)
		canvas.SetSaturation(0.9)
		canvas.Write("_examples/output/example-modulate.png")
	} else {
		t.Errorf("Failed to create blank image.")
	}
}

func TestAdaptive(t *testing.T) {

	canvas := New()
	defer canvas.Destroy()

	err := canvas.Open("_examples/input/example.png")

	if err == nil {
		canvas.AdaptiveBlur(1.2)
		canvas.AdaptiveResize(100, 100)
		canvas.Write("_examples/output/example-adaptive.png")
	} else {
		t.Errorf("Failed to create blank image.")
	}
}

func TestNoise(t *testing.T) {
	canvas := New()
	defer canvas.Destroy()

	err := canvas.Open("_examples/input/example.png")

	if err == nil {
		canvas.AddNoise()
		canvas.Write("_examples/output/example-noise.png")
	} else {
		t.Errorf("Failed to create blank image.")
	}
}

func TestChop(t *testing.T) {
	canvas := New()
	defer canvas.Destroy()

	err := canvas.Open("_examples/input/example.png")

	if err == nil {
		canvas.Chop(0, 0, 100, 50)
		canvas.Write("_examples/output/example-chop.png")
	} else {
		t.Errorf("Failed to create blank image.")
	}
}

func TestCrop(t *testing.T) {
	canvas := New()
	defer canvas.Destroy()

	err := canvas.Open("_examples/input/example.png")

	if err == nil {
		canvas.Crop(100, 200, 200, 100)
		canvas.Write("_examples/output/example-crop.png")
	} else {
		t.Errorf("Failed to create blank image.")
	}
}
