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
	maxNum := a[0][index[0]]
	i := 0
	sameCnt := 1
loop:
	for {
		tmp, ok := binarySearch(a[i][index[i]:], maxNum)
		if ok {
			index[i] = tmp + index[i]
			sameCnt++
			if sameCnt == len(index) {
				result = append(result, a[i][index[i]])
				for j := 0; j < len(index); j++ {
					index[j]++
					if index[j] >= uint32(len(a[j])) {
						break loop
					}
					if a[j][index[j]] > maxNum {
						maxNum = a[j][index[j]]
					}
				}
				sameCnt = 1
			}
		} else {
			index[i] = tmp + index[i] + 1
			if index[i] >= uint32(len(a[i])) {
				break loop
			}
			if a[i][index[i]] > maxNum {
				maxNum = a[i][index[i]]
			}
			sameCnt = 1
		}
		i = (i + 1) % len(a)
	}
	return result
}

func binarySearch(arr []uint32, target uint32) (uint32, bool) {
	low, high := 0, len(arr)-1
	for low <= high {
		mid := (low + high) / 2
		if arr[mid] == target {
			return uint32(mid), true
		} else if arr[mid] < target {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return uint32(high), false
}
