package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func DivideFlag() {
	// функция делит флаг типа -nr на 2 отдельных флага для распознавания
	var newFlags []string
	newFlags = append(newFlags, os.Args[0])

	for _, flag := range os.Args[1:] {
		if len(flag) > 1 && flag[0] == '-' && flag[1] != '-' {
			// делим флаги
			for _, ch := range flag[1:] {
				newFlags = append(newFlags, "-"+string(ch))
			}
		} else {
			newFlags = append(newFlags, flag)
		}
	}
	os.Args = newFlags
}

func main() {
	DivideFlag() // вызываем деление флага до
	reverse := flag.Bool("r", false, "reverse order")
	unique := flag.Bool("u", false, "unique order")
	numbers := flag.Bool("n", false, "numbers order")
	column := flag.Int("k", 0, "sort by column (1-based, tab-separated)")
	month := flag.Bool("m", false, "sort by month")
	flag.Parse()

	// Определяем источник ввода: файл или stdin
	var scanner *bufio.Scanner
	if flag.NArg() > 0 {
		fileName := flag.Arg(0)
		f, err := os.Open(fileName)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error opening file:", err)
			os.Exit(1)
		}
		defer f.Close()
		scanner = bufio.NewScanner(f)
	} else {
		scanner = bufio.NewScanner(os.Stdin)
	}

	//чтение строк
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		os.Exit(1)
	}
	// сортировка
	sort.Slice(lines, func(i, j int) bool {
		// (-k) выбираем колонку
		var keyI, keyJ string
		if *column > 0 {
			fieldsI := strings.Fields(lines[i])
			fieldsJ := strings.Fields(lines[j])
			if len(fieldsI) >= *column {
				keyI = fieldsI[*column-1]
			}
			if len(fieldsJ) >= *column {
				keyJ = fieldsJ[*column-1]
			}
		} else {
			keyI = lines[i]
			keyJ = lines[j]
		}

		var less bool

		// (-M) сортировка по месяцу
		if *month {
			var monthOrder = map[string]int{
				"Jan": 1, "Feb": 2, "Mar": 3,
				"Apr": 4, "May": 5, "Jun": 6,
				"Jul": 7, "Aug": 8, "Sep": 9,
				"Oct": 10, "Nov": 11, "Dec": 12,
			}
			numI, okI := monthOrder[keyI]
			numJ, okJ := monthOrder[keyJ]
			if !okI {
				numI = 0
			}
			if !okJ {
				numJ = 0
			}
			less = numI < numJ
		} else if *numbers { // (-n) сортировка по числу
			ni, err1 := strconv.ParseFloat(keyI, 64)
			nj, err2 := strconv.ParseFloat(keyJ, 64)
			if err1 == nil && err2 == nil {
				less = ni < nj
			} else {
				less = keyI < keyJ
			}
		} else {
			// обычная сортировка
			less = keyI < keyJ
		}

		// (-r) обратный порядок
		if *reverse {
			return !less
		}
		return less
	})

	// (-u) убираем дубликаты
	if *unique {
		lines = Unique(lines)
	}

	for _, line := range lines {
		fmt.Println(line)
	}
}

func Unique(str []string) []string { // функция для вывода только уникальных символов
	checker := make(map[string]bool)
	var result []string
	for _, ch := range str {
		// простая проверка на новизну символа
		if !checker[ch] {
			// если новый то ставим статус "уже есть" и добавляем в результат
			checker[ch] = true
			result = append(result, ch)
		}
	}
	return result
}
