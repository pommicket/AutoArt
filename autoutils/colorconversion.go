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

package autoutils

import (
	"math"
)

func HSVToRGB(h uint8, s uint8, v uint8) (uint8, uint8, uint8) {
	// https://en.wikipedia.org/wiki/HSL_and_HSV#HSV_to_RGB
	V := float64(v) / 256
	S := float64(s) / 256
	C := V * S
	H := float64(h) / (256 / 6)
	X := float64(C) * (1 - math.Abs(math.Mod(H, 2)-1))
	var r, g, b float64
	switch true {
	case s == 0:
		r, g, b = 0, 0, 0
	case H <= 1:
		r, g, b = C, X, 0
	case H <= 2:
		r, g, b = X, C, 0
	case H <= 3:
		r, g, b = 0, C, X
	case H <= 4:
		r, g, b = 0, X, C
	case H <= 5:
		r, g, b = X, 0, C
	default:
		r, g, b = C, 0, X
	}
	m := V - C
	r += m
	g += m
	b += m
	r *= 255
	g *= 255
	b *= 255
	return uint8(r), uint8(g), uint8(b)
}
