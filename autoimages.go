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

package main

import (
	"bufio"
	"fmt"
	"github.com/pommicket/autoart/autoart"
	"github.com/pommicket/autoart/autoutils"
	"image"
	"image/png"
	"math/rand"
	"os"
	"time"
)

// AutoImages client

func genImage(width int, height int, paletted bool, conf *autoart.Config, pconf *autoart.PaletteConfig, filename string) error {
	var img image.Image
	if paletted {
		img = autoart.GenerateImagePalette(width, height, *pconf)
	} else {
		img = autoart.GenerateImage(width, height, *conf)
	}
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	err = png.Encode(file, img)
	if err != nil {
		file.Close()
		return err
	}
	return file.Close()
}

func batchedImages(seed int64, width int, height int, paletted bool, conf *autoart.Config, pconf *autoart.PaletteConfig, number int64) error {
	// Create a directory for the images
	rand.Seed(seed)
	dir := fmt.Sprintf("autoimages%v", seed)
	err := os.MkdirAll(dir, 0700)
	if err != nil {
		return err
	}
	err = autoutils.RunInBatches(number, "Generating images...", func(n int64, errs chan<- error) {
		filename := fmt.Sprintf("%v/%09d.png", dir, n)
		errs <- genImage(width, height, paletted, conf, pconf, filename)
	})

	if err != nil {
		return err
	}
	fmt.Println("Done! Your images are in this directory:", dir)
	return nil
}

func autoImages(reader *bufio.Reader) error {
	prompt := `How many options do you want?
1. None - Just make an image
2. Some - Basic options
3. All  - Advanced options
Please enter 1, 2, or 3 (default: 1): `
	option, err := readInt64(reader, prompt, func(i int64) bool {
		return i >= 1 && i <= 3
	}, 1)
	if err != nil {
		return err
	}
	var conf autoart.Config
	var pconf autoart.PaletteConfig
	t := time.Now().UTC().UnixNano()
	if option == 1 {
		rand.Seed(t)
		fmt.Println("Generating image...")
		filename := fmt.Sprintf("autoimages%d.png", t)
		err = genImage(1920, 1080, false, &conf, &pconf, filename)
		if err != nil {
			// We're done!
			fmt.Println("Generated an image:", filename)
		}
		return err
	}
	// Basic options
	positive := func(i int64) bool { return i > 0 }
	width, err := readInt64(reader, "Width (default: 1920)? ", positive, 1920)
	if err != nil {
		return err
	}
	height, err := readInt64(reader, "Height (default: 1080)? ", positive, 1080)
	if err != nil {
		return err
	}
	number, err := readInt64(reader, "How many (default: 1)? ", positive, 1)
	if err != nil {
		return err
	}
	if option == 2 {
		return batchedImages(t, int(width), int(height), false, &conf, &pconf, number)
	}

	paletted, err := readBool(reader, "Should a palette be used (y/n, default: n)? ", false)
	if err != nil {
		return err
	}

	// Advanced options
	if paletted {
		err = readPaletteConf(reader, &pconf)
		if err != nil {
			return err
		}
	} else {
		err = readConf(reader, &conf)
		if err != nil {
			return err
		}
	}
	seed, err := readInt64(reader, "Random seed (default: current time)? ", func(i int64) bool {
		return true
	}, t)

	return batchedImages(seed, int(width), int(height), paletted, &conf, &pconf, number)
}
