package rendr

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

func (img *rndImage) rndImageAlloc(channel, size0, size1 uint, dtype rndType) {
	doalloc:= false
	if img.Channel!= channel || img.Size[0] !=size0 || img.Size[1] !=size1 || img.Dtype!=dtype {
		doalloc = true
	}
	if doalloc {
		img.Data = make([]float64, channel*size0*size1)
		img.Channel = channel
		img.Size[0] = size0
		img.Size[1] = size1
		img.Dtype = dtype
	}
}

func rndImageSave(filename string, img *rndImage) {
	width:= int(img.Size[0])
	height := int(img.Size[1])
	upleft := image.Point{X:0, Y:0}
	lowRight := image.Point{X:width, Y:height}
	outImg := image.NewRGBA(image.Rectangle{Min:upleft, Max:lowRight})

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			index:= (y*width + x)*4
			color:= color.RGBA{uint8(rndQuantize(0, img.Data[index], 1, 256)),
				uint8(rndQuantize(0, img.Data[index+1], 1, 256)),
				uint8(rndQuantize(0, img.Data[index+2], 1, 256)),
				uint8(rndQuantize(0, img.Data[index+3], 1, 256)),
			}
			outImg.Set(x, y, color)
		}
	}
	f, _ := os.Create(filename)
	if err:=png.Encode(f, outImg); err!=nil {
		fmt.Println("Failed to encode the image!")
	}
}