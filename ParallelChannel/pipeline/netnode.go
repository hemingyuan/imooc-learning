package pipeline

import (
	"bufio"
	"net"
)

func NetSource(addr string) <-chan int {
	out := make(chan int, 1024)
	go func() {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			panic(err)
		}
		// defer conn.Close()
		r := ReadSource(bufio.NewReader(conn), -1)
		// conn.
		for v := range r {
			out <- v
		}
		close(out)
	}()
	return out
}

func NetSink(addr string, in <-chan int) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	go func() {
		c, err := l.Accept()
		if err != nil {
			panic(err)
		}
		defer c.Close()

		write := bufio.NewWriter(c)
		// write.Write()
		WriteSink(write, in)
	}()
}
