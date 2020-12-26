package crazy_grafica

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"math"
	"strings"
)

type (
	Alignment int

	pen struct {
		color  color.Color
		weight int
	}
	cellAlignment struct {
		hAlign    Alignment
		vCentered bool
	}
)

const (
	AlignLeft Alignment = iota
	AlignRight
	AlignCenter
)

var defaultPen = pen{
	color:  color.Black,
	weight: 6,
}

func NewPen(color color.Color, w Measure) pen {
	return pen{
		color:  color,
		weight: w.toPixel(),
	}
}

func (p pen) textOptInt() int {
	return 0
}

func (p pen) tableColumnOptInt() int {
	return 0
}

func makeFontDrawer(
	dst draw.Image,
	fontData *truetype.Font,
	fontColor color.Color,
	fontSize float64,
	pos image.Point,
) *font.Drawer {
	return &font.Drawer{
		Dst: dst,
		Src: image.NewUniform(fontColor),
		Face: truetype.NewFace(fontData, &truetype.Options{
			Size:    fontSize,
			Hinting: font.HintingFull,
			DPI:     dpi,
		}),
		Dot: fixed.Point26_6{
			X: fixed.I(pos.X),
			Y: fixed.I(pos.Y),
		},
	}
}

func calcTextPositionX(
	rect image.Rectangle,
	textWidth fixed.Int26_6,
	align cellAlignment,
) fixed.Int26_6 {
	xPosition := fixed.I(rect.Min.X)
	switch align.hAlign {
	case AlignRight:
		xPosition = fixed.I(rect.Max.X) - textWidth
	case AlignCenter:
		xPosition += fixed.I((rect.Max.X-rect.Min.X)/2) - fixed.I(textWidth.Ceil()/2)
	}
	return xPosition
}

type textOnPos struct {
	text string
	dot  fixed.Point26_6
}

func splitAndFitToRectangle(
	drawer *font.Drawer,
	rect image.Rectangle,
	text string,
	align cellAlignment,
) []textOnPos {
	var (
		textBounds, _ = drawer.BoundString(measStr)
		textHeight    = textBounds.Max.Y - textBounds.Min.Y
		yPosition     = fixed.I(rect.Min.Y + textHeight.Ceil())
		maxYPosition  = fixed.I(rect.Min.Y + int(math.Round(float64(textHeight.Ceil())/1.5)))
	)
	if align.vCentered {
		yPosition += fixed.I((rect.Max.Y-rect.Min.Y)/2 - textHeight.Ceil()/2)
		if yPosition < maxYPosition {
			yPosition = maxYPosition
		}
	}
	calcPosition := func(width fixed.Int26_6) fixed.Point26_6 {
		return fixed.Point26_6{
			X: calcTextPositionX(rect, width, align),
			Y: yPosition,
		}
	}
	textSlice := make([]textOnPos, 0, 1)
	textWidth := drawer.MeasureString(text)
	if textWidth.Ceil() > rect.Max.X {
		textChains := strings.Split(text, " ")
		for {
			dot := fixed.Point26_6{}
			full := true
			textToWrite := ""
			for nn := 0; nn < len(textChains); nn++ {
				w := strings.Join(textChains[:nn+1], " ")
				textEnd := drawer.MeasureString(w)
				dot = calcPosition(textEnd)
				if textEnd.Ceil() > rect.Max.X {
					textToWrite = strings.Join(textChains[:nn], " ")
					textChains = textChains[nn:]
					full = false
					break
				}
			}
			if full {
				textToWrite = strings.Join(textChains, " ")
				textChains = nil
			}
			if textToWrite == "" {
				textToWrite = textChains[0]
				textChains = textChains[1:]
			}
			textSlice = append(textSlice, textOnPos{
				text: textToWrite,
				dot:  dot,
			})
			if len(textChains) == 0 {
				break
			} else {
				yPosition += fixed.I(int(math.Round(float64(textHeight.Ceil()) * lineSpacing)))
			}
		}
	} else {
		textSlice = append(textSlice, textOnPos{
			text: text,
			dot:  calcPosition(textWidth),
		})
	}
	return textSlice
}

func fillTextIntoRect(
	drawer *font.Drawer,
	text string,
	rect image.Rectangle,
	align cellAlignment,
) int {
	var yPosition fixed.Int26_6
	for _, s := range splitAndFitToRectangle(drawer, rect, text, align) {
		yPosition = s.dot.Y
		drawer.Dot = s.dot
		drawer.DrawString(s.text)
	}
	return yPosition.Ceil()
}

func drawRect(img draw.Image, rect image.Rectangle, usePen pen) {
	for x := rect.Min.X; x <= rect.Max.X; x++ {
		for t := 0; t < usePen.weight; t++ {
			img.Set(x, rect.Min.Y+t, usePen.color)
			img.Set(x, rect.Max.Y+t, usePen.color)
		}
	}
	for y := rect.Min.Y; y <= rect.Max.Y; y++ {
		for t := 0; t < usePen.weight; t++ {
			img.Set(rect.Min.X+t, y, usePen.color)
			img.Set(rect.Max.X+t, y, usePen.color)
		}
	}
}

func padRect(rect image.Rectangle, padding int) image.Rectangle {
	top := rect.Min.Y + padding
	bottom := rect.Max.Y - padding
	if !(bottom > top) {
		bottom = top + 1
	}
	newRect := image.Rect(rect.Min.X+padding, top, rect.Max.X-padding, bottom)
	return newRect
}

func padRect4(rect image.Rectangle, l, t, r, b int) image.Rectangle {
	top := rect.Min.Y + t
	bottom := rect.Max.Y - b
	if !(bottom > top) {
		bottom = top + 1
	}
	newRect := image.Rect(rect.Min.X+l, top, rect.Max.X-r, bottom)
	return newRect
}
