package main

func main() {
	var x []int
	y := []int{1, 2, 3, 4, 5}
	x = append(x, y...)
}
