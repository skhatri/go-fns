package expr

import (
	"os"
	"testing"
)

func TestSolveEnableExpr(t *testing.T) {

	switch os.Getenv("GOOS") {
	case "darwin":
		assertTrue(t, SolveEnvExpression("${env.GOOS==darwin}"), "env.GOOS==darwin should hold true when running on darwin")
	default:
		assertTrue(t, SolveEnvExpression("${env.GOOS!=darwin}"), "env.GOOS!=darwin should hold true when not running on darwin")
	}
	assertTrue(t, SolveEnvExpression("${on==on}"), "text match expected to be true when lhs and rhs is same")
	assertTrue(t, SolveEnvExpression("${on!=off}"), "text match expected to be false when lhs and rhs are different")
}

func assertTrue(t *testing.T, cond bool, msg string) {
	if !cond {
		t.Errorf(msg)
	}
}
