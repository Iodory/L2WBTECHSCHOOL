package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	// все в нижний регистр
	ignoreReg := flag.Bool("i", false, "ignore register")
	// только что не подходит
	withoutPat := flag.Bool("v", false, "without pattern")
	// номер строки
	numLine := flag.Bool("n", false, "number of lines")
	//счетчик
	countLine := flag.Bool("c", false, "count lines")
	//число строк ДО
	prevLines := flag.Int("B", 0, "previous lines")
	//число строк после
	afterLines := flag.Int("A", 0, "after lines")
	//число и до и после
	cLines := flag.Int("C", 0, "before+after")
	//фиксированная строка
	fixLines := flag.Bool("F", false, "fix lines")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("usage: go run main.go PATTERN")
	}

	pattern := args[0]
	files := args[1:]

	if len(files) > 0 {
		for _, fileName := range files {
			f, err := os.Open(fileName)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error opening file: ", err)
				continue
			}

			if *cLines > 0 {
				if *prevLines < *cLines {
					*prevLines = *cLines
				}
				if *afterLines < *cLines {
					*afterLines = *cLines
				}
			}

			prevList := make([]string, 0, *prevLines)
			afterList := 0

			lineNum := 1
			count := 0

			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := scanner.Text()
				text := line
				pat := pattern

				if *ignoreReg {
					text = strings.ToLower(line)
					pat = strings.ToLower(pattern)
				}

				prevList = append(prevList, line)
				if len(prevList) > *prevLines {
					prevList = prevList[1:] // удаляем старую строку если превышаем
				}

				contains := false
				if *fixLines {
					//поиск подстроки
					contains = text == pat
				} else {
					//поиск регулярки
					ok, err := regexp.MatchString(pat, text)
					if err != nil {
						fmt.Fprintln(os.Stderr, err)
						continue
					}
					contains = ok
				}
				if (*withoutPat && !contains) || (!*withoutPat && contains) {
					//вывод строк из prevList
					if len(prevList) > 1 {
						for _, l := range prevList[:len(prevList)-1] {
							fmt.Println(l)
						}
					}
					//вывод текущей строки
					if *numLine {
						fmt.Println(lineNum, line)
					} else {
						fmt.Println(line)
					}
					count++
					afterList = *afterLines
				} else if afterList > 0 {
					fmt.Println(line)
					afterList--
				}
				lineNum++
			}
			if *countLine {
				fmt.Println(count)
			}
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "reading file: ", err)
			}
			f.Close()
		}
	} else {

		prevList := make([]string, 0, *prevLines)
		afterList := 0

		lineNum := 1
		count := 0

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			text := line
			pat := pattern

			if *ignoreReg {
				text = strings.ToLower(line)
				pat = strings.ToLower(pattern)
			}

			prevList = append(prevList, line)
			if len(prevList) > *prevLines {
				prevList = prevList[1:] // удаляем старую строку если превышаем
			}

			contains := false
			if *fixLines {
				//поиск подстроки
				contains = text == pat
			} else {
				//поиск регулярки
				ok, err := regexp.MatchString(pat, text)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					continue
				}
				contains = ok
			}

			if (*withoutPat && !contains) || (!*withoutPat && contains) {
				//вывод строк из prevList
				if len(prevList) > 1 {
					for _, l := range prevList[:len(prevList)-1] {
						fmt.Println(l)
					}
				}
				//вывод текущей строки
				if *numLine {
					fmt.Println(lineNum, line)
				} else {
					fmt.Println(line)
				}
				count++
				afterList = *afterLines
			} else if afterList > 0 {
				fmt.Println(line)
				afterList--
			}
			lineNum++
		}
		if *countLine {
			fmt.Println(count)
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading file: ", err)
		}
	}
}
