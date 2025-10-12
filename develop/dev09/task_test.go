package main

import "testing"

func TestUnpack(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "basic unpack",
			input:   "a4bc2d5e",
			want:    "aaaabccddddde",
			wantErr: false,
		},
		{
			name:    "no digits",
			input:   "abcd",
			want:    "abcd",
			wantErr: false,
		},
		{
			name:    "empty string",
			input:   "",
			want:    "",
			wantErr: false,
		},
		{
			name:    "only digits",
			input:   "45",
			want:    "",
			wantErr: true,
		},
		{
			name:    "starts with digit",
			input:   "4a",
			want:    "",
			wantErr: true,
		},
		{
			name:    "escaped characters",
			input:   `qwe\4\5`,
			want:    "qwe45",
			wantErr: false,
		},
		{
			name:    "mixed escape and digits",
			input:   `qwe\45`,
			want:    "qwe44444",
			wantErr: false,
		},
		{
			name:    "escape at end",
			input:   `qwe\`,
			want:    "",
			wantErr: true,
		},
		{
			name:    "zero digit",
			input:   "a0",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Unpack(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unpack() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Unpack() = %v, want %v", got, tt.want)
			}
		})
	}
}
