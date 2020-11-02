package main

import (
	"fmt"
	"github.com/golang/freetype"
	cg "github.com/iv-menshenin/crazy-grafica"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// I can set the dimensions of the canvas by specifying clear measures of length
	width := cg.Millimeters(210.0)
	height := cg.Millimeters(297.0)
	var fontSize float64 = 12

	canvasRect := cg.NewRectangle(cg.ZeroPixel(), cg.ZeroPixel(), width, height)
	backColor := color.RGBA{255, 255, 255, 255}
	frontColor := color.RGBA{46, 46, 46, 255}
	accentColor := color.RGBA{198, 46, 46, 255}
	// usePen := pen{color: frontColor, weight: int(mmToPix(0.3))}

	var palette = color.Palette{
		backColor,
		frontColor,
		accentColor,
	}
	img := image.NewPaletted(canvasRect, palette)
	for nn := 0; nn < img.Bounds().Max.X; nn++ {
		for mm := 0; mm < img.Bounds().Max.Y; mm++ {
			img.Set(nn, mm, backColor)
		}
	}

	fontFace, err := freetype.ParseFont(goregular.TTF)
	if err != nil {
		panic(err)
	}
	fontFaceB, err := freetype.ParseFont(gobold.TTF)
	if err != nil {
		panic(err)
	}

	myPen := cg.NewPen(frontColor, cg.Millimeters(0.25))
	accPen := cg.NewPen(accentColor, cg.Millimeters(0.25))
	canv := cg.NewCanvas(img, img.Rect)
	fontOpt := cg.OptionFont(fontFace, fontSize, myPen)
	fontAccentOpt := cg.OptionFont(fontFace, fontSize, accPen)
	tableMiddFont := cg.OptionFont(fontFace, 9, myPen)
	tableSmallFont := cg.OptionFont(fontFace, 8, myPen)
	tableBoldFont := cg.OptionFont(fontFaceB, 8, myPen)
	lastPoint := canv.Write(cg.Padding(cg.Millimeters(5), cg.Lines(
		cg.PaddingLeftRight(cg.Millimeters(10), cg.Lines(
			cg.Cols(
				cg.Text(fmt.Sprintf("Order #%d %s", rand.Int(), time.Now().Format("01.02.2006 15:04")), cg.OptionAlignment(cg.AlignLeft), fontOpt),
				cg.Text("DRAFT", cg.OptionAlignment(cg.AlignRight), fontAccentOpt),
			),
			cg.Text("555-345-65-66 Menshenin Igor", fontOpt),
		)),
		cg.FixedY(cg.Millimeters(2)),
		cg.Table(
			[]cg.TableColumn{
				cg.Column("PRODUCT NAME", .49, cg.OptionCentered(), tableMiddFont, cg.OptionAlignment(cg.AlignLeft)),
				cg.Column("COUNT", .06, cg.OptionCentered(), tableSmallFont, cg.OptionAlignment(cg.AlignCenter)),
				cg.Column("PRICE", .08, cg.OptionCentered(), tableSmallFont, cg.OptionAlignment(cg.AlignRight)),
				cg.Column("AMOUNT", .1, cg.OptionCentered(), tableSmallFont, cg.OptionAlignment(cg.AlignRight)),
				cg.Column("PRICE WS", .09, cg.OptionCentered(), tableSmallFont, cg.OptionAlignment(cg.AlignRight)),
				cg.Column("AMOUNT WS", .1, cg.OptionCentered(), tableSmallFont, cg.OptionAlignment(cg.AlignRight)),
				cg.Column("DELIVERY", .08, cg.OptionCentered(), tableSmallFont, cg.OptionAlignment(cg.AlignRight)),
			},
			cg.Cols(
				cg.Text("Some product name. Pretty stuff with long name. 60x90x12 box everywhere. Light-green. BEER", cg.OptionCentered(), tableSmallFont),
				cg.Text("20"),
				cg.Text("36.00"),
				cg.Text("720.00"),
				cg.Text("30.00"),
				cg.Text("600.00"),
				cg.Text(""),
			),
			cg.Cols(
				cg.Text("Foo and Bar company paper", cg.OptionCentered(), tableSmallFont),
				cg.Text("20"),
				cg.Text("36.00"),
				cg.Text("720.00"),
				cg.Text("30.00"),
				cg.Text("600.00"),
				cg.Text("120.00"),
			),
			cg.Cols(
				cg.Empty(),
				cg.ColSpan(cg.Text("1440.00", tableBoldFont, cg.OptionAlignment(cg.AlignRight)), 3),
				cg.ColSpan(cg.Text("1200.00", tableBoldFont, cg.OptionAlignment(cg.AlignRight)), 2),
				cg.Text("120.00", tableBoldFont, cg.OptionAlignment(cg.AlignRight)),
			),
			cg.Cols(
				cg.Text("Total including delivery and fee (retail)", tableBoldFont, cg.OptionAlignment(cg.AlignRight)),
				cg.ColSpan(cg.Text("1977.37", tableBoldFont, cg.OptionAlignment(cg.AlignRight)), 6),
			),
			cg.Cols(
				cg.Text("Total including delivery and fee (wholesale)", tableBoldFont, cg.OptionAlignment(cg.AlignRight)),
				cg.ColSpan(cg.Text("1812.22", tableBoldFont, cg.OptionAlignment(cg.AlignRight)), 6),
			),
		),
	)))

	out := img.SubImage(image.Rect(img.Bounds().Min.X, img.Bounds().Min.Y, img.Bounds().Max.X, lastPoint.Y))
	// Кодировать как PNG.
	outFile, _ := os.Create("/home/devalio/image.png")
	png.Encode(outFile, out)
	outFile.Close()
	outFile, _ = os.Create("/home/devalio/image.jpeg")
	jpeg.Encode(outFile, out, &jpeg.Options{Quality: 100})
	outFile.Close()
}
