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

func fillTextIntoRect(
	text string,
	img draw.Image,
	fontData *truetype.Font,
	fontSize float64,
	rect image.Rectangle,
	align cellAlignment,
	usePen pen,
) int {
	fontDrawer := &font.Drawer{
		Dst: img,
		Src: image.NewUniform(usePen.color),
		Face: truetype.NewFace(fontData, &truetype.Options{
			Size:    fontSize,
			Hinting: font.HintingFull,
			DPI:     dpi,
		}),
		Dot: fixed.Point26_6{
			X: fixed.I(rect.Min.X),
			Y: fixed.I(rect.Min.Y),
		},
	}
	var (
		textBounds, _ = fontDrawer.BoundString(measStr)
		textHeight    = textBounds.Max.Y - textBounds.Min.Y
		xPosition     = fixed.I(rect.Min.X)
		yPosition     = fixed.I(rect.Min.Y + textHeight.Ceil())
		maxYPosition  = fixed.I(rect.Min.Y + int(math.Round(float64(textHeight.Ceil())/1.5)))
	)
	if align.vCentered {
		yPosition += fixed.I((rect.Max.Y-rect.Min.Y)/2 - textHeight.Ceil()/2)
		if yPosition < maxYPosition {
			yPosition = maxYPosition
		}
	}
	calcPosition := func(textToWrite string) {
		switch align.hAlign {
		case AlignRight:
			xPosition = fixed.I(rect.Max.X) - fontDrawer.MeasureString(textToWrite)
		case AlignCenter:
			xPosition += fixed.I((rect.Max.X-rect.Min.X)/2) - fixed.I(fontDrawer.MeasureString(textToWrite).Ceil()/2)
		}
		fontDrawer.Dot = fixed.Point26_6{
			X: xPosition,
			Y: yPosition,
		}
	}
	writeStr := func(textToWrite string) {
		fontDrawer.DrawString(textToWrite)
	}
	calcPosition(text)

	if fontDrawer.MeasureString(text).Ceil() > rect.Max.X {
		textChains := strings.Split(text, " ")
		for {
			full := true
			textToWrite := ""
			for nn := 0; nn < len(textChains); nn++ {
				w := strings.Join(textChains[:nn+1], " ")
				calcPosition(w)
				textBounds, _ = fontDrawer.BoundString(w)
				if textBounds.Max.X > fixed.I(rect.Max.X) {
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
			writeStr(textToWrite)
			if len(textChains) == 0 {
				break
			} else {
				yPosition += fixed.I(int(math.Round(float64(textHeight.Ceil()) * lineSpacing)))
			}
		}
	} else {
		writeStr(text)
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
