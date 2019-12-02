package waveform

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

const (
	BUFFER_SIZE           = 300
	NUMBER_OF_BYTES       = 4
	MAX_AUDIO_VALUE int64 = 65536
	MIN_AUDIO_VALUE int64 = -65536
)

func extractMinMaxValues(sourcePath string, rawFile *os.File, widthCount float64) ([]int64, []int64) {
	rawfileInfo, err := rawFile.Stat()
	if err != nil {
		log.Fatal(err)
	}

	widthInFloat := widthCount
	segmentSize := int(float64(rawfileInfo.Size())/widthInFloat+0.5) / NUMBER_OF_BYTES
	width := int(widthInFloat)
	maximumValues := make([]int64, width)
	minimumValues := make([]int64, width)
	segmentByteSize := segmentSize * NUMBER_OF_BYTES

	var wg sync.WaitGroup
	c := make(chan int, width)

	for position := 0; position < width; position++ {
		if position%BUFFER_SIZE == 0 {
			wg.Add(1)
			go func(index int) {
				startTime := time.Now().Local()
				fmt.Printf("%d: Started\n", index)
				lastIndex := index + BUFFER_SIZE
				if lastIndex > width {
					lastIndex = width
				}

				mins, maxs := getMinMaxValuesWithIndexFromFile(rawFile, index, lastIndex, segmentByteSize)
				for i := 0; i < len(mins); i++ {
					minimumValues[index+i] = mins[i]
					maximumValues[index+i] = maxs[i]
				}
				fmt.Printf("%d: [%v]\n", index, time.Now().Sub(startTime))
				wg.Done()
			}(position)
		}
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	for _ = range c {
	}
	return minimumValues, maximumValues
}

func getMinMaxValuesWithIndexFromFile(file *os.File, startIndex int, lastIndex int, segmentByteSize int) ([]int64, []int64) {
	numberOfSegments := lastIndex - startIndex
	data := make([]byte, segmentByteSize*numberOfSegments)
	mins := make([]int64, numberOfSegments)
	maxs := make([]int64, numberOfSegments)
	n, err := file.ReadAt(data, int64(startIndex*segmentByteSize))
	if err != nil {
		if err == io.EOF {
			return nil, nil
		}
		log.Fatal(err)
	}

	var start int
	var last int
	for i := 0; i < numberOfSegments; i++ {
		start = i * segmentByteSize
		last = (i + 1) * segmentByteSize

		if last > n {
			last = n
		}

		min, max := getMinMaxValue(data[start:last], last-start)
		mins[i] = min
		maxs[i] = max
	}

	return mins, maxs
}

func getMinMaxValue(data []byte, dataLength int) (int64, int64) {
	max := MIN_AUDIO_VALUE
	min := MAX_AUDIO_VALUE

	word := make([]byte, NUMBER_OF_BYTES*2)
	for index := 0; index < dataLength; index++ {
		word[index%NUMBER_OF_BYTES] = data[index]
		if (index+1)%NUMBER_OF_BYTES == 0 {
			for j := 0; j < NUMBER_OF_BYTES; j++ {
				word[NUMBER_OF_BYTES+j] = 0
			}

			var value int32
			var valueInInt64 int64
			buf := bytes.NewReader(word)
			err := binary.Read(buf, binary.LittleEndian, &value)

			valueInInt64 = int64(value)
			if err != nil {
				log.Fatal(err)
			}

			if valueInInt64 < min {
				min = valueInInt64
			}

			if valueInInt64 > max {
				max = valueInInt64
			}
		}
	}
	return min, max
}
