package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const winWidth, winHeight int = 1080, 768

// Mutate changes a single r,g,b Node

func clear(pixels []byte) {
	for i := range pixels {
		pixels[i] = 0
	}
}

func round(num uint32) uint32 {
	val := num + uint32(math.Copysign(0.5, float64(num)))
	return val
}

func add(c color.Color, num uint32) color.Color {
	r, g, b, a := c.RGBA()
	return color.RGBA{uint8(r + num), uint8(g + num), uint8(b + num), uint8(a + num)}
}

func main() {

	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Printf("could not initialize sdl: %v", err)
		return
	}

	window, err := sdl.CreateWindow("Abstract Pictures", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Printf("could not create window: %v ", err)
		return
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Printf("could not create renderer: %v ", err)
		return
	}
	defer renderer.Destroy()

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")
	rand.Seed(time.Now().UTC().UnixNano())

	kitten, err := os.Open("kitten.png")
	if err != nil {
		fmt.Printf("could not open image: %v", err)
	}
	defer kitten.Close()

	img, err := png.Decode(kitten)
	if err != nil {
		panic("could not decode image: " + err.Error())
	}

	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y

	pixels := make([][]color.Color, h)

	for y := 0; y < h; y++ {
		pixels[y] = make([]color.Color, w)
		for x := 0; x < w; x++ {
			c := img.At(x, y)
			pixels[y][x] = c
		}
	}
	// for y := 1; y < h-1; y++ {
	// 	for x := 1; x < w-1; x++ {
	// c := pixels[y][x]
	// p := color.Palette(palette.Plan9)
	// newC := p.Convert(c)
	// oldR, oldG, oldB, oldA := c.RGBA()
	// newR, newG, newB, newA := newC.RGBA()
	// quantError := (oldR + oldG + oldB + oldA) - (newR + newG + newB + newA)
	// 		pixels[y][x+1] = pixels[y][x+1]
	// 		pixels[y+1][x-1] = pixels[y+1][x-1]
	// 		pixels[y+1][x] = pixels[y+1][x]
	// 		pixels[y+1][x+1] = pixels[y+1][x+1]
	// 	}
	// }
	// for each y from top to bottom
	//  for each x from left to right
	//     oldpixel  := pixel[x][y]
	//     newpixel  := find_closest_palette_color(oldpixel)
	//     pixel[x][y]  := newpixel
	//     quant_error  := oldpixel - newpixel
	//     pixel[x + 1][y    ] := pixel[x + 1][y    ] + quant_error * 7 / 16
	//     pixel[x - 1][y + 1] := pixel[x - 1][y + 1] + quant_error * 3 / 16
	//     pixel[x    ][y + 1] := pixel[x    ][y + 1] + quant_error * 5 / 16
	//     pixel[x + 1][y + 1] := pixel[x + 1][y + 1] + quant_error * 1 / 16

	m := image.NewRGBA(image.Rect(0, 0, w, h))
	for y, row := range pixels {
		for x, col := range row {
			m.Set(x, y, col)
		}
	}
	buf := new(bytes.Buffer)

	err = png.Encode(buf, m)
	if err != nil {
		panic(err)
	}

	bytes := buf.Bytes()

	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STATIC, int32(w), int32(h))
	if err != nil {
		panic("could not create texture from pixels: " + err.Error())
	}
	err = tex.Update(nil, bytes, w*8)
	if err != nil {
		panic(err)
	}
	renderer.Copy(tex, nil, nil)
	renderer.Present()
	for {
		// Needed for macOS to display window correctly
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

	}
}
