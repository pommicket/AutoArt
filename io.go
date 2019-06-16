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

// Some functions for dealing with user input

import (
    "bufio"
    "strconv"
    "strings"
    "fmt"
	"github.com/pommicket/autoart/autoart"
)


// Reads an int64 from a buffered reader, after giving the user the given prompt
// If the number that the user entered does not satisfy valid, the user will
// be prompted again. If the user enters nothing, the default will be used.
func readInt64(reader *bufio.Reader, prompt string,
               valid func (int64) bool, def int64) (int64, error) {
    fmt.Print(prompt)
    for {
        line, err := reader.ReadString('\n')
        if err != nil {
            return 0, err
        }
        line = strings.TrimSpace(line)
        if line == "" {
            return def, nil
        }
        num, err := strconv.ParseInt(line, 0, 64)
        if err != nil {
            fmt.Println("Please enter a number.")
            fmt.Print(prompt)
            continue
        }
        if !valid(num) {
            fmt.Println("Please enter a valid option.")
            fmt.Print(prompt)
            continue
        }
        return num, nil
    }
}

// Reads a bool from the reader, prompting the user until they enter y/n.
func readBool(reader *bufio.Reader, prompt string, def bool) (bool, error) {
	for {
        fmt.Print(prompt)
        line, err := reader.ReadString('\n')
        if err != nil { return false, err }
        line = strings.ToLower(strings.TrimSpace(line))
        if line == "" {
            return def, nil
        }
        switch line[0] {
        case 'y':
            return true, nil
        case 'n':
            return false, nil
    	}
		fmt.Println("Please enter yes or no.")
    }
}

// Reads an autoart.Config
func readConf(reader *bufio.Reader, conf *autoart.Config) error {
	positive := func (i int64) bool { return i > 0 }
	functionLength, err := readInt64(reader, "Function length (default: 40)? ", positive, 40)
    if err != nil { return err }
    colorSpace, err := readInt64(reader, `Which color space should be used?
1. RGB
2. Grayscale
3. CMYK
4. HSV
5. YCbCr
Please enter a number between 1 and 5 (default: 1): `, func (i int64) bool {
        return i >= 1 && i <= 5
    }, 1)

    if err != nil { return err }
    alpha, err := readBool(reader, "Should an alpha channel be included (yes/no, default: no)? ", false)
	if err != nil { return err }


    rectifier, err := readInt64(reader, `How should out of range values be dealt with?
1. Modulo
2. Clamp
3. Sigmoid
Please enter 1, 2, or 3 (default: 1): `, func (i int64) bool {
        return i >= 1 && i <= 3
    }, 1)
	if err != nil { return err }

    coords, err := readInt64(reader, `Which coordinate system should be used?
1. x, y
2. r, theta
Please enter 1 or 2 (default: 1): `, func (i int64) bool {
        return i >= 1 && i <= 2
    }, 1)
	if err != nil { return err }

	conf.FunctionLength = int(functionLength)
	conf.ColorSpace = int(colorSpace - 1)
	conf.Alpha = alpha
	conf.Rectifier = int(rectifier - 1)
    conf.CoordinateSys = int(coords - 1)
	return nil
}

func readPaletteConf(reader *bufio.Reader, conf *autoart.PaletteConfig) error {
	positive := func (i int64) bool { return i > 0 }
	ncolors, err := readInt64(reader, "How many colors do you want (default: 10)? ", positive, 10)
	if err != nil { return err }
	alpha, err := readBool(reader, "Should an alpha channel be included (yes/no, default: no)? ", false)
	if err != nil { return err }
	functionLength, err := readInt64(reader, "Function length (default: 40)? ", positive, 40)
	if err != nil { return err }
    coords, err := readInt64(reader, `Which coordinate system should be used?
1. x, y
2. r, theta
Please enter 1 or 2 (default: 1): `, func (i int64) bool {
        return i >= 1 && i <= 2
    }, 1)
	if err != nil { return err }
	conf.NColors = int(ncolors)
	conf.Alpha = alpha
	conf.FunctionLength = int(functionLength)
    conf.CoordinateSys = int(coords - 1)
	return nil
}
