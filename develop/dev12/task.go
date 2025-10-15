package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Конфиг для парсинга флагов
type config struct {
	after      int
	before     int
	context    int
	count      bool
	ignoreCase bool
	invert     bool
	fixed      bool
	lineNumber bool
	pattern    string
	filename   string
}

func main() {
	cfg := parseFlags()

	if cfg.pattern == "" {
		fmt.Fprintln(os.Stderr, "Pattern is required")
		os.Exit(1)
	}

	lines, err := readLines(cfg.filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	matches := findMatches(lines, cfg)

	if cfg.count {
		fmt.Println(len(matches))
		return
	}

	printMatches(lines, matches, cfg)
}

// Парсинг флагов
func parseFlags() config {
	var cfg config

	flag.IntVar(&cfg.after, "A", 0, "Print N lines after match")
	flag.IntVar(&cfg.before, "B", 0, "Print N lines before match")
	flag.IntVar(&cfg.context, "C", 0, "Print N lines of context")
	flag.BoolVar(&cfg.count, "c", false, "Count matching lines only")
	flag.BoolVar(&cfg.ignoreCase, "i", false, "Ignore case")
	flag.BoolVar(&cfg.invert, "v", false, "Invert match")
	flag.BoolVar(&cfg.fixed, "F", false, "Fixed string match")
	flag.BoolVar(&cfg.lineNumber, "n", false, "Print line numbers")

	flag.Parse()

	if cfg.context > 0 {
		cfg.before = cfg.context
		cfg.after = cfg.context
	}

	args := flag.Args()
	if len(args) > 0 {
		cfg.pattern = args[0]
	}
	if len(args) > 1 {
		cfg.filename = args[1]
	}

	return cfg
}

// Чтение строк из файла или стандартного ввода
func readLines(filename string) ([]string, error) {
	var scanner *bufio.Scanner

	if filename == "" {
		scanner = bufio.NewScanner(os.Stdin)
	} else {
		file, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		scanner = bufio.NewScanner(file)
	}

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

// Поиск совпадений
func findMatches(lines []string, cfg config) []int {
	var matches []int
	matcher := createMatcher(cfg)

	for i, line := range lines {
		if matcher(line) != cfg.invert {
			matches = append(matches, i)
		}
	}

	return matches
}

// Создание функции для поиска
func createMatcher(cfg config) func(string) bool {
	if cfg.fixed {
		if cfg.ignoreCase {
			pattern := strings.ToLower(cfg.pattern)
			return func(line string) bool {
				return strings.Contains(strings.ToLower(line), pattern)
			}
		}
		return func(line string) bool {
			return strings.Contains(line, cfg.pattern)
		}
	}

	// Для регулярных выражений
	pattern := cfg.pattern
	if cfg.ignoreCase {
		pattern = "(?i)" + pattern
	}
	re := regexp.MustCompile(pattern)
	return func(line string) bool {
		return re.MatchString(line)
	}
}

// Вывод результатов
func printMatches(lines []string, matches []int, cfg config) {
	printed := make(map[int]bool)

	for _, matchIdx := range matches {

		start := max(0, matchIdx-cfg.before)
		for i := start; i < matchIdx; i++ {
			if !printed[i] {
				printLine(lines, i, cfg)
				printed[i] = true
			}
		}

		if !printed[matchIdx] {
			printLine(lines, matchIdx, cfg)
			printed[matchIdx] = true
		}

		end := min(len(lines)-1, matchIdx+cfg.after)
		for i := matchIdx + 1; i <= end; i++ {
			if !printed[i] {
				printLine(lines, i, cfg)
				printed[i] = true
			}
		}
	}
}

// Вывод строки
func printLine(lines []string, idx int, cfg config) {
	if cfg.lineNumber {
		fmt.Printf("%d:", idx+1)
	}
	fmt.Println(lines[idx])
}

// Доп. функция для макс числа
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Доп. функция для мин числа
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
