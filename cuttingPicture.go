package main

import (
	"bytes"
	"image"
	"image/draw"
	"image/png"
)

func cutPicture(imgByte []byte, x, y int) []byte {

	img, _, err := image.Decode(bytes.NewReader(imgByte))
	check(err)

	buf := new(bytes.Buffer)

	dst := image.NewRGBA64(image.Rect(0, 0, x, y))

	startPoint := image.Point{
		X: img.Bounds().Min.X,
		Y: img.Bounds().Min.Y,
	}

	endPoint := startPoint.Add(image.Point{
		X: x,
		Y: y,
	})

	rect := image.Rectangle{
		Min: startPoint,
		Max: endPoint,
	}

	draw.Draw(dst, rect, img, startPoint, draw.Over)

	err = png.Encode(buf, dst)
	check(err)

	return buf.Bytes()

}
