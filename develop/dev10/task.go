package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// Config хранит конфигурацию флагов сортировки
type Config struct {
	column       int
	numeric      bool
	reverse      bool
	unique       bool
	month        bool
	ignoreBlanks bool
	checkSorted  bool
	humanNumeric bool
}

// MonthToNumber преобразует название месяца в число
func MonthToNumber(month string) int {
	months := map[string]int{
		"jan": 1, "feb": 2, "mar": 3, "apr": 4, "may": 5, "jun": 6,
		"jul": 7, "aug": 8, "sep": 9, "oct": 10, "nov": 11, "dec": 12,
	}

	monthLower := strings.ToLower(month)
	if num, exists := months[monthLower]; exists {
		return num
	}
	return 0
}

// ParseHumanNumber парсит числа с суффиксами (K, M, G, T, P)
func ParseHumanNumber(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}

	multipliers := map[byte]float64{
		'K': 1 << 10,
		'M': 1 << 20,
		'G': 1 << 30,
		'T': 1 << 40,
		'P': 1 << 50,
	}

	lastChar := s[len(s)-1]
	if multiplier, exists := multipliers[lastChar]; exists {
		numStr := strings.TrimSpace(s[:len(s)-1])
		num, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			return 0
		}
		return num * multiplier
	}

	num, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return num
}

// GetColumn извлекает колонку из строки
func GetColumn(line string, column int, ignoreBlanks bool) string {
	if column <= 0 {
		if ignoreBlanks {
			return strings.TrimRightFunc(line, unicode.IsSpace)
		}
		return line
	}

	fields := strings.Split(line, "\t")
	if column > len(fields) {
		return ""
	}

	field := fields[column-1]
	if ignoreBlanks {
		return strings.TrimRightFunc(field, unicode.IsSpace)
	}
	return field
}

// CompareValues сравнивает значения в зависимости от типа сортировки
func CompareValues(a, b string, cfg Config) bool {
	switch {
	case cfg.numeric:
		return compareNumeric(a, b)
	case cfg.month:
		return compareMonth(a, b)
	case cfg.humanNumeric:
		return compareHumanNumeric(a, b)
	default:
		return a < b
	}
}

// compareNumeric сравнивает строки как числа
func compareNumeric(a, b string) bool {
	// Сначала проверяем специальные случаи с пустой строкой
	if a == "" && b == "" {
		return false // две пустые строки равны
	}
	if a == "" {
		return true // пустая строка всегда меньше любой непустой
	}
	if b == "" {
		return false // любая непустая строка всегда больше пустой
	}

	numA, errA := strconv.ParseFloat(a, 64)
	numB, errB := strconv.ParseFloat(b, 64)

	switch {
	case errA == nil && errB == nil:
		return numA < numB
	case errA != nil && errB != nil:
		return a < b
	case errA == nil:
		return true // a - число, b - не число: число меньше
	default:
		return false // a - не число, b - число: не число больше
	}
}

// compareMonth сравнивает строки как названия месяцев
func compareMonth(a, b string) bool {
	monthA := MonthToNumber(a)
	monthB := MonthToNumber(b)

	switch {
	case monthA > 0 && monthB > 0:
		return monthA < monthB
	case monthA > 0:
		return true
	case monthB > 0:
		return false
	default:
		return a < b
	}
}

// compareHumanNumeric сравнивает строки как человеко-читаемые числа
func compareHumanNumeric(a, b string) bool {
	numA := ParseHumanNumber(a)
	numB := ParseHumanNumber(b)
	return numA < numB
}

// ReadLines читает строки из файла или stdin
func ReadLines(filename string) ([]string, error) {
	var input io.Reader = os.Stdin

	if filename != "" {
		file, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		input = file
	}

	var lines []string
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// RemoveDuplicates удаляет дубликаты из отсортированного списка
func RemoveDuplicates(lines []string) []string {
	if len(lines) == 0 {
		return lines
	}

	result := []string{lines[0]}
	for i := 1; i < len(lines); i++ {
		if lines[i] != lines[i-1] {
			result = append(result, lines[i])
		}
	}
	return result
}

// IsSorted проверяет, отсортированы ли данные
func IsSorted(lines []string, cfg Config) bool {
	if len(lines) <= 1 {
		return true
	}

	for i := 1; i < len(lines); i++ {
		colA := GetColumn(lines[i-1], cfg.column, cfg.ignoreBlanks)
		colB := GetColumn(lines[i], cfg.column, cfg.ignoreBlanks)

		less := CompareValues(colA, colB, cfg)
		if cfg.reverse {
			less = !less
		}

		if !less && colA != colB {
			return false
		}
	}
	return true
}

// processFlags обрабатывает флаги и возвращает конфигурацию
func processFlags() (Config, string) {
	var cfg Config

	flag.IntVar(&cfg.column, "k", 0, "sort by column number")
	flag.BoolVar(&cfg.numeric, "n", false, "sort by numeric value")
	flag.BoolVar(&cfg.reverse, "r", false, "reverse sort order")
	flag.BoolVar(&cfg.unique, "u", false, "output only unique lines")
	flag.BoolVar(&cfg.month, "M", false, "sort by month names")
	flag.BoolVar(&cfg.ignoreBlanks, "b", false, "ignore trailing blanks")
	flag.BoolVar(&cfg.checkSorted, "c", false, "check if data is sorted")
	flag.BoolVar(&cfg.humanNumeric, "h", false, "sort by human-readable numbers")

	flag.Parse()

	filename := ""
	if flag.NArg() > 0 {
		filename = flag.Arg(0)
	}

	return cfg, filename
}

func main() {
	cfg, filename := processFlags()

	lines, err := ReadLines(filename)
	if err != nil {
		log.Fatalf("Error reading input: %v", err)
	}

	if cfg.checkSorted {
		if IsSorted(lines, cfg) {
			os.Exit(0)
		}
		fmt.Fprintln(os.Stderr, "Data is not sorted")
		os.Exit(1)
	}

	sort.SliceStable(lines, func(i, j int) bool {
		a := GetColumn(lines[i], cfg.column, cfg.ignoreBlanks)
		b := GetColumn(lines[j], cfg.column, cfg.ignoreBlanks)

		result := CompareValues(a, b, cfg)

		if cfg.reverse {
			return !result
		}
		return result
	})

	if cfg.unique {
		lines = RemoveDuplicates(lines)
	}

	for _, line := range lines {
		fmt.Println(line)
	}
}
