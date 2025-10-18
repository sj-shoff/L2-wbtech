package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestParseFields(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []fieldRange
		wantErr bool
	}{
		{
			name:  "single field",
			input: "1",
			want:  []fieldRange{{1, 1}},
		},
		{
			name:  "multiple fields",
			input: "1,3,5",
			want:  []fieldRange{{1, 1}, {3, 3}, {5, 5}},
		},
		{
			name:  "range",
			input: "1-3",
			want:  []fieldRange{{1, 3}},
		},
		{
			name:  "mixed fields and ranges",
			input: "1,3-5,7",
			want:  []fieldRange{{1, 1}, {3, 5}, {7, 7}},
		},
		{
			name:    "invalid field",
			input:   "abc",
			wantErr: true,
		},
		{
			name:    "invalid range",
			input:   "1-",
			wantErr: true,
		},
		{
			name:    "decreasing range",
			input:   "5-3",
			wantErr: true,
		},
		{
			name:    "zero field",
			input:   "0",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseFields(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseFields() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !equalFieldRanges(got, tt.want) {
				t.Errorf("parseFields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func equalFieldRanges(a, b []fieldRange) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestRunCut(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		fields    string
		delim     string
		separated bool
		want      string
		wantErr   bool
	}{
		{
			name:   "single field",
			input:  "a\tb\tc\td\n",
			fields: "2",
			delim:  "\t",
			want:   "b\n",
		},
		{
			name:   "multiple fields",
			input:  "a\tb\tc\td\n",
			fields: "1,3",
			delim:  "\t",
			want:   "a\tc\n",
		},
		{
			name:   "range of fields",
			input:  "a\tb\tc\td\te\n",
			fields: "2-4",
			delim:  "\t",
			want:   "b\tc\td\n",
		},
		{
			name:   "mixed fields and ranges",
			input:  "a\tb\tc\td\te\n",
			fields: "1,3-4",
			delim:  "\t",
			want:   "a\tc\td\n",
		},
		{
			name:      "separated flag skips lines without delimiter",
			input:     "a\tb\tc\nno_delimiter\nd\te\tf\n",
			fields:    "1",
			delim:     "\t",
			separated: true,
			want:      "a\nd\n",
		},
		{
			name:   "fields beyond boundary are ignored",
			input:  "a\tb\tc\n",
			fields: "1,5",
			delim:  "\t",
			want:   "a\n",
		},
		{
			name:   "custom delimiter",
			input:  "a,b,c,d\n",
			fields: "2-3",
			delim:  ",",
			want:   "b,c\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := strings.NewReader(tt.input)
			var output bytes.Buffer

			cfg := config{
				fields:    tt.fields,
				delimiter: tt.delim,
				separated: tt.separated,
			}

			err := runCut(input, &output, cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("runCut() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got := output.String()
			if got != tt.want {
				t.Errorf("runCut() = %q, want %q", got, tt.want)
			}
		})
	}
}
