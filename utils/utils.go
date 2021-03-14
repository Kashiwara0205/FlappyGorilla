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

func GetRotateValue(values []int, i int) int{
	length := len(values)
	x := (length + i) / length
	idx := length + i - length * x 

	return values[idx]
}