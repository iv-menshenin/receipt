package receipt

import (
	"image"
)

type (
	PaddingStruct interface {
		getPaddingFnc() func(DrawStruct) DrawStruct
		drawContent() DrawStruct
	}
	fixedFiller struct {
		x int
		y int
	}
	padding struct {
		paddingRight  int
		paddingTop    int
		paddingLeft   int
		paddingBottom int
		content       DrawStruct
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

// PaddingLeftRight adds padding on all sides to the drawing object
func Padding4(pad Measure, d DrawStruct) DrawStruct {
	return padding{
		paddingRight:  pad.toPixel(),
		paddingTop:    pad.toPixel(),
		paddingLeft:   pad.toPixel(),
		paddingBottom: pad.toPixel(),
		content:       d,
	}
}

// PaddingLeftRight allows to adjust the padding on each side of the object border
func Padding(l, t, r, b Measure, d DrawStruct) DrawStruct {
	return padding{
		paddingRight:  r.toPixel(),
		paddingTop:    t.toPixel(),
		paddingLeft:   l.toPixel(),
		paddingBottom: b.toPixel(),
		content:       d,
	}
}

// PaddingLeftRight adds left and right padding to the drawing object
func PaddingLeftRight(pad Measure, d DrawStruct) DrawStruct {
	return padding{
		paddingRight:  pad.toPixel(),
		paddingTop:    0,
		paddingLeft:   pad.toPixel(),
		paddingBottom: 0,
		content:       d,
	}
}

// PaddingLeftRight adds left and right padding to the drawing object
func PaddingTopBottom(pad Measure, d DrawStruct) DrawStruct {
	return padding{
		paddingRight:  0,
		paddingTop:    pad.toPixel(),
		paddingLeft:   0,
		paddingBottom: pad.toPixel(),
		content:       d,
	}
}

func (p padding) getPaddingFnc() func(DrawStruct) DrawStruct {
	return func(content DrawStruct) DrawStruct {
		return Padding(
			Pixels(p.paddingLeft),
			Pixels(p.paddingTop),
			Pixels(p.paddingRight),
			Pixels(p.paddingBottom),
			content,
		)
	}
}

func (p padding) drawContent() DrawStruct {
	return p.content
}

func (p padding) WriteTo(canvas Canvas, rect image.Rectangle) image.Point {
	result := p.content.WriteTo(canvas, padRect4(rect, p.paddingLeft, p.paddingTop, p.paddingRight, p.paddingBottom))
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
