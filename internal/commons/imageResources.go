package commons

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

func SknSelectResource(alias string) fyne.Resource {
	return sknImageByName(alias, false, false).Resource
}
func SknSelectThemedResource(alias string) fyne.Resource {
	return sknImageByName(alias, true, false).Resource
}

func SknSelectImage(alias string) *canvas.Image {
	return sknImageByName(alias, false, false)
}
func SknSelectThemedImage(alias string) *canvas.Image {
	return sknImageByName(alias, true, false)
}
func SknSelectThemedInvertedImage(alias string) *canvas.Image {
	return sknImageByName(alias, true, true)
}

// SknLoadImageFromURI loads an image or svg from filesystem
func SknLoadImageFromURI(u fyne.URI, themed bool) fyne.CanvasObject {
	read, err := storage.Reader(u)
	if err != nil {
		fmt.Println("Error opening image", err)
		return canvas.NewRectangle(color.Black)
	}
	res, err := storage.LoadResourceFromURI(read.URI())
	if err != nil {
		fmt.Println("Error reading image", err)
		return canvas.NewRectangle(color.Black)
	}
	img := canvas.NewImageFromResource(res)
	img.FillMode = canvas.ImageFillContain
	if themed {
		img.Resource = theme.NewThemedResource(img.Resource)
	}
	return img
}

func sknImageByName(alias string, themed bool, inverted bool) *canvas.Image {
	var selected fyne.Resource

	switch alias {
	case "sensorOn":
		selected = resourceSensorsOnMbr24pxSvg
	case "ThumbsUp":
		selected = resourceThumbsUpMdr24pxSvg
	case "ThumbsDown":
		selected = resourceThumbsDownMdr24pxSvg
	default:
		selected = resourceTimeLapseMbr24pxSvg
	}

	image := canvas.NewImageFromResource(selected)
	if themed {
		image = canvas.NewImageFromResource(theme.NewThemedResource(selected))
	}
	if inverted {
		image = canvas.NewImageFromResource(theme.NewInvertedThemedResource(selected))
	}

	image.FillMode = canvas.ImageFillContain
	image.ScaleMode = canvas.ImageScaleSmooth
	return image
}
