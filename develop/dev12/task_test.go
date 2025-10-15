package main

import (
	"bytes"
	"os"
	"testing"
)

func TestReadLinesFromFile(t *testing.T) {
	tempFile, err := os.CreateTemp("", "test*.txt")
	if err != nil {
		t.Fatal("Failed to create temp file")
	}
	defer os.Remove(tempFile.Name())

	testData := "first line\nsecond line\nthird line"
	_, err = tempFile.WriteString(testData)
	if err != nil {
		t.Fatal("Failed to write to file")
	}
	tempFile.Close()

	lines, err := readLines(tempFile.Name())
	if err != nil {
		t.Errorf("Error reading file: %v", err)
	}

	expected := []string{"first line", "second line", "third line"}
	if len(lines) != len(expected) {
		t.Errorf("Expected %d lines, got %d", len(expected), len(lines))
	}

	for i, line := range lines {
		if line != expected[i] {
			t.Errorf("Line %d: expected '%s', got '%s'", i, expected[i], line)
		}
	}
}

func TestReadLinesFromStdin(t *testing.T) {
	input := "stdin line\nanother line"

	r, w, _ := os.Pipe()
	oldStdin := os.Stdin
	os.Stdin = r

	defer func() {
		os.Stdin = oldStdin
		r.Close()
	}()

	go func() {
		w.WriteString(input)
		w.Close()
	}()

	lines, err := readLines("")
	if err != nil {
		t.Errorf("Error reading from stdin: %v", err)
	}

	expected := []string{"stdin line", "another line"}

	if len(lines) != len(expected) {
		t.Errorf("Expected %d lines, got %d", len(expected), len(lines))
	}

	for i := range expected {
		if lines[i] != expected[i] {
			t.Errorf("Line %d: Expected '%s', got '%s'", i, expected[i], lines[i])
		}
	}
}

func TestFixedStringSearch(t *testing.T) {
	lines := []string{
		"Hello World",
		"test line",
		"hello there",
		"HELLO everyone",
		"goodbye",
	}

	cfg := config{
		pattern: "hello",
		fixed:   true,
	}

	matches := findMatches(lines, cfg)

	// Должен найти только точное совпадение с "hello"
	expected := []int{2} // "hello there"
	if len(matches) != len(expected) {
		t.Errorf("Found %d matches, expected %d", len(matches), len(expected))
	}
}

func TestCaseInsensitiveSearch(t *testing.T) {
	lines := []string{
		"Hello World",
		"test line",
		"hello there",
		"HELLO everyone",
		"goodbye",
	}

	cfg := config{
		pattern:    "hello",
		fixed:      true,
		ignoreCase: true,
	}

	matches := findMatches(lines, cfg)

	// Должен найти все hello варианты
	expected := []int{0, 2, 3}
	if len(matches) != len(expected) {
		t.Errorf("Found %d matches, expected %d", len(matches), len(expected))
	}
}

func TestInvertedSearch(t *testing.T) {
	lines := []string{
		"apple",
		"banana",
		"apple pie",
		"cherry",
	}

	cfg := config{
		pattern: "apple",
		fixed:   true,
		invert:  true,
	}

	matches := findMatches(lines, cfg)

	// Должен найти строки без "apple"
	expected := []int{1, 3} // "banana", "cherry"
	if len(matches) != len(expected) {
		t.Errorf("Found %d matches, expected %d", len(matches), len(expected))
	}
}

func TestCountOnly(t *testing.T) {
	lines := []string{
		"test one",
		"test two",
		"no match",
		"test three",
	}

	cfg := config{
		pattern: "test",
		fixed:   true,
		count:   true,
	}

	matches := findMatches(lines, cfg)

	// Должен найти 3 строки "test"
	if len(matches) != 3 {
		t.Errorf("Found %d matches, expected 3", len(matches))
	}
}

func TestLineNumberOutput(t *testing.T) {
	lines := []string{"first", "second", "third"}
	matches := []int{1}

	// Capture output
	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cfg := config{lineNumber: true}
	printMatches(lines, matches, cfg)

	w.Close()
	os.Stdout = oldStdout
	buf.ReadFrom(r)
	output := buf.String()

	// Проверяет что строка выводится с номером строки
	expected := "2:second\n"
	if output != expected {
		t.Errorf("Output: '%s', expected: '%s'", output, expected)
	}
}

