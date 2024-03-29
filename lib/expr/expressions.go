package expr

import (
	"os"
	"regexp"
)

var exprRegex = regexp.MustCompile(`^\$\{(.+)\}$`)
var equation = regexp.MustCompile(`(env\.)?([a-zA-Z_0-9]+)\s*([=!]=)\s*([^ ]+)`)

func SolveEnvExpression(expr string) bool {
	out := exprRegex.FindStringSubmatch(expr)
	if len(out) == 1 {
		return false
	}
	condition := out[1]
	equationComponents := equation.FindStringSubmatch(condition)
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
