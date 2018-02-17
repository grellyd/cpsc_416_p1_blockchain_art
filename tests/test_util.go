package tests

import (
	"blockartlib"
	"minerlib"
	"strings"
	"testing"
)

type CanvasSettings = blockartlib.CanvasSettings
type LineSegment = minerlib.LineSegment
type Operation = blockartlib.Operation
type Point = minerlib.Point
type Shape = minerlib.Shape

func ConcatOps(ops map[string]*Operation) string {
	opSigs := ""
	for _, op := range ops {
		opSigs += op.ShapeHash + " by " + op.ArtNodePubKey + ", "
	}
	return opSigs
}

func ExpectEquals(t *testing.T, s, expected string) {
	if strings.Compare(s, expected) != 0 {
		t.Errorf("Got: %s\nExpected: %s\n", s, expected)
	}
}

func ExpectContains(t *testing.T, testStr string, strArr []string) {
	for _, str := range strArr {
		if !strings.Contains(testStr, str) {
			t.Errorf("Expected string to contain %v, but not found.\n", str)
		}
	}
}

func ExpectTrue(t *testing.T, b bool, errStr string) {
	if !b {
		t.Errorf(errStr + "Expected true. Got false.\n")
	}
}

func ExpectFalse(t *testing.T, b bool, errStr string) {
	if b {
		t.Errorf(errStr + "Expected false. Got true.\n")
	}
}
