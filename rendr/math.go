package rendr

import (
	"fmt"
	"gonum.org/v1/gonum/mat"
	"os"
	"testing"
)

/*
  rnd_m4_affine_inv: invert a given 4x4 matrix "mat", and store inverse in
  "inv", assuming that the bottom row of "mat" is [0 0 0 1].  This means the
  matrix represents an affine transform, and its inverse is simpler to
  compute than for a general 4x4 matrix.
*/
func rnd_m4_affine_inv(_mat [16]float64) ([16]float64, [9]float64) {
	upper := [9]float64{}
	copy(upper[0:3], _mat[0:3])
	copy(upper[3:6], _mat[4:7])
	copy(upper[6:9], _mat[8:11])
	inv3 := mat.NewDense(3, 3, upper[:])
	if err:= inv3.Inverse(inv3); err!=nil {
		fmt.Println("Cannot inverse the matrix")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	_V := [3]float64{-_mat[3], -_mat[7], -_mat[11]}
	V := mat.NewVecDense(3, _V[:])
	//var U mat.VecDense
	actual := make([]float64,3)
	U:= mat.NewVecDense(3, actual)
	U.MulVec(inv3, V)
	_inv := [16]float64{
		upper[0], upper[1], upper[2], actual[0],
		upper[3], upper[4], upper[5], actual[1],
		upper[6], upper[7], upper[8], actual[2],
		0, 0, 0, 1,
	}
	return _inv,upper
}

func rnd_m4_print(_mat [16]float64, t *testing.T) {
	s:= fmt.Sprintf("mat=[%.5f, %.5f, %.5f, %.5f\n    %.5f, %.5f, %.5f, %.5f\n    %.5f, %.5f, %.5f, %.5f\n    %.5f, %.5f, %.5f, %.5f]\n",
		_mat[0], _mat[1], _mat[2], _mat[3],
		_mat[4], _mat[5], _mat[6], _mat[7],
		_mat[8], _mat[9], _mat[10], _mat[11],
		_mat[12], _mat[13], _mat[14], _mat[15])
	if t==nil {
		fmt.Printf(s)
	} else {
		t.Errorf(s)
	}
}