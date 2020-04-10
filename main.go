package main

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
)

func main() {
	//test()
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	width := 1072
	height := 701
	margin := 20

	dc := gg.NewContext(width, height)

	// load wold
	background, err := gg.LoadImage("wolf.png")
	if err != nil {
		return errors.Wrap(err, "load background image")
	}

	// load logos
	logos, err := gg.LoadImage("logos.png")
	if err != nil {
		return errors.Wrap(err, "load logos image")
	}

	// resize logos
	w := float64(logos.Bounds().Dx()) * 0.5
	h := float64(logos.Bounds().Dy()) * 0.5

	resizedLogos := resize.Resize(uint(w), uint(h), logos, resize.Lanczos2)

	dc.DrawImage(background, 0, 0)
	dc.DrawImage(resizedLogos, width-resizedLogos.Bounds().Dx()-margin, height-resizedLogos.Bounds().Dy()-margin)

	// draw image
	x := float64(margin)
	y := float64(margin)
	w = float64(dc.Width()/2) - float64(margin)
	h = float64(dc.Height()) - float64(2.0*margin)
	dc.SetColor(color.RGBA{0, 111, 123, 0})
	dc.DrawRectangle(x, y, w, h)
	//dc.Fill()

	// load fonts
	fontSize := 90.0
	fontPath := filepath.Join("fonts", "RobotoCondensed-Bold.ttf")
	if err := dc.LoadFontFace(fontPath, fontSize); err != nil {
		return errors.Wrap(err, "load font")
	}

	// write text
	dc.SetColor(color.White)
	s := "La nit del llop"
	textWidth, textHeight := dc.MeasureString(s)
	x = float64(margin)
	textHeight, _, err = adjustFont(fontSize, textWidth, w, dc, fontPath, textHeight, s, 10)
	if err != nil {
		panic(err)
	}
	y = textHeight
	dc.DrawStringAnchored(s, float64(margin)+(w/2), y, 0.5, 0.5)
	y = float64(margin) + textHeight

	// load fonts
	fontPath = filepath.Join("fonts", "RobotoCondensed-Regular.ttf")
	if err := dc.LoadFontFace(fontPath, 35); err != nil {
		return errors.Wrap(err, "load font")
	}
	s = "presenta"
	textWidth, textHeight = dc.MeasureString(s)
	y += textHeight + 15
	dc.DrawStringAnchored(s, float64(margin)+(w/2), y, 0.5, 0.5)

	s = `"Antes de los años terribles"`
	fontPath = filepath.Join("fonts", "RobotoCondensed-Bold.ttf")
	if err := dc.LoadFontFace(fontPath, 200); err != nil {
		return errors.Wrap(err, "load font")
	}
	textWidth, textHeight = dc.MeasureString(s)
	textHeight, fontSize, err = adjustFont(200, textWidth, w, dc, fontPath, textHeight, s, 10)
	if err != nil {
		panic(err)
	}
	y += textHeight + float64(margin)*12
	dc.DrawStringAnchored(s, float64(margin)+w/2, y, 0.5, 0.5)

	fontPath = filepath.Join("fonts", "RobotoCondensed-Regular.ttf")
	if err := dc.LoadFontFace(fontPath, 35); err != nil {
		return errors.Wrap(err, "load font")
	}
	s = "amb"
	_, textHeight = dc.MeasureString(s)
	y += textHeight + float64(margin)
	dc.DrawStringAnchored(s, float64(margin)+w/2, y, 0.5, 0.5)

	fontPath = filepath.Join("fonts", "RobotoCondensed-Bold.ttf")
	if err := dc.LoadFontFace(fontPath, fontSize); err != nil {
		return errors.Wrap(err, "load font")
	}
	s = "Victor del Arbol"
	textWidth, textHeight = dc.MeasureString(s)
	textHeight, fontSize, err = adjustFont(fontSize, textWidth, w, dc, fontPath, textHeight, s, 10)
	if err != nil {
		panic(err)
	}
	y += textHeight + float64(margin)
	dc.DrawStringAnchored(s, float64(margin)+w/2, y, 0.5, 0.5)

	s = "Dijous 19 de març a les 21h"
	textWidth, textHeight = dc.MeasureString(s)
	textWidth, textHeight = dc.MeasureString(s)
	textHeight, fontSize, err = adjustFont(fontSize, textWidth, w, dc, fontPath, textHeight, s, 10)
	if err != nil {
		panic(err)
	}
	y += textHeight + float64(margin)*2
	dc.DrawStringAnchored(s, float64(margin)+w/2, y, 0.5, 0.5)

	s = "a l'Orfeó Catalònia, sopar tertúlia amb l'autor"
	textWidth, textHeight = dc.MeasureString(s)
	textWidth, textHeight = dc.MeasureString(s)
	fmt.Println(fontSize, s)
	textHeight, fontSize, err = adjustFont(fontSize, textWidth, w, dc, fontPath, textHeight, s, 1)
	fmt.Println(fontSize, s)
	if err != nil {
		panic(err)
	}
	y += textHeight + float64(margin)
	dc.DrawStringAnchored(s, float64(margin)+w/2, y, 0.5, 0.5)

	s = "Contactar amb ernesto@projecte-loc.org"
	textWidth, textHeight = dc.MeasureString(s)
	y += textHeight + float64(margin)
	dc.DrawStringAnchored(s, float64(margin)+w/2, y, 0.5, 0.5)

	if err := dc.SavePNG("output.jpg"); err != nil {
		return errors.Wrap(err, "save png")
	}

	return nil
}

