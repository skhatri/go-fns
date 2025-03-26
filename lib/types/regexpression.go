// Package types provides custom types and their implementations, including
// specialized types for regular expressions and other data structures.
package types

import (
	"regexp"
)

// Regex represents a regular expression pattern with methods for
// compilation and text marshaling/unmarshaling.
type Regex struct {
	*regexp.Regexp
}

// Compile creates a new RegexExpression from a pattern string.
// Returns the compiled RegexExpression and any error that occurred during compilation.
func Compile(expr string) (*Regex, error) {
	reg, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	return &Regex{reg}, nil
}

// MarshalText implements the encoding.TextMarshaler interface for RegexExpression.
// It marshals the regular expression pattern to a text representation.
func (re *Regex) MarshalText() ([]byte, error) {
	return []byte(re.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for RegexExpression.
// It unmarshals a text representation into a regular expression pattern.
func (re *Regex) UnmarshalText(text []byte) error {
	newRE, err := Compile(string(text))
	if err != nil {
		return err
	}
	*re = *newRE
	return nil
}

// MustCompile creates a new RegexExpression from a pattern string.
// It panics if the pattern cannot be compiled.
// This is a convenience function for use when the pattern is known to be valid.
func MustCompile(pattern string) *Regex {
	reg, err := Compile(pattern)
	if err != nil {
		panic(err)
	}
	return reg
}
