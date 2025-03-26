// Package expr provides utilities for evaluating expressions, particularly environment variable
// expressions in strings.
package expr

import (
	"os"
	"regexp"
)

var exprRegex = regexp.MustCompile(`^\$\{(.+)\}$`)
var equation = regexp.MustCompile(`(env\.)?([a-zA-Z_0-9]+)\s*([=!]=)\s*([^ ]*)(?:\s*)$`)

// SolveEnvExpression evaluates a string containing environment variable expressions
// in the format ${VARIABLE_NAME} and replaces them with their corresponding values.
// If a variable is not found in the environment, it will be replaced with an empty string.
// Returns the evaluated string with all environment variables resolved.
func SolveEnvExpression(expr string) bool {
	out := exprRegex.FindStringSubmatch(expr)
	if len(out) < 2 {
		return false
	}
	condition := out[1]
	equationComponents := equation.FindStringSubmatch(condition)
	if len(equationComponents) < 5 {
		return false
	}
	isEnvVar := equationComponents[1] == "env."
	lhs := equationComponents[2]
	if isEnvVar {
		lhs = os.Getenv(lhs)
	}
	op := equationComponents[3]
	rhs := equationComponents[4]
	switch op {
	case "==":
		return lhs == rhs
	default:
		return lhs != rhs
	}
}
