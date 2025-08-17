package main

import "fmt"

func main() {
	a := [6]int{76, 77, 78, 79, 80}
	var b []int = a[1:4] // в b передается срез из массива a от 1 до 4 не включительно
	fmt.Println(b)       // вывод {77 78 79}
}
