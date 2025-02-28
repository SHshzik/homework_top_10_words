package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

type Entry struct {
	text  string
	count int
}

func findTop(filename string) []Entry {
	var result = map[string]int{}
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err.Error())
		return []Entry{}
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		words := strings.Split(line, " ")
		for _, word := range words {
			if word == "\r" {
				continue
			}
			if _, ok := result[word]; !ok {
				result[word] = 0
			}
			result[word] += 1
		}
	}
	entries := make([]Entry, 0, len(result))
	for key, value := range result {
		entries = append(entries, Entry{key, value})
	}
	slices.SortFunc(entries, func(a, b Entry) int {
		return b.count - a.count
	})
	m := len(result)
	if m > 10 {
		m = 10
	}
	return entries[:m]
}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Usage: go run main.go <command>")
		return
	}
	files := arguments[1:]
	for _, file := range files {
		result := findTop(file)
		fmt.Println(result)
	}
}
