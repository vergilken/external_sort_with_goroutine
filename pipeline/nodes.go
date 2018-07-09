package pipeline

import (
	"sort"
	"io"
	"encoding/binary"
	"math/rand"
	"os"
	"bufio"
	"fmt"
	"time"
)

var startTime time.Time
func Init() {
	startTime = time.Now()
}

func ArraySource(a ...int) <-chan int{
	out := make(chan int)
	go func() {
		for _, v := range a {
			out <- v
		}
		close(out)
	} ()
	return out
}

// In memory sort by quick sorting
func InMemSort(in <-chan int) chan int {
	out := make(chan int, 1024)
	go func() {
		a := []int{}
		for v := range in {
			a = append(a, v)
		}
		fmt.Println("Read Done: ", time.Now().Sub(startTime))

		sort.Ints(a)
		fmt.Println("In Memory Sort Done: ", time.Now().Sub(startTime))

		for _, v := range a {
			out <- v
		}
		close(out)
	}()

	return out
}

// Merge two channel
func Merge(in1, in2 <-chan int) chan int {
	out := make(chan int, 1024)
	go func() {
		v1, ok1 := <- in1
		v2, ok2 := <- in2
		for ok1 || ok2 {
			if !ok2 || (ok1 && v1 <= v2) {
				out <- v1
				v1, ok1 = <- in1
			} else {
				out <- v2
				v2, ok2 = <- in2
			}
		}
		close(out)
		fmt.Println("Merge Done: ", time.Now().Sub(startTime))

	}()

	return out
}


// TODO: corner case on imcompleted chunk
func ReaderSource(reader io.Reader, chunkSize int) <-chan int {
	out := make(chan int, 1024)
	go func() {
		buffer := make([]byte, 8)
		bytesRead := 0
		for {
			n, err := reader.Read(buffer)
			bytesRead += n
			if n > 0 {
				v := int(binary.BigEndian.Uint64(buffer))
				out <- v
			}
			if err != nil || (chunkSize != -1 && bytesRead >= chunkSize) {
				break
			}
		}
		close(out)
		}()

	return out
}

func WriterSink(writer io.Writer, in <-chan int) {
	for v := range in {
		buffer := make([]byte, 8)
		binary.BigEndian.PutUint64(buffer, uint64(v))
		writer.Write(buffer)
	}
}

func RandomSource(count int) <-chan int {
	out := make(chan int, 1024)
	go func() {
		for i := 0; i < count; i++ {
			out <- rand.Int()
		}
		close(out)
	}()

	return out
}

func MergeN(inputs ...<-chan int) <-chan int {
	m := len(inputs)
	if m == 1 {
		return inputs[0]
	}
	m /= 2
	return Merge(MergeN(inputs[:m]...), MergeN(inputs[m:]...))
}

// Outside function
func GenerateRandomFile (filename string, count int) {
	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	p := RandomSource(count)
	writer := bufio.NewWriter(file)
	WriterSink(writer, p)
	writer.Flush()
}

func CreatePipeline(filename string,
					fileSize, chunkCount int) <-chan int{
	chunkSize := fileSize / chunkCount
	sortResults := []<-chan int{}
	Init()

	for i := 0; i < chunkCount; i++ {
		file, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		file.Seek(int64(i * chunkSize), 0)

		source := ReaderSource(bufio.NewReader(file), chunkSize)

		sortResults = append(sortResults, InMemSort(source))
	}

	return MergeN(sortResults...)
}

func WriteToFile(p <-chan int, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	WriterSink(writer, p)
}

func PrintFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	count := 0
	p := ReaderSource(file, -1)
	for v := range p {
		fmt.Println(v)
		if count++; count >= 100 {
			break
		}
	}
}