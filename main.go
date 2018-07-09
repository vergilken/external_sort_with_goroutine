package main

import(
	"./pipeline"
)

const (
	filename = "large"
	count = 100000000
	outputFilename = "Sorted_large"
)

func main() {
	pipeline.GenerateRandomFile(filename, count)
	p := pipeline.CreatePipeline(filename, 8 * count, 4)
	pipeline.WriteToFile(p, outputFilename)
	// pipeline.PrintFile(outputFilename)
}