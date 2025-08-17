package main

import (
	"fmt"
	"os"
)

func Foo() error {
	var err *os.PathError = nil // значение переменной == nil
	return err                  // но возвращаем переменную как интерфейс (происходит "упаковка" в error)
	// это работает потому что у PathError есть error метод
}

func main() {
	err := Foo()
	fmt.Println(err)        // выводит nil потому что мы присвоили VALUE nil
	fmt.Println(err == nil) // выводит false ведь v == nil НО type == os.PathError
}
