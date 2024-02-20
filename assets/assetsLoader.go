package assets

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	"strings"

	"gioui.org/op/paint"
	"gioui.org/widget"
	"golang.org/x/image/draw"
)

var (
	MarkedFieldImg widget.Image
)

//go:embed imgs
var embededImgsDir embed.FS

var Images map[string]image.Image
var customWidgetCache map[string]widget.Image

func InitializeAssets() {
	fmt.Println("Initializing Assets")

	if Images == nil {
		Images = make(map[string]image.Image)
	}

	if customWidgetCache == nil {
		customWidgetCache = make(map[string]widget.Image)
	}

	imgsDir, err := embededImgsDir.ReadDir("imgs")
	if err != nil {
		panic(err)
	}

	for _, entry := range imgsDir {
		if entry.IsDir() {
			continue
		}

		indexOfDot := strings.LastIndex(entry.Name(), ".")
		key := entry.Name()[:indexOfDot]

		_, hasItem := Images[key]
		if !hasItem {
			file, _ := embededImgsDir.ReadFile(fmt.Sprintf("imgs/%s", entry.Name()))
			image, _, err := image.Decode(bytes.NewBuffer(file))

			if err != nil {
				fmt.Printf("Couldn't decode image (%s), due to this error: %s\n", entry.Name(), err.Error())
				continue
			}

			Images[key] = image
		}
	}

	initializeImages()
}

func initializeImages() {
	MarkedFieldImg = widget.Image{Src: paint.NewImageOp(Images["marked"]), Fit: widget.Cover}
}

func GetWidgetImage(name string, size int) widget.Image {
	key := fmt.Sprintf("%s-%d", name, size)

	res, has := customWidgetCache[key]
	if !has {
		customWidgetCache[key] = widget.Image{Src: paint.NewImageOp(resizeImage(name, size)), Fit: widget.Cover}
		res = widget.Image{Src: paint.NewImageOp(resizeImage(name, size)), Fit: widget.Cover}
	}

	return res
}

func resizeImage(name string, size int) image.Image {
	tempImg, hasImg := Images[name]
	if !hasImg {
		return nil
	}

	dst := image.NewRGBA(image.Rect(0, 0, size, size))

	draw.NearestNeighbor.Scale(dst, dst.Rect, tempImg, tempImg.Bounds(), draw.Over, nil)

	return dst
}
