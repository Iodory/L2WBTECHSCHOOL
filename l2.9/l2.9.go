package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func UnpackString(str string) (string, error) { // функция принимает строку на ввод и распаковывает ее
	if str == "" {                              // сразу проверка на пустую строку
		return "", nil // возвращение прустой строки
	}

	var result strings.Builder // тип strings.Builder полезен он создан для таких ситуаций и не пересоздает строку каждую иттерацию
	var prev rune              // тип рун последней буквы
	var hasPrev bool           // тип бул на проверку true если есть текущая буква для повторения
	escaped := false           // проверка на escape

	for i, r := range str { // цикл в котором будет все проверяться
		switch {
		case escaped: // вариант с escape
			if hasPrev { // проверка на есть ли предыдущая руна
				result.WriteRune(prev) // если есть будет добавлять ее к строке для вывода
			}
			prev = r
			hasPrev = true
			escaped = false
		case r == '\\': // экранируем следующий символ
			if hasPrev {
				result.WriteRune(prev) // если была буква перед то добавим ее
				hasPrev = false        // буква "израсходована"
			}
			escaped = true // ставим флаг, что следующий символ нужно воспринимать как обычный даже если цифра
		case unicode.IsLetter(r): // вариант если буква
			if hasPrev { // проверка на есть ли предыдущая руна
				result.WriteRune(prev)
			}
			prev = r
			hasPrev = true
		case unicode.IsDigit(r): // варинат если цифра
			if !hasPrev { // проверка на есть ли предыдущая руна что бы если НЕ было то была ошибка (мы не можем прибавить к пустоте)
				return "", fmt.Errorf("некорректная строка цифра в позиции %d", i)
			}

			count, _ := strconv.Atoi(string(r))                     // конвертируем цифру из рун в int
			result.WriteString(strings.Repeat(string(prev), count)) // через метод из пакета strings повторяем последнюю букву count раз
			hasPrev = false                                         // буква "израсходована"
		default: // случай если ни один не подходит
			return "", fmt.Errorf("некорректный символ %q", r)
		}
	}

	if escaped {
		return "", fmt.Errorf("некорректная строка окончание на escape")
	}
	if hasPrev {
		result.WriteRune(prev)
	}
	return result.String(), nil
}

func main() {
	tests := []string{
		"a4bc2d5e", // обычная распаковка
		"abcd",     // без цифр
		"45",       // начинается с цифры → ошибка
		"",         // пустая строка
		`qwe\4\5`,  // escape → цифры трактуются как буквы
		`qwe\45`,   // escape + цифра
	}

	for _, test := range tests {
		result, err := UnpackString(test)
		if err != nil {
			fmt.Printf("In: %-10q -> error: %v\n", test, err)
		} else {
			fmt.Printf("In: %-10q -> Out: %q\n", test, result)
		}
	}
}
