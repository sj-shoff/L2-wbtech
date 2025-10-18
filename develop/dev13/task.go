package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type fieldRange struct {
	start int
	end   int
}

type config struct {
	fields    string
	delimiter string
	separated bool
}

func main() {
	cfg := parseFlags()

	if err := runCut(os.Stdin, os.Stdout, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// Парсинг флагов
func parseFlags() config {
	var cfg config

	flag.StringVar(&cfg.fields, "f", "", "list of fields to extract (e.g., 1,3-5)")
	flag.StringVar(&cfg.delimiter, "d", "\t", "delimiter character")
	flag.BoolVar(&cfg.separated, "s", false, "only output lines containing delimiter")

	flag.Parse()

	if cfg.fields == "" {
		fmt.Fprintln(os.Stderr, "Error: fields list must be specified with -f")
		os.Exit(1)
	}

	return cfg
}

// Парсинг полей
func parseFields(fieldsStr string) ([]fieldRange, error) {
	var ranges []fieldRange

	parts := strings.Split(fieldsStr, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if strings.Contains(part, "-") {
			rangeParts := strings.SplitN(part, "-", 2)
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("invalid range format: %s", part)
			}

			start, err := strconv.Atoi(strings.TrimSpace(rangeParts[0]))

			if err != nil || start < 1 {
				return nil, fmt.Errorf("invalid start value in range: %s", rangeParts[0])
			}

			end, err := strconv.Atoi(strings.TrimSpace(rangeParts[1]))

			if err != nil || end < 1 {
				if err != nil {
					return nil, fmt.Errorf("invalid end value in range: %s", rangeParts[1])
				}
			}

			if start > end {
				return nil, fmt.Errorf("invalid range: start (%d) is greater than end (%d)", start, end)
			}

			ranges = append(ranges, fieldRange{start, end})

		} else {
			field, err := strconv.Atoi(part)
			if err != nil || field < 1 {
				return nil, fmt.Errorf("invalid field value: %s", part)
			}

			ranges = append(ranges, fieldRange{field, field})
		}
	}

	return ranges, nil
}

// Реализация утилиты cut
func runCut(input io.Reader, output io.Writer, cfg config) error {
	fieldRanges, err := parseFields(cfg.fields)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(input)
	delimiter := cfg.delimiter

	for scanner.Scan() {
		line := scanner.Text()

		if !strings.Contains(line, delimiter) {
			if !cfg.separated {
				fmt.Fprintln(output, line)
			}
			continue
		}

		fields := strings.Split(line, delimiter)
		var outputFields []string

		for _, r := range fieldRanges {
			for i := r.start; i <= r.end; i++ {
				if i <= len(fields) {
					outputFields = append(outputFields, fields[i-1])
				}
			}
		}

		if len(outputFields) > 0 {
			fmt.Fprintln(output, strings.Join(outputFields, delimiter))
		}
	}

	return scanner.Err()
}
