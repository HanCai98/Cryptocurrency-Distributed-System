/*
 *  Brown University, CS138, Spring 2022
 *
 *  Purpose: contains all interval related logic.
 */

package pkg

// Interval represents [Lower, Upper)
type Interval struct {
	Lower uint64 // Inclusive
	Upper uint64 // Exclusive
}

// GenerateIntervals divides the range [0, upperBound] into numIntervals intervals.
func GenerateIntervals(upperBound uint64, numIntervals int) (intervals []Interval) {
	// TODO: Students should implement this.
	// create a slice
	if numIntervals <= 0 {
		intervals = append(intervals, Interval{Lower: 0, Upper: upperBound})
		return intervals
	}

	if numIntervals > int(upperBound) {
		numIntervals = int(upperBound)
	}

	intervals = make([]Interval, numIntervals)
	remainder := upperBound % uint64(numIntervals)
	stepSize := upperBound / uint64(numIntervals)
	index := 0

	for i := uint64(0); i < upperBound; i += stepSize {
		lower := i
		upper := i + stepSize

		if remainder != 0 {
			upper += 1
			remainder -= 1
			i += 1
		}

		intervals[index] = Interval{Lower: lower, Upper: upper}
		index += 1

	}

	return intervals
}
