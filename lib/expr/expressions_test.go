package expr

import (
	"os"
	"testing"
)

func TestSolveEnableExpr(t *testing.T) {
	// Test environment variable cases
	switch os.Getenv("GOOS") {
	case "darwin":
		assertTrue(t, SolveEnvExpression("${env.GOOS==darwin}"), "env.GOOS==darwin should hold true when running on darwin")
	default:
		assertTrue(t, SolveEnvExpression("${env.GOOS!=darwin}"), "env.GOOS!=darwin should hold true when not running on darwin")
	}

	// Test non-env variable cases
	assertTrue(t, SolveEnvExpression("${on==on}"), "text match expected to be true when lhs and rhs is same")
	assertTrue(t, SolveEnvExpression("${on!=off}"), "text match expected to be false when lhs and rhs are different")

	// Test invalid expression format
	assertFalse(t, SolveEnvExpression("invalid"), "invalid expression format should return false")
	assertFalse(t, SolveEnvExpression("${invalid}"), "invalid equation format should return false")
	assertFalse(t, SolveEnvExpression("${env.INVALID_VAR==value}"), "non-existent env var should return false")

	// Test with empty environment variable
	os.Setenv("EMPTY_VAR", "")
	assertTrue(t, SolveEnvExpression("${env.EMPTY_VAR==}"), "empty env var should match empty string")
	assertTrue(t, SolveEnvExpression("${env.EMPTY_VAR!=nonempty}"), "empty env var should not match non-empty string")
}

func assertTrue(t *testing.T, cond bool, msg string) {
	if !cond {
		t.Errorf(msg)
	}
}

func assertFalse(t *testing.T, cond bool, msg string) {
	if cond {
		t.Errorf(msg)
	}
}
