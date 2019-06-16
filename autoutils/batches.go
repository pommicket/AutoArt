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
    "fmt"
)

/*
This function runs f in batches of batchSize. n will be the number in the
sequence. f should send an error to errs if it wants this function to return
that error, and otherwise should send nil when it is done. Before each batch,
a message will be printed, starting with progress, and showing how many batches
have been completed so far out of the total number of batches.
*/
const batchSize = 32
func RunInBatches(number int64, progress string, f func (n int64, errs chan<- error)) error {
    nBatches := number / batchSize
    errs := make(chan error)
    for batch := int64(0); batch <= nBatches; batch++ {
        fmt.Println(progress, batch+1, "/", nBatches+1)
        thisBatchSize := batchSize
        if batch == nBatches {
            // Deal with case of last batch
            thisBatchSize = int(number - nBatches * batchSize)
        }
        for task := 0; task < thisBatchSize; task++ {
            go f(int64(task) + batchSize * batch, errs)
        }

        for completed := 0; completed < thisBatchSize; completed++ {
            err := <-errs
            if err != nil {
                return err
            }
        }
    }
    return nil
}