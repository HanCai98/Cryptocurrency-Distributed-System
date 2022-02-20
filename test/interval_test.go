package test

import (
	liteminer "liteminer/pkg"
	"testing"
)

func TestInterval1(t *testing.T) {
	upperBound := uint64(20)
	numIntervals := 5
	intervals := liteminer.GenerateIntervals(upperBound, numIntervals)

	target := make([]liteminer.Interval, numIntervals)
	target[0] = liteminer.Interval{Lower: 0, Upper: 4}
	target[1] = liteminer.Interval{Lower: 4, Upper: 8}
	target[2] = liteminer.Interval{Lower: 8, Upper: 12}
	target[3] = liteminer.Interval{Lower: 12, Upper: 16}
	target[4] = liteminer.Interval{Lower: 16, Upper: 20}

	if len(intervals) != len(target) {
		t.Errorf("Expected  %d, but received %d", target, intervals)
	}

	for i, interval := range target {
		if interval != intervals[i] {
			t.Errorf("Expected  %d, but received %d", target, intervals)
		}
	}

}

func TestInterval2(t *testing.T) {
	upperBound := uint64(22)
	numIntervals := 5
	intervals := liteminer.GenerateIntervals(upperBound, numIntervals)

	target := make([]liteminer.Interval, numIntervals)
	target[0] = liteminer.Interval{Lower: 0, Upper: 5}
	target[1] = liteminer.Interval{Lower: 5, Upper: 10}
	target[2] = liteminer.Interval{Lower: 10, Upper: 14}
	target[3] = liteminer.Interval{Lower: 14, Upper: 18}
	target[4] = liteminer.Interval{Lower: 18, Upper: 22}

	if len(intervals) != len(target) {
		t.Errorf("Expected  %d, but received %d", target, intervals)
	}

	for i, interval := range target {
		if interval != intervals[i] {
			t.Errorf("Expected  %d, but received %d", target, intervals)
		}
	}

}

func TestInterval3(t *testing.T) {
	upperBound := uint64(5)
	numIntervals := 6
	intervals := liteminer.GenerateIntervals(upperBound, numIntervals)

	target := make([]liteminer.Interval, upperBound)
	target[0] = liteminer.Interval{Lower: 0, Upper: 1}
	target[1] = liteminer.Interval{Lower: 1, Upper: 2}
	target[2] = liteminer.Interval{Lower: 2, Upper: 3}
	target[3] = liteminer.Interval{Lower: 3, Upper: 4}
	target[4] = liteminer.Interval{Lower: 4, Upper: 5}

	if len(intervals) != len(target) {
		t.Errorf("Expected  %d, but received %d", target, intervals)
	}

	for i, interval := range target {
		if interval != intervals[i] {
			t.Errorf("Expected  %d, but received %d", target, intervals)
		}
	}
}

func TestInterval4(t *testing.T) {
	upperBound := uint64(10)
	numIntervals := -2
	intervals := liteminer.GenerateIntervals(upperBound, numIntervals)

	var target []liteminer.Interval
	target = append(target, liteminer.Interval{Lower: 0, Upper: 10})

	if len(intervals) != len(target) {
		t.Errorf("Expected  %d, but received %d", target, intervals)
	}

	for i, interval := range target {
		if interval != intervals[i] {
			t.Errorf("Expected  %d, but received %d", target, intervals)
		}
	}
}
