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

// process обрабатывает ввод по строкам и выводит результат
func process(scanner *bufio.Scanner, pattern string, ignoreReg, withoutPat, numLine, countLine, fixLines bool, prevLines, afterLines int) {
	prevList := make([]string, 0, prevLines)
	afterList := 0
	lineNum := 1
	count := 0

	for scanner.Scan() {
		line := scanner.Text()
		text := line
		pat := pattern

		if ignoreReg {
			text = strings.ToLower(line)
			pat = strings.ToLower(pattern)
		}

		prevList = append(prevList, line)
		if len(prevList) > prevLines {
			prevList = prevList[1:]
		}

		contains := false
		if fixLines {
			contains = text == pat
		} else {
			ok, err := regexp.MatchString(pat, text)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			contains = ok
		}

		if (withoutPat && !contains) || (!withoutPat && contains) {
			if len(prevList) > 1 {
				for _, l := range prevList[:len(prevList)-1] {
					fmt.Println(l)
				}
			}
			if numLine {
				fmt.Println(lineNum, line)
			} else {
				fmt.Println(line)
			}
			count++
			afterList = afterLines
		} else if afterList > 0 {
			fmt.Println(line)
			afterList--
		}
		lineNum++
	}

	if countLine {
		fmt.Println(count)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading error:", err)
	}
}

func main() {
	// Флаги
	ignoreReg := flag.Bool("i", false, "ignore case")
	withoutPat := flag.Bool("v", false, "select non-matching lines")
	numLine := flag.Bool("n", false, "print line number")
	countLine := flag.Bool("c", false, "print count of matching lines")
	prevLines := flag.Int("B", 0, "number of lines to print before match")
	afterLines := flag.Int("A", 0, "number of lines to print after match")
	cLines := flag.Int("C", 0, "number of lines to print before and after match")
	fixLines := flag.Bool("F", false, "match fixed strings (exact match)")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("usage: go run main.go [OPTIONS] PATTERN [FILE...]")
	}

	pattern := args[0]
	files := args[1:]

	if *cLines > 0 {
		if *prevLines < *cLines {
			*prevLines = *cLines
		}
		if *afterLines < *cLines {
			*afterLines = *cLines
		}
	}

	if len(files) > 0 {
		for _, fileName := range files {
			f, err := os.Open(fileName)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error opening file:", err)
				continue
			}
			defer f.Close()
			scanner := bufio.NewScanner(f)
			process(scanner, pattern, *ignoreReg, *withoutPat, *numLine, *countLine, *fixLines, *prevLines, *afterLines)
		}
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		process(scanner, pattern, *ignoreReg, *withoutPat, *numLine, *countLine, *fixLines, *prevLines, *afterLines)
	}
}
