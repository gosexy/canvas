package canvas

import "testing"
import "math"

func TestOpenWrite(t *testing.T) {
  canvas := New()

  opened := canvas.Open("examples/input/example.png")

  if opened {
    canvas.SetQuality(90)
    canvas.Write("examples/output/example.jpg")
  }

  canvas.Destroy()
}

func TestResize(t *testing.T) {
  canvas := New()

  opened := canvas.Open("examples/input/example.png")

  if opened {
    canvas.Resize(100, 100)
    canvas.Write("examples/output/example-100x100.png")
  }

  canvas.Destroy()
}

func TestBlank(t *testing.T) {
  canvas := New()

  canvas.SetBackgroundColor("#00ff00")

  success := canvas.Blank(400, 400)

  if success {
    canvas.Write("examples/output/example-blank.png")
  }

  canvas.Destroy()
}

func TestDrawLine(t *testing.T) {
  canvas := New()

  canvas.SetBackgroundColor("#000000")

  success := canvas.Blank(400, 400)

  if success {
    canvas.Translate(200, 200)
    canvas.SetStrokeWidth(10)
    canvas.SetStrokeColor("#ffffff")
    canvas.Line(100, 100)

    canvas.Write("examples/output/example-line.png")
  }

  canvas.Destroy()
}

func TestDrawCircle(t *testing.T) {
  canvas := New()

  canvas.SetBackgroundColor("#000000")

  success := canvas.Blank(400, 400)

  if success {

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

    canvas.Write("examples/output/example-circle.png")
  }

  canvas.Destroy()
}

func TestDrawRectangle(t *testing.T) {
  canvas := New()

  canvas.SetBackgroundColor("#000000")

  success := canvas.Blank(400, 400)

  if success {

    canvas.SetFillColor("#ff0000")

    canvas.Translate(200 - 50, 200 + 75)
    canvas.SetStrokeWidth(5)
    canvas.SetStrokeColor("#ffffff")
    canvas.Rectangle(100, -150)
    
    canvas.Write("examples/output/example-rectangle.png")
  }

  canvas.Destroy()
}

func TestDrawEllipse(t *testing.T) {
  canvas := New()

  success := canvas.Blank(400, 400)

  if success {

    canvas.SetFillColor("#ff0000")

    canvas.PushDrawing()
      canvas.Translate(200, 200)
      canvas.Rotate(math.Pi/3)
      canvas.Ellipse(50, 180)
    canvas.PopDrawing()
    
    canvas.SetFillColor("#ff00ff")
    
    canvas.PushDrawing()
      canvas.Translate(200, 200)
      canvas.Rotate(-math.Pi/3)
      canvas.Ellipse(25, 90)
    canvas.PopDrawing()
    
    canvas.Write("examples/output/example-ellipse.png")
  }

  canvas.Destroy()
}

func TestBlur(t *testing.T) {
  canvas := New()
  defer canvas.Destroy()

  opened := canvas.Open("examples/input/example.png")

  if (opened) {
    canvas.Blur(0.1)
    canvas.Write("examples/output/example-blur.png")
  }
}

func TestModulate(t *testing.T) {
  canvas := New()
  defer canvas.Destroy()

  opened := canvas.Open("examples/input/example.png")

  if (opened) {
    canvas.SetBrightness(-0.5)
    canvas.SetHue(0.2)
    canvas.SetSaturation(0.9)
    canvas.Write("examples/output/example-modulate.png")
  }
}

func TestAdaptive(t *testing.T) {

  canvas := New()
  defer canvas.Destroy()

  opened := canvas.Open("examples/input/example.png")

  if (opened) {
    canvas.AdaptiveBlur(1.2)
    canvas.AdaptiveResize(100, 100)
    canvas.Write("examples/output/example-adaptive.png")
  }
}

func TestNoise(t *testing.T) {
  canvas := New()
  defer canvas.Destroy()

  opened := canvas.Open("examples/input/example.png")

  if (opened) {
    canvas.AddNoise()
    canvas.Write("examples/output/example-noise.png")
  }
}

func TestChop(t *testing.T) {
  canvas := New()
  defer canvas.Destroy()

  opened := canvas.Open("examples/input/example.png")

  if (opened) {
    canvas.Chop(0, 0, 100, 50)
    canvas.Write("examples/output/example-chop.png")
  }
}

func TestCrop(t *testing.T) {
  canvas := New()
  defer canvas.Destroy()

  opened := canvas.Open("examples/input/example.png")

  if (opened) {
    canvas.Crop(100, 200, 200, 100)
    canvas.Write("examples/output/example-crop.png")
  }
}
