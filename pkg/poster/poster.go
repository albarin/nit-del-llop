package poster

import (
	"fmt"
	"image/color"

	"github.com/albarin/poster/pkg/webhooks"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
)

const (
	width          = 1072
	height         = 701
	margin         = 20.0
	LobsterTwoBold = "assets/fonts/LobsterTwo-Bold.ttf"
	RobotoLight    = "assets/fonts/RobotoCondensed-Light.ttf"
	RobotoBold     = "assets/fonts/RobotoCondensed-Bold.ttf"
)

func Run(poster webhooks.Poster, background, logos, pic string) error {
	ctx := gg.NewContext(width, height)

	err := drawBackground(ctx, background)
	if err != nil {
		return err
	}

	err = drawLogos(ctx, logos)
	if err != nil {
		return err
	}

	err = drawPic(ctx, pic)
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

func drawText(ctx *gg.Context, poster webhooks.Poster) error {
	ctx.SetColor(color.White)

	lines := []Line{
		{
			text:      "La nit del llop",
			marginTop: margin,
			fontSize:  90.0,
			fontPath:  LobsterTwoBold,
		},
		{
			text:      "presenta",
			marginTop: margin,
			fontSize:  25,
			fontPath:  RobotoLight,
		},
		{
			text:      fmt.Sprintf(`"%s"`, poster.Title),
			marginTop: 250,
			fontSize:  45,
			fontPath:  RobotoBold,
		},
		{
			text:      "amb",
			marginTop: 15,
			fontSize:  25,
			fontPath:  RobotoLight,
		},
		{
			text:      fmt.Sprintf("%s", poster.Guest),
			marginTop: 15,
			fontSize:  45,
			fontPath:  RobotoBold,
		},
		{
			text:      poster.When(),
			marginTop: margin + 10,
			fontSize:  35,
			fontPath:  RobotoLight,
		},
		{
			text:      poster.Where(),
			marginTop: margin,
			fontSize:  35,
			fontPath:  RobotoLight,
		},
		{
			text:      "Contactar amb ernesto@projecte-loc.org",
			marginTop: margin,
			fontSize:  35,
			fontPath:  RobotoLight,
		},
	}

	contentWidth := float64(ctx.Width()/2 - margin*2)
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

func drawPic(ctx *gg.Context, file string) error {
	pic, err := gg.LoadImage(file)
	if err != nil {
		return err
	}

	resizedPic := resize.Thumbnail(
		uint(pic.Bounds().Dx()),
		210,
		pic,
		resize.Lanczos3,
	)

	contentWidth := ctx.Width()/2 + margin
	ctx.DrawImageAnchored(resizedPic, margin+contentWidth/2, 165, 0.5, 0)

	return nil
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
		uint(logos.Bounds().Dx()/2),
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
