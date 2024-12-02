package util

func FindUniqueElemens(arr []uint32) []uint32 {
	m := map[uint32]struct{}{}
	result := []uint32{}
	for _, v := range arr {
		if _, ok := m[v]; !ok {
			m[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

func FindCommonElements(a [][]uint32) []uint32 {
	result := []uint32{}
	index := make([]uint32, len(a))
	for {
		isEnd := false
		checkIndexList := []uint32{}
		for i := 0; i < len(a); i++ {
			if len(a[i]) == 0 {
				return []uint32{}
			}
			checkIndexList = append(checkIndexList, a[i][index[i]])
		}
		isEqual := allElementsEqual(checkIndexList)
		if !isEqual {
			maxNum, _ := maxElementWithIndex(checkIndexList)
			for i := 0; i < len(index); i++ {
				index[i] = binarySearch(a[i], maxNum)
				checkIndexList[i] = a[i][index[i]]
			}

			if allElementsEqual(checkIndexList) {
				isEqual = true
			} else {
				_, minNumIndex := minElementWithIndex(checkIndexList)
				index[minNumIndex]++
				if index[minNumIndex] >= uint32(len(a[minNumIndex])) {
					isEnd = true
				}
			}
		}
		if isEqual {
			result = append(result, checkIndexList[0])
			for i := 0; i < len(index); i++ {
				index[i]++
				if index[i] >= uint32(len(a[i])) {
					isEnd = true
				}
			}
		}
		if isEnd {
			break
		}
	}
	return result
}

func maxElementWithIndex(arr []uint32) (uint32, uint32) {
	max := arr[0]
	index := 0
	for i, v := range arr {
		if v > max {
			max = v
			index = i
		}
	}
	return max, uint32(index)
}

func minElementWithIndex(arr []uint32) (uint32, uint32) {
	min := arr[0]
	index := 0
	for i, v := range arr {
		if v < min {
			min = v
			index = i
		}
	}
	return min, uint32(index)
}

func allElementsEqual[T comparable](arr []T) bool {
	if len(arr) == 0 {
		return true
	}
	first := arr[0]
	for _, v := range arr {
		if v != first {
			return false
		}
	}
	return true
}

func binarySearch(arr []uint32, target uint32) uint32 {
	low := 0
	high := len(arr) - 1
	for low <= high {
		mid := (low + high) / 2
		if arr[mid] == target {
			return uint32(mid)
		} else if arr[mid] < target {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return uint32(low)
}
