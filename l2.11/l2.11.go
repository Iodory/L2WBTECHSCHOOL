package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	words := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}
	result := Anagramm(words)

	for k, v := range result {
		fmt.Println(k, ":", v)
	}
}

func Anagramm(words []string) map[string][]string {
	// временная мапа: ключ = отсортированные буквы, значение = список слов
	temp := make(map[string][]string)
	// чтобы помнить порядок появления слов
	order := make(map[string]string)

	for _, w := range words {
		lower := strings.ToLower(w) // пятак
		runes := []rune(lower)      // п я т а к

		// сортировка букв
		sort.Slice(runes, func(i, j int) bool { // а к п т я
			return runes[i] < runes[j]
		})
		key := string(runes)

		// сохраняем слова в группу
		temp[key] = append(temp[key], lower)

		// если это первое слово с таким ключом запомним его
		if _, ok := order[key]; !ok {
			order[key] = lower
		}
	}

	// финальная обработка
	result := make(map[string][]string)
	for key, group := range temp {
		if len(group) > 1 {
			sort.Strings(group)        // сортируем слова в группе
			result[order[key]] = group // ключ = первое встретившееся слово
		}
	}

	return result
}
