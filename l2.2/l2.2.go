package main

import "fmt"

func test() (x int) {
	defer func() { // первый
		x++ // 1 + 1 = 2
	}()
	x = 1 // = 1
	return
}

func anotherTest() int {
	var x int      // = 0
	defer func() { // второй
		x++ // = 2 но локально (не идет в ретёрн)
	}()
	x = 1    // = 1
	return x // 1 потому что он в ретёрн копируется и не важно что x меняется копия была ДО
}

func main() {
	fmt.Println(test())        // 2
	fmt.Println(anotherTest()) // 1
}
