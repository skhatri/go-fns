package types

import (
	"regexp"
)

type Regex struct {
	*regexp.Regexp
}

func Compile(expr string) (*Regex, error) {
	reg, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	return &Regex{reg}, nil
}

func (re *Regex) MarshalText() ([]byte, error) {
	return []byte(re.String()), nil
}

func (re *Regex) UnmarshalText(text []byte) error {
	newRE, err := Compile(string(text))
	if err != nil {
		return err
	}
	*re = *newRE
	return nil
}

