package main

import (
	"flag"
	"os"
	"reflect"
	"testing"
)

func TestMonthToNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"January", "jan", 1},
		{"February", "feb", 2},
		{"March", "mar", 3},
		{"April", "apr", 4},
		{"May", "may", 5},
		{"June", "jun", 6},
		{"July", "jul", 7},
		{"August", "aug", 8},
		{"September", "sep", 9},
		{"October", "oct", 10},
		{"November", "nov", 11},
		{"December", "dec", 12},
		{"Uppercase", "JAN", 1},
		{"Mixed case", "Feb", 2},
		{"Invalid month", "invalid", 0},
		{"Empty string", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MonthToNumber(tt.input)
			if result != tt.expected {
				t.Errorf("MonthToNumber(%q) = %d, expected %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseHumanNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{"Kilobytes", "1K", 1024},
		{"Megabytes", "2M", 2 * 1024 * 1024},
		{"Gigabytes", "1.5G", 1.5 * 1024 * 1024 * 1024},
		{"Terabytes", "1T", 1 * 1024 * 1024 * 1024 * 1024},
		{"Petabytes", "1P", 1 * 1024 * 1024 * 1024 * 1024 * 1024},
		{"Regular number", "100", 100},
		{"Decimal number", "3.14", 3.14},
		{"With spaces", " 2 K ", 2 * 1024},
		{"Invalid suffix", "1X", 0},
		{"Empty string", "", 0},
		{"Only letters", "abc", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseHumanNumber(tt.input)
			if result != tt.expected {
				t.Errorf("ParseHumanNumber(%q) = %f, expected %f", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetColumn(t *testing.T) {
	tests := []struct {
		name         string
		line         string
		column       int
		ignoreBlanks bool
		expected     string
	}{
		{"First column", "apple\tbanana\tcherry", 1, false, "apple"},
		{"Second column", "apple\tbanana\tcherry", 2, false, "banana"},
		{"Third column", "apple\tbanana\tcherry", 3, false, "cherry"},
		{"Column out of range", "apple\tbanana", 5, false, ""},
		{"Zero column returns whole line", "apple\tbanana", 0, false, "apple\tbanana"},
		{"Ignore blanks with spaces", "apple \tbanana\tcherry  ", 1, true, "apple"},
		{"Ignore blanks without spaces", "apple\tbanana\tcherry", 1, true, "apple"},
		{"Empty line", "", 1, false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetColumn(tt.line, tt.column, tt.ignoreBlanks)
			if result != tt.expected {
				t.Errorf("GetColumn(%q, %d, %t) = %q, expected %q",
					tt.line, tt.column, tt.ignoreBlanks, result, tt.expected)
			}
		})
	}
}

func TestCompareNumeric(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		expected bool
	}{
		{"Smaller number", "5", "10", true},
		{"Larger number", "10", "5", false},
		{"Equal numbers", "5", "5", false},
		{"Negative numbers", "-5", "-3", true},
		{"Decimal numbers", "3.14", "2.71", false},
		{"Number vs non-number", "5", "apple", true},  // Число должно быть меньше не-числа
		{"Non-number vs number", "apple", "5", false}, // Не-число должно быть больше числа
		{"Both non-numbers", "apple", "banana", true}, // Обычное строковое сравнение
		{"Empty string as number", "", "5", true},     // Пустая строка должна быть меньше числа
		{"Number vs empty string", "5", "", false},    // Число должно быть больше пустой строки
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compareNumeric(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("compareNumeric(%q, %q) = %t, expected %t",
					tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestCompareMonth(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		expected bool
	}{
		{"January before March", "jan", "mar", true},
		{"December after January", "dec", "jan", false},
		{"Same month", "feb", "feb", false},
		{"Month vs non-month", "jan", "apple", true},
		{"Non-month vs month", "apple", "jan", false},
		{"Both non-months", "apple", "banana", true},
		{"Uppercase months", "JAN", "FEB", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compareMonth(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("compareMonth(%q, %q) = %t, expected %t",
					tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestCompareHumanNumeric(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		expected bool
	}{
		{"KB smaller than MB", "1K", "1M", true},
		{"MB larger than KB", "1M", "1K", false},
		{"Same values", "2K", "2K", false},
		{"Different sizes same value", "1024K", "1M", false},
		{"Regular numbers", "100", "200", true},
		{"Invalid numbers", "abc", "def", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compareHumanNumeric(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("compareHumanNumeric(%q, %q) = %t, expected %t",
					tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestCompareValues(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		cfg      Config
		expected bool
	}{
		{
			name:     "Default string comparison",
			a:        "apple",
			b:        "banana",
			cfg:      Config{},
			expected: true,
		},
		{
			name:     "Numeric comparison",
			a:        "5",
			b:        "10",
			cfg:      Config{numeric: true},
			expected: true,
		},
		{
			name:     "Month comparison - January before February",
			a:        "jan",
			b:        "feb",
			cfg:      Config{month: true},
			expected: true, // jan(1) < feb(2) = true
		},
		{
			name:     "Month comparison - February after January",
			a:        "feb",
			b:        "jan",
			cfg:      Config{month: true},
			expected: false, // feb(2) < jan(1) = false
		},
		{
			name:     "Human numeric comparison",
			a:        "1K",
			b:        "2K",
			cfg:      Config{humanNumeric: true},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CompareValues(tt.a, tt.b, tt.cfg)
			if result != tt.expected {
				t.Errorf("CompareValues(%q, %q, %+v) = %t, expected %t",
					tt.a, tt.b, tt.cfg, result, tt.expected)
			}
		})
	}
}

func TestRemoveDuplicates(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "No duplicates",
			input:    []string{"apple", "banana", "cherry"},
			expected: []string{"apple", "banana", "cherry"},
		},
		{
			name:     "With duplicates",
			input:    []string{"apple", "apple", "banana", "banana", "cherry"},
			expected: []string{"apple", "banana", "cherry"},
		},
		{
			name:     "All duplicates",
			input:    []string{"apple", "apple", "apple"},
			expected: []string{"apple"},
		},
		{
			name:     "Empty slice",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "Single element",
			input:    []string{"apple"},
			expected: []string{"apple"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RemoveDuplicates(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("RemoveDuplicates(%v) = %v, expected %v",
					tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsSorted(t *testing.T) {
	tests := []struct {
		name     string
		lines    []string
		cfg      Config
		expected bool
	}{
		{
			name:     "Sorted strings",
			lines:    []string{"apple", "banana", "cherry"},
			cfg:      Config{column: 0},
			expected: true,
		},
		{
			name:     "Unsorted strings",
			lines:    []string{"cherry", "apple", "banana"},
			cfg:      Config{column: 0},
			expected: false,
		},
		{
			name:     "Single line",
			lines:    []string{"apple"},
			cfg:      Config{column: 0},
			expected: true,
		},
		{
			name:     "Empty slice",
			lines:    []string{},
			cfg:      Config{column: 0},
			expected: true,
		},
		{
			name:     "Sorted numbers",
			lines:    []string{"1", "2", "3"},
			cfg:      Config{column: 0, numeric: true},
			expected: true,
		},
		{
			name:     "Unsorted numbers",
			lines:    []string{"3", "1", "2"},
			cfg:      Config{column: 0, numeric: true},
			expected: false,
		},
		{
			name:     "Reverse sorted with reverse flag",
			lines:    []string{"cherry", "banana", "apple"},
			cfg:      Config{column: 0, reverse: true},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsSorted(tt.lines, tt.cfg)
			if result != tt.expected {
				t.Errorf("IsSorted(%v, %+v) = %t, expected %t",
					tt.lines, tt.cfg, result, tt.expected)
			}
		})
	}
}

func TestProcessFlags(t *testing.T) {

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	tests := []struct {
		name           string
		args           []string
		expectedConfig Config
		expectedFile   string
	}{
		{
			name:           "Default flags",
			args:           []string{"cmd", "file.txt"},
			expectedConfig: Config{},
			expectedFile:   "file.txt",
		},
		{
			name: "Numeric sort",
			args: []string{"cmd", "-n", "file.txt"},
			expectedConfig: Config{
				column:  0,
				numeric: true,
			},
			expectedFile: "file.txt",
		},
		{
			name: "Reverse sort",
			args: []string{"cmd", "-r", "file.txt"},
			expectedConfig: Config{
				column:  0,
				reverse: true,
			},
			expectedFile: "file.txt",
		},
		{
			name: "Column sort",
			args: []string{"cmd", "-k", "2", "file.txt"},
			expectedConfig: Config{
				column: 2,
			},
			expectedFile: "file.txt",
		},
		{
			name: "Multiple flags",
			args: []string{"cmd", "-n", "-r", "-u", "-k", "3", "file.txt"},
			expectedConfig: Config{
				column:  3,
				numeric: true,
				reverse: true,
				unique:  true,
			},
			expectedFile: "file.txt",
		},
		{
			name: "Month sort",
			args: []string{"cmd", "-M", "file.txt"},
			expectedConfig: Config{
				column: 0,
				month:  true,
			},
			expectedFile: "file.txt",
		},
		{
			name: "Human numeric sort",
			args: []string{"cmd", "-h", "file.txt"},
			expectedConfig: Config{
				column:       0,
				humanNumeric: true,
			},
			expectedFile: "file.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			flag.CommandLine = flag.NewFlagSet(tt.args[0], flag.ExitOnError)
			os.Args = tt.args

			cfg, filename := processFlags()

			if cfg.column != tt.expectedConfig.column ||
				cfg.numeric != tt.expectedConfig.numeric ||
				cfg.reverse != tt.expectedConfig.reverse ||
				cfg.unique != tt.expectedConfig.unique ||
				cfg.month != tt.expectedConfig.month ||
				cfg.ignoreBlanks != tt.expectedConfig.ignoreBlanks ||
				cfg.checkSorted != tt.expectedConfig.checkSorted ||
				cfg.humanNumeric != tt.expectedConfig.humanNumeric {
				t.Errorf("processFlags() config = %+v, expected %+v", cfg, tt.expectedConfig)
			}

			if filename != tt.expectedFile {
				t.Errorf("processFlags() filename = %s, expected %s", filename, tt.expectedFile)
			}
		})
	}
}
