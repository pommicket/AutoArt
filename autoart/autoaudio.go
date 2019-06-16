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
    "github.com/pommicket/autoart/autoutils"
    "io"
)

func GenerateAudio(output io.Writer, duration float64, sampleRate int32,
                   functionLength int, rectifier int) error {
    samples := int64(duration * float64(sampleRate))
    err := autoutils.WriteAudioHeader(output, samples, 1, sampleRate)
    if err != nil { return err }

    vars := make([]float64, 1)
    const sampleBufferSize = 4096
    sampleBuffer := make([]uint8, sampleBufferSize)
    sampleBufferIndex := 0

    var function autoutils.Function
    function.Generate(1, functionLength)

    for s := int64(0); s < samples; s++ {
        t := float64(s) / float64(sampleRate)
        vars[0] = t
        value := rectify(function.Evaluate(vars), rectifier)
        sampleBuffer[sampleBufferIndex] = uint8(255 * value)
        sampleBufferIndex++
        if sampleBufferIndex == sampleBufferSize {
            err = autoutils.WriteAudioSamples(output, sampleBuffer)
            if err != nil {
                return err
            }
            sampleBufferIndex = 0
        }
    }
    return autoutils.WriteAudioSamples(output, sampleBuffer[:sampleBufferIndex])
}