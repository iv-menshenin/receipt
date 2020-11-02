package crazy_grafica

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
	"image"
)

type (
	DrawText interface {
		WriteTo(Canvas, image.Rectangle) image.Point
		replaceOptions(options ...TextOption) DrawStruct
		defaultOptions(options ...TextOption) DrawStruct
	}
	TextOption interface {
		textOptInt() int
		tableColumnOptInt() int
	}
	textAlignment struct {
		alignment Alignment
	}
	textCentered struct{}
	textFont     struct {
		font     *truetype.Font
		fontSize float64
		usePen   pen
	}
)

// Text renders text to Canvas
func Text(s string, options ...TextOption) DrawStruct {
	return text{
		text:    s,
		options: options,
	}
}

func (t text) replaceOptions(options ...TextOption) DrawStruct {
	return text{
		text:    t.text,
		options: options,
	}
}

func (t text) defaultOptions(options ...TextOption) DrawStruct {
	if len(t.options) == 0 {
		return text{
			text:    t.text,
			options: options,
		}
	}
	return t
}

// OptionFont contains font settings
func OptionFont(font *truetype.Font, fontSize float64, usePen pen) TextOption {
	return textFont{
		font:     font,
		fontSize: fontSize,
		usePen:   usePen,
	}
}

// OptionAlignment lets you set horizontal alignment
//  AlignLeft, AlignRight, AlignCenter
func OptionAlignment(a Alignment) TextOption {
	return textAlignment{
		alignment: a,
	}
}

// OptionCentered means that the text needs to be vertically aligned in the center
func OptionCentered() TextOption {
	return textCentered{}
}

func (_ textAlignment) textOptInt() int {
	return 0
}

func (_ textCentered) textOptInt() int {
	return 0
}

func (_ textFont) textOptInt() int {
	return 0
}

func (_ textAlignment) tableColumnOptInt() int {
	return 0
}

func (_ textCentered) tableColumnOptInt() int {
	return 0
}

func (_ textFont) tableColumnOptInt() int {
	return 0
}

type (
	text struct {
		text    string
		options []TextOption
	}
)

func getDefaultFont() *truetype.Font {
	fontFace, err := freetype.ParseFont(goregular.TTF)
	if err != nil {
		panic(err)
	}
	return fontFace
}

func (t text) WriteTo(canvas Canvas, rect image.Rectangle) image.Point {
	var (
		font      *truetype.Font
		fontSize  float64 = defaultFontSize
		usePen            = defaultPen
		alignment         = cellAlignment{
			hAlign:    AlignLeft,
			vCentered: false,
		}
	)
	for _, opt := range t.options {
		if f, ok := opt.(textFont); ok {
			font = f.font
			usePen = f.usePen
			fontSize = f.fontSize
		}
		if a, ok := opt.(textAlignment); ok {
			alignment.hAlign = a.alignment
		}
		if _, ok := opt.(textCentered); ok {
			alignment.vCentered = true
		}
	}
	if font == nil {
		font = getDefaultFont()
	}
	y := fillTextIntoRect(t.text, canvas.img, font, fontSize, rect, alignment, usePen)
	return image.Point{X: rect.Max.X, Y: y}
}
