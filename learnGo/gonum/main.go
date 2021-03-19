package main

import "fmt"

func tryme(_inv []float64) {
	_inv[1] = 100
}

func main() {
	//u := mat.NewVecDense(3, []float64{1, 2, 3})
	//v := mat.NewVecDense(3, []float64{4, 5, 6})
	a := [3]float64{4,2,1}
	tryme(a[:])
	fmt.Println(a)
}
