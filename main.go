package main

import(
	"./pipeline"
	"os"
	"fmt"
)

const (
	filename = "large"
	count = 100000000
	outputFilename = "Sorted_large"
)

func main() {
	pipeline.GenerateRandomFile(filename, count)
	fileInfo, _ := os.Stat(filename)
	filesize:= fileInfo.Size()
	fmt.Printf("Filename: %s, Size: %d\n", filename, filesize)

	p := pipeline.CreatePipeline(filename, 8 * count, 4)
	pipeline.WriteToFile(p, outputFilename)

	fileInfo, _ = os.Stat(outputFilename)
	filesize = fileInfo.Size()
	fmt.Printf("Filename: %s, Size: %d\n", outputFilename, filesize)
	// pipeline.PrintFile(outputFilename)
}