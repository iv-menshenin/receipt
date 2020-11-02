package crazy_grafica

import "image"

type (
	fixedFiller struct {
		x int
		y int
	}
	padding struct {
		paddingRight  int
		paddingTop    int
		paddingLeft   int
		paddingBottom int
		d             DrawStruct
	}
	empty struct{}
)

func Fixed(x, y Measure) DrawStruct {
	return fixedFiller{x: x.toPixel(), y: y.toPixel()}
}

func (f fixedFiller) WriteTo(canvas Canvas, rect image.Rectangle) image.Point {
	return image.Point{
		X: rect.Min.X + f.x,
		Y: rect.Min.Y + f.y,
	}
}

func FixedY(y Measure) DrawStruct {
	return Fixed(pixels(0), y)
}

func Padding(pad Measure, d DrawStruct) DrawStruct {
	return padding{
		paddingRight:  pad.toPixel(),
		paddingTop:    pad.toPixel(),
		paddingLeft:   pad.toPixel(),
		paddingBottom: pad.toPixel(),
		d:             d,
	}
}

func PaddingLeftRight(pad Measure, d DrawStruct) DrawStruct {
	return padding{
		paddingRight:  pad.toPixel(),
		paddingTop:    0,
		paddingLeft:   pad.toPixel(),
		paddingBottom: 0,
		d:             d,
	}
}

func (p padding) WriteTo(canvas Canvas, rect image.Rectangle) image.Point {
	result := p.d.WriteTo(canvas, padRect4(rect, p.paddingLeft, p.paddingTop, p.paddingRight, p.paddingBottom))
	return image.Point{X: rect.Min.X, Y: result.Y + p.paddingBottom}
}

func Empty() DrawStruct {
	return empty{}
}

func (e empty) WriteTo(canvas Canvas, rect image.Rectangle) image.Point {
	return image.Point{
		X: rect.Min.X,
		Y: rect.Min.Y,
	}
}