func TestContextOutput(t *testing.T) {
	lines := []string{
		"line 1",
		"line 2",
		"MATCH",
		"line 4",
		"line 5",
	}

	matches := []int{2}

	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cfg := config{
		before: 1,
		after:  1,
	}
	printMatches(lines, matches, cfg)

	w.Close()
	os.Stdout = oldStdout
	buf.ReadFrom(r)
	output := buf.String()

	// Должен вывести 1 строку перед, саму строку, 1 строку после
	expected := "line 2\nMATCH\nline 4\n"
	if output != expected {
		t.Errorf("Output: '%s', expected: '%s'", output, expected)
	}
}

func TestBeginningOfFile(t *testing.T) {
	lines := []string{
		"MATCH",
		"line 2",
		"line 3",
	}

	matches := []int{0}

	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cfg := config{before: 2} // Должен выводить 2 строки до, но их нет
	printMatches(lines, matches, cfg)

	w.Close()
	os.Stdout = oldStdout
	buf.ReadFrom(r)
	output := buf.String()

	// Должен вывести только первую строку
	expected := "MATCH\n"
	if output != expected {
		t.Errorf("Output: '%s', expected: '%s'", output, expected)
	}
}

func TestRegexSearch(t *testing.T) {
	lines := []string{
		"test123",
		"hello",
		"456test",
		"no match",
	}

	cfg := config{
		pattern: `\d+`,
		fixed:   false,
	}

	matches := findMatches(lines, cfg)

	// Должен найти строки с числами
	expected := []int{0, 2}
	if len(matches) != len(expected) {
		t.Errorf("Found %d matches, expected %d", len(matches), len(expected))
	}
}

func TestMaxMinFunctions(t *testing.T) {
	if max(5, 3) != 5 {
		t.Error("max(5, 3) should return 5")
	}
	if max(2, 8) != 8 {
		t.Error("max(2, 8) should return 8")
	}

	if min(5, 3) != 3 {
		t.Error("min(5, 3) should return 3")
	}
	if min(2, 8) != 2 {
		t.Error("min(2, 8) should return 2")
	}
}

func TestMultipleMatches(t *testing.T) {
	lines := []string{
		"apple",
		"banana",
		"apple",
		"cherry",
		"apple",
	}

	cfg := config{
		pattern: "apple",
		fixed:   true,
	}

	matches := findMatches(lines, cfg)

	// Должен найти все три "apple"
	if len(matches) != 3 {
		t.Errorf("Found %d matches, expected 3", len(matches))
	}

	expected := []int{0, 2, 4}
	for i, idx := range matches {
		if idx != expected[i] {
			t.Errorf("Match %d: index %d, expected %d", i, idx, expected[i])
		}
	}
}

func TestOverlappingContext(t *testing.T) {
	lines := []string{
		"line 1",
		"MATCH 1",
		"line 3",
		"MATCH 2",
		"line 5",
	}

	matches := []int{1, 3}

	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cfg := config{before: 1, after: 1}
	printMatches(lines, matches, cfg)

	w.Close()
	os.Stdout = oldStdout
	buf.ReadFrom(r)
	output := buf.String()

	// Должен вывести все строки без дубликатов
	expected := "line 1\nMATCH 1\nline 3\nMATCH 2\nline 5\n"
	if output != expected {
		t.Errorf("Output: '%s', expected: '%s'", output, expected)
	}
}

func TestEndOfFileContext(t *testing.T) {
	lines := []string{
		"line 1",
		"line 2",
		"MATCH",
	}

	matches := []int{2}

	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cfg := config{after: 2} // Требуется 2 строки после последней, но их нет
	printMatches(lines, matches, cfg)

	w.Close()
	os.Stdout = oldStdout
	buf.ReadFrom(r)
	output := buf.String()

	expected := "MATCH\n"
	if output != expected {
		t.Errorf("Output: '%s', expected: '%s'", output, expected)
	}
}
