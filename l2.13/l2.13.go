package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	fields := flag.String("f", "", "display fields")
	delimiter := flag.String("d", "\t", "display fields delimiter")
	separated := flag.Bool("s", false, "display all fields")
	flag.Parse()

	// -f
	var nums []int
	if *fields != "" {
		parts := strings.Split(*fields, ",")
		for _, part := range parts {
			if strings.Contains(part, "-") {
				lPart := strings.SplitN(part, "-", 2)
				lN, err := strconv.Atoi(lPart[0])
				if err != nil {
					log.Fatal(err)
				}
				rN, err := strconv.Atoi(lPart[1])
				if err != nil {
					log.Fatal(err)
				}
				for i := lN; i <= rN; i++ {
					nums = append(nums, i)
				}
			} else {
				num, err := strconv.Atoi(part)
				if err != nil {
					log.Fatal(err)
				}
				nums = append(nums, num)
			}
		}
	}
	// читаем ввод
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		// -s
		if !strings.Contains(line, *delimiter) {
			if *separated {
				continue
			}
			if *fields == "" {
				fmt.Println(line)
			}
			continue
		}

		parts := strings.Split(line, *delimiter)

		var out []string
		for _, num := range nums {
			if num-1 >= 0 && num-1 < len(parts) {
				out = append(out, parts[num-1])
			}
		}
		if len(out) > 0 {
			fmt.Println(strings.Join(out, *delimiter))
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
