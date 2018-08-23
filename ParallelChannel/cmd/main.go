package main

import (
	"bufio"
	"fmt"
	"imooc_learning/ParallelChannel/pipeline"
	"os"
)

func main() {
	// demo()
	p := createPipeline("large.in", 80000000, 8)
	writeToFile(p, "large.out")
}

func readFile(filename string) {
	f, _ := os.Open(filename)
	p := pipeline.ReadSource(f, -1)
	for v := range p {
		fmt.Println(v)
	}
}

func writeToFile(p <-chan int, filename string) {
	f, _ := os.Create(filename)
	pipeline.WriteSink(f, p)
}

func createPipeline(filename string, filesize, chunkCount int) <-chan int {
	readSize := filesize / chunkCount
	result := []<-chan int{}
	pipeline.Init()
	for i := 0; i < chunkCount; i++ {
		f, _ := os.Open(filename)
		f.Seek(int64(i*readSize), 0)
		p := pipeline.ReadSource(f, readSize)
		result = append(result, pipeline.InMemSort(p))
	}
	return pipeline.MergeN(result...)
}

func createNetPipeline(filename string, filesize, chunkCount int) <-chan int {
	readSize := filesize / chunkCount

	sortAddr := []string{}
	pipeline.Init()
	for i := 0; i < chunkCount; i++ {
		f, _ := os.Open(filename)
		f.Seek(int64(i*readSize), 0)
		p := pipeline.ReadSource(bufio.NewReader(f), readSize)
		addr := fmt.Sprintf(":%d", 7000+i)
		pipeline.NetSink(addr, pipeline.InMemSort(p))

		sortAddr = append(sortAddr, addr)
		// result = append(result, pipeline.InMemSort(p))
	}
	result := []<-chan int{}

	for _, v := range sortAddr {
		result = append(result, pipeline.NetSource(v))
	}
	return pipeline.MergeN(result...)
}

func demo() {
	const (
		file = "large.in"
		n    = 10000000
	)
	f, _ := os.Create(file)
	w := bufio.NewWriter(f)
	defer f.Close()
	p := pipeline.RandomSource(n)

	pipeline.WriteSink(w, p)
	w.Flush()

	// f, _ := os.Open(file)
	// defer f.Close()

	// p1 := pipeline.ReadSource(f, -1)
	// for v := range p1 {
	// 	fmt.Println(v)
	// }
}

func megerDemo() {
	p1 := pipeline.InMemSort(
		pipeline.ArraySource(3, 5, 7, 6, 9, 1, 2),
	)
	p2 := pipeline.InMemSort(
		pipeline.ArraySource(7, 0, 9, 4, 2, 7, 8),
	)

	p := pipeline.Merge(p1, p2)
	for k := range p {
		fmt.Println(k)
	}
}
