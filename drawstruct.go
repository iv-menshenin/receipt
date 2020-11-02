package crazy_grafica

import (
	"image"
	"image/draw"
)

type (
	Canvas struct {
		img  draw.Image
		rect image.Rectangle
	}
	DrawStruct interface {
		WriteTo(Canvas, image.Rectangle) image.Point
	}
	DrawColumns interface {
		WriteTo(Canvas, image.Rectangle) image.Point
		getColumnByNum(int) DrawStruct
	}
	lines struct {
		lines []DrawStruct
	}
	cols struct {
		cols []DrawStruct
	}
)

func NewCanvas(img draw.Image, rect image.Rectangle) Canvas {
	return Canvas{
		img:  img,
		rect: rect,
	}
}

func (c Canvas) Write(d DrawStruct) image.Point {
	return d.WriteTo(c, c.rect)
}

func Lines(d ...DrawStruct) DrawStruct {
	return lines{}.Add(d...)
}

func (l lines) Add(d ...DrawStruct) lines {
	return lines{
		lines: append(l.lines, d...),
	}
}

func (l lines) WriteTo(canvas Canvas, rect image.Rectangle) image.Point {
	rect.Max.Y = rect.Min.Y
	for _, d := range l.lines {
		point := d.WriteTo(canvas, rect)
		rect.Min.Y = point.Y
		rect.Max.Y = point.Y
	}
	return image.Point{X: rect.Max.X, Y: rect.Max.Y}
}

func Cols(d ...DrawStruct) DrawColumns {
	return cols{}.Add(d...)
}

func (c cols) Add(d ...DrawStruct) cols {
	return cols{
		cols: append(c.cols, d...),
	}
}

func (c cols) WriteTo(canvas Canvas, rect image.Rectangle) image.Point {
	bottom := rect.Max.Y
	for _, d := range c.cols {
		point := d.WriteTo(canvas, rect)
		rect.Min.X = point.X
		if bottom < point.Y {
			bottom = point.Y
			rect.Max.Y = bottom
		}
	}
	return image.Point{X: rect.Max.X, Y: rect.Max.Y}
}

func (c cols) getColumnByNum(num int) DrawStruct {
	if num < len(c.cols) {
		return c.cols[num]
	}
	return Empty()
}
