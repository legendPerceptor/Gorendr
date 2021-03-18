package rendr

import (
	"math"
	"testing"
)

func TestM4ai(t *testing.T) {
	_mat := [16]float64{
		-0.846443892, -0.268435329, -0.728768885, 1.74738681,
		0.351376086, -0.246699125, -0.860544026, -1.32462478,
		1.23825359, -0.530447423, -0.0935271308, 1.33915234,
		0, 0, 0, 1,
	}
	rnd_m4_print(_mat, nil)
	inv, _ := rnd_m4_affine_inv(_mat)
	invStd := [16]float64{
		-0.777712464,   0.648633242,   0.0919002295,    2.09509182,
		-1.85313749,    1.76136518,   -1.76658344,     7.9370203,
		0.21369946,  -1.40215135,   0.543965399,   -2.95919275,
		0,             0,             0,             1,
	}
	rnd_m4_print(inv, nil)
	for i:=0; i<16; i++ {
		if math.Abs(inv[i]-invStd[i])>0.0001 {
			t.Errorf("The Inverse of matrix is not correct")
		}
	}
}