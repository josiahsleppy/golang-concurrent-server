package collatz

import "sort"

//Collatz returns the maximum Collatz count for all numbers 1 through n.
//An secondary parameter specifies whether partitioning should be used in the calculations.
func Collatz(n int, useChunk bool) int {
	if n < 1 {
		return -1
	}
	var retVal int
	if useChunk {
		const partitionSize = 1000
		array := []int{}
		for i := 1; i <= n; i += partitionSize {
			endIndex := i + partitionSize - 1
			if endIndex > n {
				endIndex = n
			}
			localMax := chunk(i, endIndex)
			go chunk(i, endIndex)
			array = append(array, localMax)
		}
		//sort
		sort.Slice(array, func(i, j int) bool {
			return array[i] > array[j]
		})
		retVal = array[0]
	} else {
		retVal = chunk(1, n)
	}
	return retVal
}

func chunk(startIndex, endIndex int) int {
	array := []int{}
	for i := startIndex; i <= endIndex; i++ {
		count := loop(i)
		array = append(array, count)
	}
	//sort
	sort.Slice(array, func(i, j int) bool {
		return array[i] > array[j]
	})
	return array[0]
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
