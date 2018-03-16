package collatz

import (
	"sort"
	"time"

	"github.com/fatih/stopwatch"
)

//Collatz returns the maximum Collatz count for all numbers 1 through n and
//the time it took to perform this calculation. A parameter specifies
//whether to calculate this value using concurrent goroutines.
func Collatz(n int, concurrent bool) (int, time.Duration) {
	if n < 1 {
		return -1, 0
	}

	/* if true {
		time.Sleep(1000 * time.Millisecond)
		return 1, 0
	} */

	sw := stopwatch.Start(0)
	numSlice := []int{}
	const partitionSize = 10000

	if concurrent {
		values := make(chan int, 1000)
		chunkCalls := 0
		for i := n; i >= 1; i -= partitionSize {
			startIndex := i - partitionSize + 1
			if startIndex < 1 {
				startIndex = 1
			}
			go chunkConcurrent(startIndex, i, values)
			chunkCalls++
		}
		for i := 0; i < chunkCalls; i++ {
			numSlice = append(numSlice, <-values)
		}
	} else {
		for i := 1; i <= n; i += partitionSize {
			endIndex := i + partitionSize - 1
			if endIndex > n {
				endIndex = n
			}
			localMax := chunk(i, endIndex)
			numSlice = append(numSlice, localMax)
		}
	}

	sort.Slice(numSlice, func(i, j int) bool {
		return numSlice[i] > numSlice[j]
	})

	sw.Stop()
	return numSlice[0], sw.ElapsedTime()
}

func chunk(startIndex, endIndex int) int {
	numSlice := []int{}
	for i := startIndex; i <= endIndex; i++ {
		numSlice = append(numSlice, loop(i))
	}
	sort.Slice(numSlice, func(i, j int) bool {
		return numSlice[i] > numSlice[j]
	})
	return numSlice[0]
}

func chunkConcurrent(startIndex, endIndex int, values chan<- int) {
	numSlice := []int{}
	for i := startIndex; i <= endIndex; i++ {
		numSlice = append(numSlice, loop(i))
	}
	sort.Slice(numSlice, func(i, j int) bool {
		return numSlice[i] > numSlice[j]
	})
	values <- numSlice[0]
}

func loop(n int) int {
	count := 0
	for n > 1 {
		if n%2 == 0 {
			n = n / 2
		} else {
			n = (3 * n) + 1
		}
		count++
	}
	return count
}
