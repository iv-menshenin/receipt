package crazy_grafica

import (
	"image"
	"math"
)

const (
	dpi    = 960.0
	mmInch = 25.4

	measStr         = ",`Ц|@"
	lineSpacing     = 1.25
	cellPadding     = 1.5
	defaultFontSize = 12
)

type (
	Measure interface {
		toPixel() int
		toInch() float64
		toMillimeter() float64
	}
	pixels      int
	inch        float64
	millimeters float64
)

func ZeroPixel() Measure {
	return Pixel(0)
}

func NewRectangle(x0, y0, x1, y1 Measure) image.Rectangle {
	return image.Rect(x0.toPixel(), y0.toPixel(), x1.toPixel(), y1.toPixel())
}

func Pixel(pix int) Measure {
	return pixels(pix)
}

func (p pixels) toPixel() int {
	return int(p)
}

func (p pixels) toInch() float64 {
	return float64(p) / dpi
}

func (p pixels) toMillimeter() float64 {
	return (float64(p) / dpi) * mmInch
}

func Inch(i float64) Measure {
	return inch(i)
}

func (i inch) toPixel() int {
	return int(math.Round(inchToPix(float64(i))))
}

func (i inch) toInch() float64 {
	return float64(i)
}

func (i inch) toMillimeter() float64 {
	return float64(i) * mmInch
}

func Millimeters(m float64) Measure {
	return millimeters(m)
}

func (m millimeters) toPixel() int {
	return int(math.Round(mmToPix(float64(m))))
}

func (m millimeters) toInch() float64 {
	return mmToInch(float64(m))
}

func (m millimeters) toMillimeter() float64 {
	return float64(m)
}

func mmToInch(mm float64) float64 {
	return mm / mmInch
}

func inchToPix(i float64) float64 {
	return i * dpi
}

func mmToPix(mm float64) float64 {
	return inchToPix(mmToInch(mm))
}
