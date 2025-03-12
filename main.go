package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"
	"sync"
)

type Entry struct {
	text  string
	count int
}

func findTop(filename string) ([]Entry, error) {
	var result = make(map[string]int)
	file, err := os.Open(filename)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			slog.Error(err.Error())
		}
	}(file)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		words := strings.Split(line, " ")
		for _, word := range words {
			if word == "\r" {
				continue
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
	return entries[:m], nil
}

func main() {
	var wg sync.WaitGroup
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Usage: go run main.go file1 file2 file3")
		return
	}
	files := arguments[1:]
	for _, file := range files {
		wg.Add(1)
		go func(file string) {
			result, err := findTop(file)
			if err != nil {
				slog.Error(err.Error())
			}
			fmt.Println(file, result)
			wg.Done()
		}(file)
	}
	wg.Wait()
}
