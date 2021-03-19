package rendr

import (
	"gonum.org/v1/gonum/mat"
	"math"
)

func rndConvoNew(ctx *rndCtx) *rndConvo{
	support := ctx.Kern.Support
	cnv := rndConvo{
		Ipos:         [3]float64{},
		Value:        0,
		Gradient:     [3]float64{},
		Inside:       0,
		Kern_values:  make([]float64, 3*support),
		Deriv_values: make([]float64, 3*support),
	}
	return &cnv
}

func rndConvoEval(xw, yw, zw float64, cnv *rndConvo, ctx *rndCtx) {
	input:= [4]float64{xw, yw, zw, 1}
	var output [4]float64
	inVec := mat.NewVecDense(4, input[:])
	outVec:= mat.NewVecDense(4, output[:])
	wtoi := mat.NewDense(4,4,ctx.WtoI[:])
	outVec.MulVec(wtoi, inVec)
	cnv.Ipos = [3]float64{output[0], output[1], output[2]}
	var n1,n2,n3,start,end int
	var alpha1, alpha2, alpha3 float64
	if ctx.Kern.Support %2==0 {
		n1 = int(math.Floor(output[0]))
		n2 = int(math.Floor(output[1]))
		n3 = int(math.Floor(output[2]))
		start = 1-ctx.Kern.Support/2
		end = ctx.Kern.Support/2
	} else {
		n1 = int(math.Floor(output[0]+0.5))
		n2 = int(math.Floor(output[1]+0.5))
		n3 = int(math.Floor(output[2]+0.5))
		start = (1-ctx.Kern.Support)/2
		end = (ctx.Kern.Support-1)/2
	}
	alpha1 = output[0]-float64(n1)
	alpha2 = output[1]-float64(n2)
	alpha3 = output[2]-float64(n3)
	var i,j,k int // loop variables
	outside:=0
	if n1+start <0 || uint(n1+end) >=ctx.Vol.Size[0] {
		outside = outside + 1
	} else if n2+start<0 || uint(n2+end) >=ctx.Vol.Size[1] {
		outside = outside + 1
	} else if n3+start<0 || uint(n3+end) >= ctx.Vol.Size[2] {
		outside = outside + 1
	}
	if outside>0 {
		cnv.Value = ctx.OutsideValue
		cnv.Gradient = [3]float64{ctx.OutsideValue, ctx.OutsideValue, ctx.OutsideValue}
		cnv.Inside = 0
		return
	}
	cnv.Inside = 1
	result:=float64(0)
	support:= ctx.Kern.Support
	k1_val := cnv.Kern_values[:support]
	k2_val := cnv.Kern_values[support:2*support]
	k3_val := cnv.Kern_values[2*support:]

	ctx.Kern.Apply(k1_val, alpha1)
	ctx.Kern.Apply(k2_val, alpha2)
	ctx.Kern.Apply(k3_val, alpha3)
	volSize := ctx.Vol.Size
	for i=start; i<=end; i++ {
		for j=start; j<=end; j++ {
			for k=start; k<=end; k++ {
				dTmp:=ctx.Vol.Data[uint(n3+i)*volSize[1]*volSize[0]+uint(n2+j)*volSize[0]+uint(n1+k)]
				result += dTmp*k1_val[k-start] * k2_val[j-start] * k3_val[i-start]
			}
		}
	}
	cnv.Value = result
	//Deriv
	k1_deriv:= cnv.Deriv_values[:support]
	k2_deriv:= cnv.Deriv_values[support:2*support]
	k3_deriv := cnv.Deriv_values[2*support:]
	ctx.Kern.Deriv.Apply(k1_deriv, alpha1)
	ctx.Kern.Deriv.Apply(k2_deriv, alpha2)
	ctx.Kern.Deriv.Apply(k3_deriv, alpha3)
	var tmp_gradient [3]float64
	for i=start; i<=end; i++ {
		for j = start; j <= end; j++ {
			for k = start; k <= end; k++ {
				dTmp:=ctx.Vol.Data[uint(n3+i)*volSize[1]*volSize[0]+uint(n2+j)*volSize[0]+uint(n1+k)]
				tmp_gradient[0] += dTmp*k1_deriv[k-start]*k2_val[j-start]*k3_val[i-start]
				tmp_gradient[1] += dTmp*k1_val[k-start]*k2_deriv[j-start]*k3_val[i-start]
				tmp_gradient[2] += dTmp*k1_val[k-start]*k2_val[j-start]*k3_deriv[i-start]
			}
		}
	}
	// Conver from index space to world space
	tgradVec := mat.NewVecDense(3, tmp_gradient[:])
	gradVec := mat.NewVecDense(3, cnv.Gradient[:])
	mt := mat.NewDense(3, 3, ctx.MT[:])
	gradVec.MulVec(mt, tgradVec)
}