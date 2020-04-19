package poster

import (
	"fmt"
	"image/color"
	"os"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
)

const (
	width          = 1280
	height         = 800
	margin         = 30.0
	LobsterTwoBold = "assets/fonts/LobsterTwo-Bold.ttf"
	RobotoLight    = "assets/fonts/RobotoCondensed-Light.ttf"
	RobotoBold     = "assets/fonts/RobotoCondensed-Bold.ttf"
)

func Draw(poster Poster) error {
	ctx := gg.NewContext(width, height)

	err := drawBackground(ctx, "assets/images/background.png")
	if err != nil {
		return err
	}

	err = drawLogos(ctx, "assets/images/logos.png")
	if err != nil {
		return err
	}

	err = drawPicture(ctx, poster)
	if err != nil {
		return err
	}

	err = drawText(ctx, poster)
	if err != nil {
		return err
	}

	err = ctx.SavePNG("cartel.png")
	if err != nil {
		return err
	}

	return nil
}

type Line struct {
	text      string
	marginTop float64
	fontSize  float64
	fontPath  string
}

func drawBackground(ctx *gg.Context, file string) error {
	background, err := gg.LoadImage(file)
	if err != nil {
		return err
	}

	ctx.DrawImage(background, 0, 0)

	return nil
}

func drawLogos(ctx *gg.Context, file string) error {
	logos, err := gg.LoadImage(file)
	if err != nil {
		return err
	}

	logos = resize.Thumbnail(
		uint(float64(logos.Bounds().Dx())*0.8),
		uint(logos.Bounds().Dx()),
		logos,
		resize.Lanczos3,
	)

	ctx.DrawImage(
		logos,
		width-logos.Bounds().Dx()-margin,
		height-logos.Bounds().Dy()-margin,
	)

	return nil
}

func drawPicture(ctx *gg.Context, poster Poster) error {
	filepath, err := poster.Picture()
	if err != nil {
		return err
	}

	pic, err := gg.LoadImage(filepath)
	if err != nil {
		return err
	}

	resizedPic := resize.Thumbnail(
		uint(pic.Bounds().Dx()),
		250,
		pic,
		resize.Lanczos3,
	)

	contentWidth := ctx.Width()/2 - margin
	ctx.DrawImageAnchored(resizedPic, margin+contentWidth/2, 185, 0.5, 0)

	err = os.Remove(filepath)
	if err != nil {
		return err
	}

	return nil
}

func drawText(ctx *gg.Context, poster Poster) error {
	ctx.SetColor(color.White)

	lines := []Line{
		{
			text:      "La nit del llop",
			marginTop: 25,
			fontSize:  90.0,
			fontPath:  LobsterTwoBold,
		},
		{
			text:      "presenta",
			marginTop: 25,
			fontSize:  25,
			fontPath:  RobotoLight,
		},
		{
			text:      fmt.Sprintf(`"%s"`, poster.Title),
			marginTop: 290,
			fontSize:  45,
			fontPath:  RobotoBold,
		},
		{
			text:      "amb",
			marginTop: 25,
			fontSize:  25,
			fontPath:  RobotoLight,
		},
		{
			text:      fmt.Sprintf("%s", poster.Guest),
			marginTop: 20,
			fontSize:  45,
			fontPath:  RobotoBold,
		},
		{
			text:      poster.When(),
			marginTop: 35,
			fontSize:  45,
			fontPath:  RobotoLight,
		},
		{
			text:      poster.Where(),
			marginTop: 20,
			fontSize:  45,
			fontPath:  RobotoLight,
		},
		{
			text:      "Contactar amb ernesto@projecte-loc.org",
			marginTop: 20,
			fontSize:  45,
			fontPath:  RobotoLight,
		},
	}

	contentWidth := float64(ctx.Width()/2 - margin)
	positionX := margin + contentWidth/2
	positionY := margin

	for _, line := range lines {
		err := ctx.LoadFontFace(line.fontPath, line.fontSize)
		if err != nil {
			return err
		}

		err = adjustFontSize(ctx, line, contentWidth)
		if err != nil {
			return err
		}

		positionY = calculatePositionY(ctx, line, positionY)
		ctx.DrawStringAnchored(line.text, positionX, positionY, 0.5, 0)
	}

	return nil
}

func calculatePositionY(ctx *gg.Context, line Line, lastPosition float64) float64 {
	_, textHeight := ctx.MeasureString(line.text)

	return lastPosition + line.marginTop + textHeight
}

func adjustFontSize(ctx *gg.Context, line Line, maxWidth float64) error {
	textWidth, _ := ctx.MeasureString(line.text)
	if textWidth <= maxWidth {
		return nil
	}

	fontSize := line.fontSize
	for fontSize = line.fontSize; textWidth > maxWidth; fontSize = fontSize - 1 {
		err := ctx.LoadFontFace(line.fontPath, fontSize)
		if err != nil {
			return err
		}

		textWidth, _ = ctx.MeasureString(line.text)
	}

	return nil
}
