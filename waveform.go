package waveform

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
)

func Generate(sourcePath string, tmpPath string, widthCount float64) string {
	sourceFilename, err := filepath.Abs(sourcePath)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	filename := filepath.Base(sourcePath)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	tempFileName := fmt.Sprintf("%s/%s.raw", tmpPath, filename)
	generateRawFile(sourceFilename, tempFileName)

	rawFile, err := os.Open(tempFileName)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	minimumValues, maximumValues := extractMinMaxValues(sourcePath, rawFile, widthCount)
	os.Remove(tempFileName)
	percents := convertToPercentage(minimumValues, maximumValues)

	result, err := json.Marshal(percents)

	if err != nil {
		log.Fatal(err)
		return ""
	}

	return string(result)
}

func convertToPercentage(minimumValues []int64, maximumValues []int64) []float64 {
	width := len(maximumValues)
	heightsInInt64 := make([]int64, width)
	heights := make([]float64, width)
	highestHeight := maximumValues[0] - minimumValues[0]
	heightsInInt64[0] = 0
	for i := 1; i < width; i++ {
		heightsInInt64[i] = maximumValues[i] - minimumValues[i]
		if highestHeight < heightsInInt64[i] {
			highestHeight = heightsInInt64[i]
		}
	}

	highestHeightInFloat64 := float64(highestHeight)

	for i := 0; i < width; i++ {
		heights[i] = toFixed(float64(heightsInInt64[i])/highestHeightInFloat64, 2)
	}
	return heights
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
