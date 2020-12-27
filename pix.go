package receipt

import (
	"image"
	"math"
)

const (
	mmInch          = 25.4
	measStr         = ",`Ц|@"
	lineSpacing     = 1.25
	cellPadding     = 1.5
	defaultFontSize = 12
)

var dpi = 340.0

type (
	// Measure is a universal unit of measurement on canvas. Use this:
	//  Millimeters, Inches, Pixels
	Measure interface {
		toPixel() int
		toInch() float64
		toMillimeter() float64
	}
	pixels      int
	inch        float64
	millimeters float64
)

// SetDPI allows you to set the resolution (dots per inch)
func SetDPI(newDpi float64) {
	dpi = newDpi
}

// ZeroPixel is a distance of 0 pixels long
func ZeroPixel() Measure {
	return Pixels(0)
}

// NewRectangle allows you to create a rectangle using length measures:
//	Millimeters, Inches, Pixels
// example
//	NewRectangle(ZeroPixel(), ZeroPixel(), Inches(1.0), Inches(1.0)) // box with sides 1 inch
func NewRectangle(x0, y0, x1, y1 Measure) image.Rectangle {
	return image.Rect(x0.toPixel(), y0.toPixel(), x1.toPixel(), y1.toPixel())
}

// Pixels - measures the distance in pixels
func Pixels(pix int) Measure {
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

// Inches - measures the distance in inches
func Inches(i float64) Measure {
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

// Millimeters - measures the distance in millimeters
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
