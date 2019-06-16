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
	"math/rand"
	"os"
	"time"
)

func generateAudio(seed int64, length int64, sampleRate int64, functionLength int64, number int64) error {
	rand.Seed(seed)
	dir := fmt.Sprintf("autoaudio%v", seed)
	err := os.MkdirAll(dir, 0700)
	if err != nil {
		return err
	}
	err = autoutils.RunInBatches(number, "Generating audio...", func(n int64, errs chan<- error) {
		filename := fmt.Sprintf("%v/%09d.wav", dir, n)
		file, err := os.Create(filename)
		if err != nil {
			errs <- err
			return
		}
		err = autoart.GenerateAudio(file, float64(length), int32(sampleRate), int(functionLength), autoart.MOD)
		if err != nil {
			errs <- err
			return
		}
		errs <- file.Close()
	})
	if err != nil {
		return err
	}
	fmt.Println("Done. Your audio is in this directory:", dir)
	return nil
}

func autoAudio(reader *bufio.Reader) error {
	prompt := `How many options do you want?
1. None - Just make some audio
2. Some - Basic options
3. All  - Advanced options
Please enter 1, 2, or 3 (default: 1): `
	option, err := readInt64(reader, prompt, func(i int64) bool {
		return i >= 1 && i <= 3
	}, 1)
	if err != nil {
		return err
	}
	t := time.Now().UTC().UnixNano()
	if option == 1 {
		filename := fmt.Sprintf("autoaudio%v.wav", t)
		rand.Seed(t)
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		err = autoart.GenerateAudio(file, 60, 44100, 80, autoart.MOD)
		if err != nil {
			return err
		}
		fmt.Println("Generated audio:", filename)
		return nil
	}
	positive := func(i int64) bool { return i > 0 }
	length, err := readInt64(reader, "Length in seconds (default: 60)? ", positive, 60)
	if err != nil {
		return err
	}
	number, err := readInt64(reader, "Number (default: 1)? ", positive, 1)
	if err != nil {
		return err
	}
	if option == 2 {
		return generateAudio(t, length, 44100, 80, number)
	}
	sampleRate, err := readInt64(reader, "Sample rate (default: 44100)? ", positive, 44100)
	if err != nil {
		return err
	}
	functionLength, err := readInt64(reader, "Function length (default: 80)? ", positive, 80)
	if err != nil {
		return err
	}
	seed, err := readInt64(reader, "Random seed (default: current time)? ", positive, t)
	if err != nil {
		return err
	}
	return generateAudio(seed, length, sampleRate, functionLength, number)

}
