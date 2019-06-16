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
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	prompt := `Please select one of the following:
1. AutoImages
2. AutoVideos
3. AutoAudio
Please enter 1, 2, or 3 (default: 1): `

	option, err := readInt64(reader, prompt, func(i int64) bool {
		return i >= 1 && i <= 3
	}, 1)
	if err != nil {
		fmt.Println("Error reading user input:", err)
	}

	switch option {
	case 1:
		err = autoImages(reader)
	case 2:
		err = autoVideos(reader)
	case 3:
		err = autoAudio(reader)
	}

	if err != nil {
		fmt.Println("An error occured:", err)
	}

}
