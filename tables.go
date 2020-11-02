package crazy_grafica

import (
	"github.com/golang/freetype/truetype"
	"image"
	"math"
)

type (
	ColumnOption interface {
		tableColumnOptInt() int
	}
	TableRow interface {
		getColumnByNum(int) DrawStruct
	}
	ColumnSpan interface {
		spanCount() int
		drawContent() DrawStruct
	}
	TableColumn struct {
		caption   string
		alignment textAlignment
		centered  bool
		font      *truetype.Font
		fontSize  float64
		usePen    pen
		pie       float64
	}
	table struct {
		columns []TableColumn
		rows    []TableRow
	}
	colSpan struct {
		draw DrawStruct
		span int
	}
)

func Column(caption string, pie float64, options ...ColumnOption) TableColumn {
	var (
		font      *truetype.Font
		fontSize  float64 = defaultFontSize
		usePen            = defaultPen
		alignment         = cellAlignment{
			hAlign:    AlignLeft,
			vCentered: false,
		}
	)
	for _, opt := range options {
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
	return TableColumn{
		caption:   caption,
		alignment: textAlignment{alignment: alignment.hAlign},
		centered:  alignment.vCentered,
		font:      font,
		fontSize:  fontSize,
		usePen:    usePen,
		pie:       pie,
	}
}

func ColSpan(d DrawStruct, span int) DrawStruct {
	return colSpan{
		draw: d,
		span: span,
	}
}

func (c colSpan) WriteTo(canvas Canvas, rect image.Rectangle) image.Point {
	return c.draw.WriteTo(canvas, rect)
}

func (c colSpan) spanCount() int {
	return c.span
}

func (c colSpan) drawContent() DrawStruct {
	return c.draw
}

func (c TableColumn) getTextOptions() []TextOption {
	if c.centered {
		return []TextOption{
			OptionCentered(),
			OptionFont(c.font, c.fontSize, c.usePen),
			OptionAlignment(c.alignment.alignment),
		}
	} else {
		return []TextOption{
			OptionFont(c.font, c.fontSize, c.usePen),
			OptionAlignment(c.alignment.alignment),
		}
	}
}

func Table(columns []TableColumn, data ...TableRow) DrawStruct {
	return table{
		columns: columns,
		rows:    data,
	}
}

func (t table) WriteTo(canvas Canvas, rect image.Rectangle) image.Point {
	tableWidth := rect.Max.X - rect.Min.X
	left := rect.Min.X
	top := rect.Min.Y
	bottom := top
	cellPadding := int(mmToPix(cellPadding))

	var headRects = make([]image.Rectangle, 0, len(t.columns))
	for _, col := range t.columns {
		colWidth := int(math.Round(float64(tableWidth) * col.pie))
		colRect := image.Rect(left, top, colWidth+left, top+int(mmToPix(5)))
		b := fillTextIntoRect(
			col.caption,
			canvas.img,
			col.font,
			col.fontSize,
			padRect(colRect, cellPadding),
			cellAlignment{
				hAlign:    AlignCenter,
				vCentered: true,
			},
			col.usePen,
		)
		if b+cellPadding > bottom {
			bottom = b + cellPadding
		}
		left += colWidth
		headRects = append(headRects, colRect)
	}
	for i, rect := range headRects {
		rect.Max.Y = bottom
		drawRect(canvas.img, rect, t.columns[i].usePen)
	}
	for _, row := range t.rows {
		headRects = make([]image.Rectangle, 0, len(t.columns))
		left = rect.Min.X
		top = bottom
		bottom += int(mmToPix(2))
		spanned := 0
		colObjIdx := 0
		colWidth := 0
		for _, col := range t.columns {
			if spanned < 2 {
				colWidth += int(math.Round(float64(tableWidth) * col.pie))
				colRect := image.Rect(left, top, colWidth+left, bottom)
				draw := row.getColumnByNum(colObjIdx)
				if d, ok := draw.(ColumnSpan); ok {
					draw = d.drawContent()
					if spanned < 1 {
						spanned = d.spanCount() - 1
						continue
					} else {
						spanned = 0
					}
				}
				if tx, ok := draw.(DrawText); ok {
					draw = tx.defaultOptions(col.getTextOptions()...)
				}
				end := Padding(pixels(cellPadding), draw).WriteTo(canvas, colRect)
				// DEBUG
				//drawRect(canvas.img, padRect(colRect, cellPadding), pen{
				//	color:  color.RGBA{198, 46, 46, 255},
				//	weight: 4,
				//})
				left += colWidth
				if end.Y > bottom {
					bottom = end.Y
				}
				headRects = append(headRects, colRect)
				colWidth = 0
				colObjIdx++
			} else {
				colWidth += int(math.Round(float64(tableWidth) * col.pie))
				spanned--
			}
		}
		for i, rect := range headRects {
			rect.Max.Y = bottom
			drawRect(canvas.img, rect, t.columns[i].usePen)
		}
	}
	return image.Point{
		X: left,
		Y: bottom,
	}
}
