package main

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func loadImage(filepath string, x int, y int) (Image, error) {
	image, err := gg.LoadImage(filepath)
	if err != nil {
		return Image{}, err
	}

	return Image{
		image: image,
		x:     x,
		y:     y,
	}, nil
}

type Image struct {
	image image.Image
	x     int
	y     int
}

func (i Image) Image() image.Image {
	return i.image
}

func (i Image) Width() int {
	return i.image.Bounds().Dx()
}

func (i Image) Height() int {
	return i.image.Bounds().Dy()
}

func loadLogos() (Image, error) {
	logos, err := loadImage("images/logos.png", 0, 0)
	if err != nil {
		return Image{}, err
	}

	logos.image = resize.Resize(
		uint(logos.Width()/2),
		uint(logos.Height()/2),
		logos.image,
		resize.Lanczos2,
	)

	return logos, nil
}

const (
	width  = 1072
	height = 701
	margin = 20.0
	title  = "La nit del llop"
)

func run() error {
	ctx := gg.NewContext(width, height)

	err := drawBackground(ctx)
	if err != nil {
		return err
	}

	pic, err := loadImage("images/foto.png", 0, 0)
	if err != nil {
		return err
	}

	// draw image
	x := margin
	y := margin
	contentWidth := float64(ctx.Width()/2) - (margin)
	h := float64(ctx.Height()) - (2.0 * margin)
	ctx.SetColor(color.RGBA{0, 111, 123, 0})
	ctx.DrawRectangle(x, y, contentWidth, h)
	//ctx.Fill()

	// load fonts
	fontSize := 90.0
	fontPath := filepath.Join("fonts", "LobsterTwo-Bold.ttf")
	if err := ctx.LoadFontFace(fontPath, fontSize); err != nil {
		return err
	}

	ctx.SetColor(color.White)
	y += margin
	_, textHeight := ctx.MeasureString(title)
	textHeight, _, err = adjustFont(ctx, title, contentWidth, fontSize, fontPath, 10)
	if err != nil {
		return err
	}
	y += textHeight
	ctx.DrawStringAnchored(title, margin+contentWidth/2, y, 0.5, 0)

	fontPath = filepath.Join("fonts", "RobotoCondensed-Light.ttf")
	if err := ctx.LoadFontFace(fontPath, 25); err != nil {
		return errors.Wrap(err, "load font")
	}
	s := "presenta"
	_, textHeight = ctx.MeasureString(s)
	y += margin + textHeight
	ctx.DrawStringAnchored(s, margin+contentWidth/2, y, 0.5, 0)

	picSpace := 210.0
	resizedPic := resize.Thumbnail(uint(pic.Width()), uint(picSpace), pic.Image(), resize.Lanczos3)
	y += margin
	ctx.DrawImageAnchored(resizedPic, int(margin+contentWidth/2), int(y), 0.5, 0)

	s = `"Antes de los años terribles"`
	fontPath = filepath.Join("fonts", "RobotoCondensed-Bold.ttf")
	if err := ctx.LoadFontFace(fontPath, 200); err != nil {
		return err
	}
	//_, textHeight = ctx.MeasureString(s)
	textHeight, fontSize, err = adjustFont(ctx, s, contentWidth, 200, fontPath, 10)
	if err != nil {
		panic(err)
	}
	fmt.Println("f1", fontSize, s)
	y += textHeight + picSpace + margin
	ctx.DrawStringAnchored(s, margin+contentWidth/2, y, 0.5, 0)
	//
	fontPath = filepath.Join("fonts", "RobotoCondensed-Light.ttf")
	if err := ctx.LoadFontFace(fontPath, 25); err != nil {
		return errors.Wrap(err, "load font")
	}
	s = "amb"
	_, textHeight = ctx.MeasureString(s)
	y += textHeight + margin
	ctx.DrawStringAnchored(s, margin+contentWidth/2, y, 0.5, 0)

	s = "Victor del Arbol"
	fontPath = filepath.Join("fonts", "RobotoCondensed-Bold.ttf")
	if err := ctx.LoadFontFace(fontPath, 200); err != nil {
		return err
	}
	textHeight, fontSize, err = adjustFont(ctx, s, contentWidth, fontSize, fontPath, 10)
	fmt.Println("f3", fontSize, s)
	if err != nil {
		return err
	}
	y += textHeight + margin
	ctx.DrawStringAnchored(s, margin+contentWidth/2, y, 0.5, 0)

	fontPath = filepath.Join("fonts", "RobotoCondensed-Light.ttf")
	if err := ctx.LoadFontFace(fontPath, fontSize); err != nil {
		return errors.Wrap(err, "load font")
	}
	s = "Dijous 19 de març a les 21h"
	textHeight, fontSize, err = adjustFont(ctx, s, contentWidth, fontSize, fontPath, 10)
	if err != nil {
		return err
	}
	y += textHeight + (margin)*1 + 10
	ctx.DrawStringAnchored(s, float64(margin)+contentWidth/2, y, 0.5, 0)

	s = "a l'Orfeó Catalònia, sopar tertúlia amb l'autor"
	fmt.Println("fs", fontSize)
	textHeight, fontSize, err = adjustFont(ctx, s, contentWidth, fontSize, fontPath, 3)
	fmt.Println("fs", fontSize)
	if err != nil {
		return err
	}
	y += textHeight + (margin) + 5
	ctx.DrawStringAnchored(s, (margin)+contentWidth/2, y, 0.5, 0)

	s = "Contactar amb ernesto@projecte-loc.org"
	_, textHeight = ctx.MeasureString(s)
	y += textHeight + (margin)
	ctx.DrawStringAnchored(s, (margin)+contentWidth/2, y, 0.5, 0)

	if err := ctx.SavePNG("output.jpg"); err != nil {
		return errors.Wrap(err, "save png")
	}

	return nil
}

func drawBackground(ctx *gg.Context) error {
	background, err := loadImage("images/background.png", 0, 0)
	if err != nil {
		return err
	}

	logos, err := loadLogos()
	if err != nil {
		return err
	}

	ctx.DrawImage(background.image, 0, 0)
	ctx.DrawImage(logos.Image(), width-logos.Width()-margin, height-logos.Height()-margin)

	return nil
}

func adjustFont(ctx *gg.Context, text string, maxWidth float64, initialFontSize float64, fontPath string, step float64) (float64, float64, error) {
	textWidth, textHeight := ctx.MeasureString(text)
	fmt.Println("---", textWidth, maxWidth)
	if textWidth <= maxWidth {
		return textHeight, initialFontSize, nil
	}
	fontSize := initialFontSize
	endPosition := 0.0
	fmt.Println(text)
	for fontSize = initialFontSize; textWidth > maxWidth; fontSize = fontSize - step {
		fmt.Println("ft", fontSize)
		err := ctx.LoadFontFace(fontPath, fontSize)
		if err != nil {
			return 0, 0, err
		}

		textWidth, endPosition = ctx.MeasureString(text)
	}

	return endPosition, fontSize + step, nil
}
