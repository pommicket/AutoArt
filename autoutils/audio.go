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
	"encoding/binary"
	"io"
)

// Write a header to writer. You need to decide ahead of time how many samples
// you want.
func WriteAudioHeader(writer io.Writer, nSamples int64, channels, sampleRate int32) error {
	w := func(data interface{}) error {
		return binary.Write(writer, binary.LittleEndian, data)
	}
	var err error
	if err = w([]byte("RIFF")); err != nil {
		return err
	}
	var chunkSize1 uint32 = 36 + uint32(nSamples)
	if err = w(chunkSize1); err != nil {
		return err
	}
	if err = w([]byte("WAVEfmt ")); err != nil {
		return err
	}
	var subchunk1size uint32 = 16
	if err = w(subchunk1size); err != nil {
		return err
	}
	var audioFormat uint16 = 1
	if err = w(audioFormat); err != nil {
		return err
	}
	var nChannels uint16 = uint16(channels)
	if err = w(nChannels); err != nil {
		return err
	}
	var srate uint32 = uint32(sampleRate)
	if err = w(srate); err != nil {
		return err
	}
	var byteRate uint32 = srate * uint32(nChannels)
	if err = w(byteRate); err != nil {
		return err
	}
	var blockAlign uint16 = nChannels
	if err = w(blockAlign); err != nil {
		return err
	}
	var bitsPerSample uint16 = 8
	if err = w(bitsPerSample); err != nil {
		return err
	}
	if err = w([]byte("data")); err != nil {
		return err
	}
	var chunkSize2 uint32 = uint32(nSamples) * uint32(nChannels)
	if err = w(chunkSize2); err != nil {
		return err
	}
	return nil
}

// Writes some samples to the writer. You will need to write a header before
// any samples.
func WriteAudioSamples(writer io.Writer, samples []uint8) error {
	return binary.Write(writer, binary.LittleEndian, samples)
}

/*
Writes audio data in WAV format to the writer. There is only support for 8-bit
audio. If there are multiple channels, audio[0] should refer to the first sample
for the first channel, audio[1] should refer to the first sample for the second
channel, etc.
WAV does not support more then 65535 channels (but what are you doing if
you're using that many?!)
*/
func WriteAudio(writer io.Writer, audio []uint8, channels, sampleRate int32) error {

	err := WriteAudioHeader(writer, int64(len(audio)), channels, sampleRate)
	if err != nil {
		return err
	}
	if err = WriteAudioSamples(writer, audio); err != nil {
		return err
	}
	return nil
}
