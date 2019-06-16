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
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

func autoVideos(reader *bufio.Reader) error {
	// Check if the user has ffmpeg
	cmd := exec.Command("ffmpeg", "-version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Is ffmpeg installed? (%v)", err)
	}

	tmp, err := ioutil.TempFile("", "frame*.png")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	tmp.Close()
	os.Remove(tmpName)
	fmt.Println("Warning: If you stop AutoVideos while it's running, there will be some temporary files, which might take up quite a bit of space. These will be deleted if you just let it finish running, and they will probably be deleted when you reboot your computer.")

	fmt.Println("The file names will look something like this: ")
	fmt.Println(tmpName)
	prompt := `How many options do you want?
1. None - Just make a video
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
	t := time.Now().UTC().UnixNano()
	if option == 1 {
		rand.Seed(t)
		filename := fmt.Sprintf("autovideos%v.mp4", t)
		err := autoart.GenerateVideo(1440, 900, conf, 10, 24, filename, true)
		fmt.Println("Generated video:", filename)
		return err
	}
	positive := func(i int64) bool { return i > 0 }
	width, err := readInt64(reader, "Width (default: 1440)? ", positive, 1440)
	if err != nil {
		return err
	}
	height, err := readInt64(reader, "Height (default: 900)? ", positive, 900)
	if err != nil {
		return err
	}
	length, err := readInt64(reader, "Length in seconds (default: 10)? ", positive, 10)
	if err != nil {
		return err
	}
	number, err := readInt64(reader, "Number (default: 1)? ", positive, 1)
	if err != nil {
		return err
	}
	if option == 2 {
		dir := fmt.Sprintf("autovideos%v", t)
		err = os.MkdirAll(dir, 0700)
		if err != nil {
			return err
		}
		rand.Seed(t)
		for i := int64(0); i < number; i++ {
			err = autoart.GenerateVideo(int(width), int(height), conf,
				float64(length), 24,
				fmt.Sprintf("%v/%09d.mp4", dir, i), true)
			if err != nil {
				return err
			}
		}
		fmt.Println("Done. Your videos are in this directory:", dir)
		return nil
	}

	framerate, err := readInt64(reader, "Frame rate (default: 24)? ", positive, 24)

	var pconf autoart.PaletteConfig
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

	dir := fmt.Sprintf("autovideos%v", seed)
	err = os.MkdirAll(dir, 0700)
	if err != nil {
		return err
	}
	rand.Seed(seed)
	if paletted {
		for i := int64(0); i < number; i++ {
			err = autoart.GenerateVideoPalette(int(width), int(height), pconf,
				float64(length), int(framerate),
				fmt.Sprintf("%v/%09d.mp4", dir, i), true)
			if err != nil {
				return err
			}
		}
	} else {
		for i := int64(0); i < number; i++ {
			err = autoart.GenerateVideo(int(width), int(height), conf,
				float64(length), int(framerate),
				fmt.Sprintf("%v/%09d.mp4", dir, i), true)
			if err != nil {
				return err
			}
		}
	}
	fmt.Println("Done. Your videos are in this directory:", dir)
	return nil

}
