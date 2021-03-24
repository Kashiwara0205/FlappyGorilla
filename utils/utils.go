package utils

func FloorDiv(x, y int) int {
	d := x / y
	if d*y == x || x >= 0 {
		return d
	}
	return d - 1
}

func FloorMod(x, y int) int {
	return x - FloorDiv(x, y)*y
}

func GetRotateIndex(size int, i int) int{
	x := (size + i) / size
	idx := size + i - size * x 

	return idx
}

func GetRotateValue(values []int, i int) int{
	return values[GetRotateIndex(len(values), i)]
}