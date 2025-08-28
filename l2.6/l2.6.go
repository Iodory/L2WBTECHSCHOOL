package main

import (
	"fmt"
)

func main() {
	var s = []string{"1", "2", "3"} // 1 2 3
	modifySlice(s)
	fmt.Println(s)              // 3 2 3 (ВЫВОДИМ ОРИГИНАЛЬНЫЙ СЛАЙС)
	fmt.Println(modifySlice(s)) // ВЫВЕДЕТ НОВЫЙ МАССИВ С НОВЫМ БАЗОВЫМ
}

func modifySlice(i []string) []string {
	i[0] = "3"         // 3 2 3
	i = append(i, "4") // 3 2 3 4 (НОВЫЙ БАЗОВЫЙ МАССИВ)
	i[1] = "5"         // 3 5 3 4
	i = append(i, "6") // 3 5 3 4 6
	return i
}