func adjustFont(fontSize float64, textWidth float64, w float64, dc *gg.Context, fontPath string, textHeight float64, s string, jump int) (float64, float64, error) {
	fmt.Println(textWidth, w)
	for i := int(fontSize); textWidth > w; i = i - jump {
		fontSize = float64(i)
		if err := dc.LoadFontFace(fontPath, fontSize); err != nil {
			return 0, 0, errors.Wrap(err, "load font")
		}
		textWidth, textHeight = dc.MeasureString(s)
	}
	return textHeight, fontSize, nil
}

const TEXT = "Call me Ishmael"

func test() {
	const W = 500
	const H = 700
	const P = 16
	dc := gg.NewContext(W, H)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.DrawLine(W/2, 0, W/2, H)
	dc.DrawLine(0, H/2, W, H/2)
	dc.DrawRectangle(P, P, W-P-P, H-P-P)
	dc.SetRGBA(0, 0, 1, 0.25)
	dc.SetLineWidth(3)
	dc.Stroke()
	dc.SetRGB(0, 0, 0)

	if err := dc.LoadFontFace("/Library/Fonts/Arial Bold.ttf", 18); err != nil {
		panic(err)
	}

	//x := float64(0)
	y := float64(P)
	when := "Dijous 19 de març a les 21h"
	_, textHeight := dc.MeasureString(when)

	dc.DrawStringAnchored(when, W/2, y, 0.5, 0.5)

	where := "a l'Orfeó Catalònia, sopar tertúlia amb l'autor"
	_, textHeight = dc.MeasureString(where)
	y += textHeight + 10
	dc.DrawStringAnchored(where, W/2, y, 0.5, 0.5)

	contact := "Contactar amb ernesto@projecte-loc.org"
	_, textHeight = dc.MeasureString(contact)
	y += textHeight + 10
	dc.DrawStringAnchored(contact, W/2, y, 0.5, 0.5)
	//
	//dc.DrawStringWrapped("UPPER MIDDLE", W/2, P, 0.5, 0, 0, 1.5, gg.AlignCenter)
	//if err := dc.LoadFontFace("/Library/Fonts/Arial.ttf", 12); err != nil {
	//	panic(err)
	//}

	dc.SavePNG("out.png")
}
