package pipeline

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"sort"
	"time"
)

var (
	startTime time.Time
)

func Init() {
	startTime = time.Now()
}

func ArraySource(a ...int) <-chan int {
	out := make(chan int, 1024)
	go func() {
		for _, v := range a {
			out <- v
		}
		close(out)
	}()
	return out
}

func InMemSort(in <-chan int) <-chan int {
	out := make(chan int, 1024)
	go func() {
		temp := []int{}
		for v := range in {
			temp = append(temp, v)
		}
		fmt.Println("Read Done:", time.Now().Sub(startTime))

		sort.Ints(temp)
		fmt.Println("Sort Done:", time.Now().Sub(startTime))

		for _, v := range temp {
			out <- v
		}
		close(out)
	}()
	return out
}

func Merge(in1, in2 <-chan int) <-chan int {
	out := make(chan int, 1024)
	go func() {
		v1, ok1 := <-in1
		v2, ok2 := <-in2
		for ok1 || ok2 {
			if !ok2 || (ok1 && v1 <= v2) {
				out <- v1
				v1, ok1 = <-in1
			} else {
				out <- v2
				v2, ok2 = <-in2
			}
		}
		close(out)
		fmt.Println("Merge Done:", time.Now().Sub(startTime))

	}()
	return out
}

func ReadSource(reader io.Reader, chunkSize int) <-chan int {
	out := make(chan int, 1024)
	go func() {
		buffer := make([]byte, 8)
		byteRead := 0
		for {
			n, err := reader.Read(buffer)
			byteRead += n
			if n > 0 {
				out <- int(binary.BigEndian.Uint64(buffer))
			}
			if err != nil || (chunkSize != -1 && byteRead >= chunkSize) {
				break
			}
		}
		close(out)
	}()
	return out
}

func WriteSink(writer io.Writer, in <-chan int) {
	buffer := make([]byte, 8)
	for v := range in {
		binary.BigEndian.PutUint64(buffer, uint64(v))
		writer.Write(buffer)
	}
}

func RandomSource(count int) <-chan int {
	out := make(chan int, 1024)
	go func() {
		for i := 0; i < count; i++ {
			out <- rand.Intn(10000)
		}
		close(out)
	}()
	return out
}

func MergeN(inputs ...<-chan int) <-chan int {
	if len(inputs) == 1 {
		return inputs[0]
	}
	m := len(inputs) / 2
	return Merge(MergeN(inputs[:m]...), MergeN(inputs[m:]...))
}
