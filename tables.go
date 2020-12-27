package receipt

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"image"
	"image/draw"
	"math"
)

type (
	ColumnOption interface {
		tableColumnOptInt() int
	}
	TableRow interface {
		getColumnByNum(int) DrawStruct
	}
	TableColumn interface {
		calculateWidth(tableWidth int) int
		getTextOptions() []TextOption
		getCaption() string
		getPen() pen
		makeFontDrawer(img draw.Image) *font.Drawer
		extractDrawStruct(DrawStruct) (DrawStruct, func(DrawStruct) DrawStruct)
	}
	ColumnSpan interface {
		spanCount() int
		drawContent() DrawStruct
	}
	tableColumn struct {
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
	return tableColumn{
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

func (c tableColumn) getTextOptions() []TextOption {
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

func (t tableColumn) calculateWidth(tableWidth int) int {
	return int(math.Round(float64(tableWidth) * t.pie))
}

func (t tableColumn) getCaption() string {
	return t.caption
}

func (t tableColumn) getPen() pen {
	return t.usePen
}

func (t tableColumn) makeFontDrawer(img draw.Image) *font.Drawer {
	return makeFontDrawer(img, t.font, t.usePen.color, t.fontSize)
}

func Table(columns []TableColumn, data ...TableRow) DrawStruct {
	return table{
		columns: columns,
		rows:    data,
	}
}

func (t tableColumn) extractDrawStruct(draw DrawStruct) (DrawStruct, func(DrawStruct) DrawStruct) {
	var padFunc = func(d DrawStruct) DrawStruct {
		// default cell padding
		return Padding4(millimeters(cellPadding), d)
	}
	for {
		switch v := draw.(type) {
		case PaddingStruct:
			draw = v.drawContent()
			// apply padding to cell
			padFunc = v.getPaddingFnc()
		case DrawText:
			draw = v.defaultOptions(t.getTextOptions()...)
			return draw, padFunc
		default:
			return draw, padFunc
		}
	}
}

func writeTableRow(
	t table,
	tableWidth int,
	left, top int,
	canvas Canvas,
	getColumnStruct func(int) DrawStruct,
) int {
	var (
		spanned   int
		colObjIdx int
		colWidth  int
		headRects = make([]image.Rectangle, 0, len(t.columns))
		bottom    = top + int(mmToPix(2))
	)
	for _, col := range t.columns {
		if spanned < 2 {
			colWidth += col.calculateWidth(tableWidth)
			colRect := image.Rect(left, top, colWidth+left, bottom)
			draw := getColumnStruct(colObjIdx)
			if d, ok := draw.(ColumnSpan); ok {
				draw = d.drawContent()
				if spanned < 1 {
					spanned = d.spanCount() - 1
					continue
				} else {
					spanned = 0
				}
			}
			draw, padFunc := col.extractDrawStruct(draw)
			end := padFunc(draw).WriteTo(canvas, colRect)
			left += colWidth
			if end.Y > bottom {
				bottom = end.Y
			}
			headRects = append(headRects, colRect)
			colWidth = 0
			colObjIdx++
		} else {
			colWidth += col.calculateWidth(tableWidth)
			spanned--
		}
	}
	for i, rect := range headRects {
		rect.Max.Y = bottom
		drawRect(canvas.img, rect, t.columns[i].getPen())
	}
	return bottom
}

func writeTableHeader(
	t table,
	canvas Canvas,
	rect image.Rectangle,
) int {
	var (
		tableWidth  = rect.Dx()
		left        = rect.Min.X
		top         = rect.Min.Y
		bottom      = top
		cellPadding = int(mmToPix(cellPadding))
		headRects   = make([]image.Rectangle, 0, len(t.columns))
	)
	for _, col := range t.columns {
		colWidth := col.calculateWidth(tableWidth)
		colRect := image.Rect(left, top, colWidth+left, top+int(mmToPix(5)))
		fontDrawer := col.makeFontDrawer(canvas.img)
		b := fillTextIntoRect(
			fontDrawer,
			col.getCaption(),
			colRect.Inset(cellPadding),
			cellAlignment{
				hAlign:    AlignCenter,
				vCentered: true,
			},
		)
		if b+cellPadding > bottom {
			bottom = b + cellPadding
		}
		left += colWidth
		headRects = append(headRects, colRect)
	}
	for i, rect := range headRects {
		rect.Max.Y = bottom
		drawRect(canvas.img, rect, t.columns[i].getPen())
	}
	return bottom
}

func (t table) WriteTo(canvas Canvas, rect image.Rectangle) image.Point {
	bottom := writeTableHeader(t, canvas, rect)
	for _, row := range t.rows {
		bottom = writeTableRow(t, rect.Dx(), rect.Min.X, bottom, canvas, row.getColumnByNum)
	}
	return image.Point{
		X: rect.Min.X,
		Y: bottom,
	}
}
