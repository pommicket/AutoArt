/*
Copyright (C) 2019 Leo Tenenbaum

This file is part of AutoArt.

AutoArt is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

AutoArt is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with AutoArt.  If not, see <https://www.gnu.org/licenses/>.
*/

package autoart

import (
    "image"
    "image/color"
    "math"
	"math/rand"
	"fmt"
    "github.com/pommicket/autoart/autoutils"
)

const (
    XY = iota
    RTHETA
)

const (
    RGB = iota
    GRAYSCALE
    CMYK
    HSV
    YCbCr
)

const (
    MOD = iota
    CLAMP
    SIGMOID
)

type Config struct {
    FunctionLength int
    ColorSpace int
    CoordinateSys int
    Alpha bool
    Rectifier int // What to do with out-of-bounds values
}

func sigmoid(x float64) float64 {
    return 1 / (1 + math.Exp(-x))
}

func rectify(x float64, rectifier int) float64 {
    switch rectifier {
    case MOD:
        return math.Mod(x, 1)
    case CLAMP:
        if x > 1 {
            return 1
        } else if x < 0 {
            return 0
        }
    case SIGMOID:
        return sigmoid(x)
    }
    return 0
}

func (conf *Config) nFunctions() int {
    a := 0
    if conf.Alpha { a = 1 }
    switch conf.ColorSpace {
    case GRAYSCALE:
        return a + 1
    case RGB, HSV, YCbCr:
        return a + 3
    case CMYK:
        return a + 4
    }
    panic("Invalid color space!")
    return a
}

const defaultFunctionLength = 40

func GenerateImageFromFunctions(width int, height int, config Config,
                                functions []autoutils.Function,
                                vars []float64) image.Image {
    var rect = image.Rectangle{image.Point{0, 0}, image.Point{width, height}}
    img := image.NewRGBA(rect)
    colorSpace := config.ColorSpace
    alpha := config.Alpha
    rectifier := config.Rectifier
    nfunctions := len(functions)
    rets := make([]uint8, nfunctions)
    fwidth, fheight := float64(width), float64(height)
    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            switch config.CoordinateSys {
            case XY:
                vars[0], vars[1] = float64(x)/fwidth, float64(y)/fheight
            case RTHETA:
                dx, dy := float64(x - width/2), float64(y - height/2)
                vars[0] = math.Sqrt(dx * dx + dy * dy) / ((fwidth+fheight)/2) // r
                vars[1] = math.Atan2(dy, dx) // theta
            }
            for i := range rets {
                ret := rectify(functions[i].Evaluate(vars), rectifier)
                rets[i] = uint8(255 * ret)
            }
            var r, g, b, a uint8
            a = 255
            switch (colorSpace) {
            case RGB:
                r, g, b = rets[0], rets[1], rets[2]
            case GRAYSCALE:
                r, g, b = rets[0], rets[0], rets[0]
            case CMYK:
                r, g, b = color.CMYKToRGB(rets[0], rets[1], rets[2], rets[3])
            case HSV:
                r, g, b = autoutils.HSVToRGB(rets[0], rets[1], rets[2])
            case YCbCr:
                r, g, b = color.YCbCrToRGB(rets[0], rets[1], rets[2])
            }
            if (alpha) {
                a = rets[nfunctions-1]
            }
            img.Set(x, y, color.RGBA{r, g, b, a})
        }
    }
    return img
}

func GenerateImage(width int, height int, config Config) image.Image {
    if config.FunctionLength == 0 {
        // 0 value of config shouldn't have empty functions
        config.FunctionLength = defaultFunctionLength
    }

    functionLength := config.FunctionLength

    nfunctions := config.nFunctions()
    functions := make([]autoutils.Function, nfunctions)
    for i := range functions {
        functions[i].Generate(2, functionLength)
    }
    vars := []float64{0, 0}
    return GenerateImageFromFunctions(width, height, config, functions, vars)
}


func GenerateImages(width int, height int, config Config, number int, verbose bool) []image.Image {
    c := make(chan image.Image)
    for i := 0; i < number; i++ {
        go func () {
            c <- GenerateImage(width, height, config)
        }()
    }
    imgs := make([]image.Image, number)
    for i := range imgs {
        imgs[i] = <-c
		if verbose {
			fmt.Println("Generating images...", i+1, "/", number)
		}
    }
    return imgs
}

type PaletteConfig struct {
	NColors int
	Alpha bool
	FunctionLength int
    CoordinateSys int
}

func GenerateImagePaletteFrom(width int, height int, conf PaletteConfig,
        funcs []autoutils.Function, vars []float64,
		palette []color.RGBA) image.Image {
	img := image.NewRGBA(image.Rectangle{image.Point{0,0}, image.Point{width, height}})
    fwidth, fheight := float64(width), float64(height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
            switch conf.CoordinateSys {
            case XY:
                vars[0], vars[1] = float64(x)/fwidth, float64(y)/fheight
            case RTHETA:
                dx, dy := float64(x - width/2), float64(y - height/2)
                vars[0] = math.Sqrt(dx * dx + dy * dy) / ((fwidth+fheight)/2) // r
                vars[1] = math.Atan2(dy, dx) // theta
            }
			for i := range palette {
				if i == conf.NColors - 1 {
					// Background color
					img.Set(x, y, palette[i])
				} else if funcs[i].Evaluate(vars) < 0 {
					img.Set(x, y, palette[i])
					break
				}
			}
		}
	}
	return img
}

func GenerateImagePalette(width int, height int, conf PaletteConfig) image.Image {
	nColors := conf.NColors
	alpha := conf.Alpha
	functionLength := conf.FunctionLength

	funcs := make([]autoutils.Function, nColors - 1)
	palette := make([]color.RGBA, nColors)

	// Choose palette
	for i := range palette {
		r, g, b := rand.Intn(256), rand.Intn(256), rand.Intn(256)
		var a int
		if alpha {
			a = rand.Intn(256)
		} else {
			a = 255
		}
		palette[i] = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	}
	// Choose functions
	for i := range funcs {
		funcs[i].Generate(2, functionLength)
	}

	vars := make([]float64, 2)
	return GenerateImagePaletteFrom(width, height, conf, funcs, vars, palette)
}

func GenerateImagesPalette(width int, height int, conf PaletteConfig, number int, verbose bool) []image.Image {
	c := make(chan image.Image)
	for i := 0; i < number; i++ {
		go func() {
			c <- GenerateImagePalette(width, height, conf)
		}()
	}
	images := make([]image.Image, number)
	for i := 0; i < number; i++ {
		images[i] = <-c
		if verbose {
			fmt.Println("Generating images...", i+1, "/", number)
		}
	}
	return images
}
