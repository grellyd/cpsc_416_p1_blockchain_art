package testlib

import (
	"minerlib"
	"testing"
)

type LineSegment = minerlib.LineSegment
type Point = minerlib.Point

func TestIntersectionsVertical(t *testing.T) {
	l1 := LineSegment{Point{5, 0}, Point{5, 5}}
	l2 := LineSegment{Point{10, 0}, Point{10, 5}}
	l3 := LineSegment{Point{5, 4}, Point{5, 7}}

	t1 := minerlib.IsIntersecting(l1, l2) // Expect false
	if t1 {
		t.Errorf("intersect(l1, l2). Expected false, got %v\n", t1)
	}

	t2 := minerlib.IsIntersecting(l1, l3) // Expect true
	if !t2 {
		t.Errorf("intersect(l1, l3). Expected false, got %v\n", t2)
	}
}
