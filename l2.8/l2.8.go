package main

import (
	"fmt"
	"log"

	"github.com/beevik/ntp"
)

func main() { // вход в функцию и начало работы
	curTime, err := ntp.Time("pool.ntp.org") // вызов функциюю из пакета ntp для вывода текущего времени по ntp
	if err != nil {                          // проверка на ошибку в случае != nil
		log.Fatal(err) // логирование ошибки
	}
	fmt.Println("Текущее время: ", curTime) // вывод текущего времени
}
