package types

import (
	"reflect"
	"testing"
)

func TestCompile(t *testing.T) {
	tests := []struct {
		name        string
		pattern     string
		wantErr     bool
		matchTest   string
		shouldMatch bool
	}{
		{
			name:        "valid regex",
			pattern:     "^test[0-9]+$",
			wantErr:     false,
			matchTest:   "test123",
			shouldMatch: true,
		},
		{
			name:        "valid regex no match",
			pattern:     "^test[0-9]+$",
			wantErr:     false,
			matchTest:   "test",
			shouldMatch: false,
		},
		{
			name:    "invalid regex",
			pattern: "[",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Compile(tt.pattern)
			if (err != nil) != tt.wantErr {
				t.Errorf("Compile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if match := got.MatchString(tt.matchTest); match != tt.shouldMatch {
				t.Errorf("Regex.MatchString() = %v, want %v", match, tt.shouldMatch)
			}
		})
	}
}

func TestMustCompile(t *testing.T) {
	tests := []struct {
		name      string
		pattern   string
		wantPanic bool
	}{
		{
			name:      "valid regex",
			pattern:   "^test[0-9]+$",
			wantPanic: false,
		},
		{
			name:      "invalid regex",
			pattern:   "[",
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("MustCompile() panic = %v, wantPanic %v", r, tt.wantPanic)
				}
			}()
			got := MustCompile(tt.pattern)
			if tt.wantPanic {
				t.Error("Expected panic did not occur")
				return
			}
			if got == nil {
				t.Error("MustCompile() returned nil for valid pattern")
			}
		})
	}
}

func TestRegex_MarshalText(t *testing.T) {
	pattern := "^test[0-9]+$"
	re := MustCompile(pattern)

	got, err := re.MarshalText()
	if err != nil {
		t.Errorf("Regex.MarshalText() error = %v", err)
		return
	}
	if string(got) != pattern {
		t.Errorf("Regex.MarshalText() = %v, want %v", string(got), pattern)
	}
}

func TestRegex_UnmarshalText(t *testing.T) {
	tests := []struct {
		name    string
		text    []byte
		wantErr bool
	}{
		{
			name:    "valid regex",
			text:    []byte("^test[0-9]+$"),
			wantErr: false,
		},
		{
			name:    "invalid regex",
			text:    []byte("["),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			re := &Regex{}
			err := re.UnmarshalText(tt.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("Regex.UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			// Verify the regex was properly unmarshaled
			if re.String() != string(tt.text) {
				t.Errorf("Regex.String() = %v, want %v", re.String(), string(tt.text))
			}
		})
	}
}

func TestRegex_MarshalUnmarshalRoundTrip(t *testing.T) {
	pattern := "^test[0-9]+$"
	original := MustCompile(pattern)

	// Marshal
	text, err := original.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText() error = %v", err)
	}

	// Unmarshal into new regex
	var unmarshaled Regex
	if err := unmarshaled.UnmarshalText(text); err != nil {
		t.Fatalf("UnmarshalText() error = %v", err)
	}

	// Compare the patterns
	if !reflect.DeepEqual(original.String(), unmarshaled.String()) {
		t.Errorf("Round trip failed: got %v, want %v", unmarshaled.String(), original.String())
	}

	// Test actual regex behavior
	testStr := "test123"
	if original.MatchString(testStr) != unmarshaled.MatchString(testStr) {
		t.Error("Round-tripped regex behaves differently from original")
	}
}
