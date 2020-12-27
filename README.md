# receipt

## Introduction
Once I needed to develop a system in a short time that could print a receipt. To do this, you can use the standard library "image/draw".
But I always don't like to write a lot of the same type of code by hand and I tried to implement basic primitives such as "text in a rectangle", "rows" and "columns".
The result is this library. Perhaps the approach that is used here does not fit into GO canons, but it does its little functionality well.

## Usages

All drawing operations on the canvas are performed using the document structure description. Similar to DOM for HTML.
Canvas is the main object which simply stores a link to the image and the dimensions of the work area rectangle.

```go
    import cg "github.com/iv-menshenin/crazy-grafica"

	width := cg.Millimeters(210.0)
	height := cg.Millimeters(297.0)
	canvasRect := cg.NewRectangle(cg.ZeroPixel(), cg.ZeroPixel(), width, height)
	img := image.NewRGBA(canvasRect)
	canv := cg.NewCanvas(img, img.Rect)

```

In order to render the necessary information, we must parse the document into structures: rows and columns.

```go
    pt := canv.Write(
        cg.Padding(
            cg.Millimeters(5),
            cg.Lines(
    		    cg.PaddingLeftRight(
                    cg.Millimeters(10),
                    cg.Lines(
    			        cg.Cols(
                            cg.Text(fmt.Sprintf("Order #%d %s", rand.Int(), time.Now().Format("01.02.2006 15:04")), cg.OptionAlignment(cg.AlignLeft), fontOpt),
                            cg.Text("DRAFT", cg.OptionAlignment(cg.AlignRight), fontAccentOpt),
    			        ),
    			        cg.Text("555-345-65-66 Menshenin Igor", fontOpt),
    		        )
                ),
    		    cg.FixedY(cg.Millimeters(2)),
    ...
```

canv.Write renders the object to the Canvas, drawing is performed sequentially - from top to bottom from left to right. The canv.Write method returns the coordinate of the bottom-right point at which drawing ended.    
You can see how it works with an [example](https://github.com/iv-menshenin/receipt/blob/main/example/main.go).

## Elements

### Measurements

Three objects represent distance measures on the Canvas:
* Millimeters
* Inches
* Pixels

### Lines and Columns

The two objects represent vertical and horizontal layouts on the Canvas. They are inherently containers.
* Lines
* Cols

### Text

This object renders the text. It takes as arguments the text itself to be placed and options that allow you to control alignment and font.
* Text

Options:
* OptionFont
* OptionAlignment
* OptionCentered

Please note that if the text does not fit in length into the container in which it is located, then the lines will wrap by words.

### Fillers and Paddings

The following structures allow you to set padding inside the container
* Padding
* PaddingLeftRight

The following structures allow you to fill a container with a void. Behavior similar to previous functions, but not containers.
* Fixed
* FixedY
* Empty

### Table

The most complex (to date) structure allows you to implement the drawing of a table while maintaining the width of the columns in each row and applying fonts to all cells of the column.
* Table

Font, heading, and relative column width settings are performed by the function:
* Column

Line-by-line filling of the table with data (in text format) is performed using the following combinations of structures:
* Cols > Text
* Cols > ColSpan > Text

## Example

![example result](https://github.com/iv-menshenin/receipt/blob/main/example/image.png)
