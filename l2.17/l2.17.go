package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	if flag.NArg() < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s host port [--timeout=10s]\n", os.Args[0])
		os.Exit(1)
	}

	host, port := flag.Arg(0), flag.Arg(1)
	address := net.JoinHostPort(host, port)

	conn, err := net.DialTimeout("tcp", address, *timeout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Printf("connected to %s\n", address)

	done := make(chan struct{})

	// горутина для чтения из сокета
	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		fmt.Println("connection closed by server")
		close(done)
	}()

	// горутина для отправки пользовательского ввода
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			text := scanner.Text() + "\n"
			_, err := conn.Write([]byte(text))
			if err != nil {
				fmt.Println("write error:", err)
				break
			}
		}
		// Ctrl+D (EOF) → выходим
		fmt.Println("closing connection...")
		conn.Close()
		close(done)
	}()

	<-done
}
