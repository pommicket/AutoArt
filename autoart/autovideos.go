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

/*
NOTE: AutoVideos requires Go 1.11 or newer (for temp file patterns).
*/

package autoart

import (
	"fmt"
	"github.com/pommicket/autoart/autoutils"
	"image"
	"image/color"
	"image/png"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
)

func generateFrame(width int, height int, paletted bool, config Config,
	pconfig PaletteConfig, palette []color.RGBA,
	functions []autoutils.Function, time float64,
	frameNumber int64, file *os.File) error { // NOTE: file is closed by this function
	vars := []float64{0, 0, time}
	var img image.Image
	if paletted {
		img = GenerateImagePaletteFrom(width, height, pconfig, functions, vars, palette)
	} else {
		img = GenerateImageFromFunctions(width, height, config, functions, vars)
	}
	if err := png.Encode(file, img); err != nil {
		file.Close()
		return err
	}
	return file.Close()
}

func generateVideo(width int, height int, paletted bool, config Config,
	pconfig PaletteConfig, time float64,
	framerate int, filename string, verbose bool) error {

	var palette []color.RGBA
	if paletted {
		// Generate palette
		palette = make([]color.RGBA, pconfig.NColors)
		for i := range palette {
			r := uint8(rand.Intn(256))
			g := uint8(rand.Intn(256))
			b := uint8(rand.Intn(256))
			a := uint8(255)
			if pconfig.Alpha {
				a = uint8(rand.Intn(256))
			}
			palette[i] = color.RGBA{r, g, b, a}
		}
	}

	if config.FunctionLength == 0 {
		// 0 value of config shouldn't have empty functions
		config.FunctionLength = defaultFunctionLength
	}

	var functionLength int
	if paletted {
		functionLength = pconfig.FunctionLength
	} else {
		functionLength = config.FunctionLength
	}

	var nfunctions int
	if paletted {
		nfunctions = pconfig.NColors
	} else {
		nfunctions = config.nFunctions()
	}
	functions := make([]autoutils.Function, nfunctions)
	for i := range functions {
		functions[i].Generate(3, functionLength)
	}

	frames := int64(time * float64(framerate))

	files := make([]*os.File, frames)
	defer func() {
		// Delete all temporary files
		for _, file := range files {
			if file != nil {
				os.Remove(file.Name())
			}
		}
	}()
	// Create temporary frame files
	for i := range files {
		var err error
		files[i], err = ioutil.TempFile("", "frame*.png")
		if err != nil {
			return err
		}
	}

	autoutils.RunInBatches(frames, "Generating video...", func(n int64, errs chan<- error) {
		t := float64(n) / float64(framerate)
		errs <- generateFrame(width, height, paletted, config, pconfig, palette, functions, t, n, files[n])
	})

	ffmpegInputFile, err := ioutil.TempFile("", "input*.txt")
	if err != nil {
		return err
	}
	defer os.Remove(ffmpegInputFile.Name())

	for _, file := range files {
		name := file.Name()
		info := fmt.Sprintf("file '%v'\nduration %v\n", name, 1/float64(framerate))
		if _, err = io.WriteString(ffmpegInputFile, info); err != nil {
			return err
		}
	}

	fmt.Println(ffmpegInputFile.Name())

	if err = ffmpegInputFile.Close(); err != nil {
		return err
	}

	if verbose {
		fmt.Println("ffmpeg", "-y", "-f", "concat", "-safe", "0", "-i", ffmpegInputFile.Name(), filename)
	}

	cmd := exec.Command("ffmpeg", "-y", "-f", "concat", "-safe", "0", "-i", ffmpegInputFile.Name(), filename)
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg failed: %v", err)
	}

	return nil
}

func GenerateVideo(width int, height int, config Config, time float64,
	framerate int, filename string, verbose bool) error {
	var pconfig PaletteConfig
	return generateVideo(width, height, false, config, pconfig, time, framerate, filename, verbose)
}

func GenerateVideoPalette(width int, height int, pconfig PaletteConfig,
	time float64, framerate int, filename string, verbose bool) error {
	var config Config
	return generateVideo(width, height, true, config, pconfig, time, framerate, filename, verbose)
}
