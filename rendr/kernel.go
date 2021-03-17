package rendr

import (
	"fmt"
	"math"
	"strings"
)

func ZeroEval(xx float64) float64 {
	return 0
}

func Zero1Apply(ww []float64, xa float64) {
	ww[0] = 0
}

func Zero2Apply(ww []float64, xa float64) {
	ww[0] = 0
	ww[1] = 0
}

func Zero3Apply(ww []float64, xa float64) {
	ww[0] = 0
	ww[1] = 0
	ww[2] = 0
}

func Zero4Apply(ww []float64, xa float64) {
	ww[0] = 0
	ww[1] = 0
	ww[2] = 0
	ww[3] = 0
}

func Zero5Apply(ww []float64, xa float64) {
	ww[0] = 0
	ww[1] = 0
	ww[2] = 0
	ww[3] = 0
	ww[4] = 0
}

func Zero6Apply(ww []float64, xa float64) {
	ww[0] = 0
	ww[1] = 0
	ww[2] = 0
	ww[3] = 0
	ww[4] = 0
	ww[5] = 0
}

// ============ Ctmr to the third derivative ============

func ctmrEval(x float64) float64 {
	var ret float64
	x = math.Abs(x)
	if x<1 {
		ret = 1+x*x*(-2.5 + x*1.5);
	} else if x<2 {
		x -= 1;
		ret = x*(-0.5 + x*(1-x/2))
	} else {
		ret = 0
	}
	return ret
}

func ctmrApply(ww []float64, xa float64) {
	ww[0] = -((xa-1)*(xa-1)*xa)/2
	ww[1] = (2 + xa*xa*(-5 + 3*xa))/2
	ww[2] = (xa + xa*xa*(4 - 3*xa))/2
	ww[3] = ((-1 + xa)*xa*xa)/2
}



func dctmrEval(_x float64) float64 {
	var ret float64
	x := math.Abs(_x)
	if x<1 {
		ret = x*(-5 + x*4.5)
	} else if x<2 {
		x -= 1
		ret = -0.5 + x * (2-x*1.5)
	} else {
		ret = 0
	}
	if _x<0 {
		return -ret
	} else {
		return ret
	}
}

func dctmrApply(ww []float64, xa float64) {
	ww[0] = (-1 + (4 - 3*xa)*xa)/2
	ww[1] = (xa*(-10 + 9*xa))/2
	ww[2] = -((-1 + xa)*(1 + 9*xa))/2
	ww[3] = (xa*(-2 + 3*xa))/2
}

func ddctmrEval(x float64) float64 {
	var ret float64
	x = math.Abs(x)
	if x<1 {
		ret = -5 + 9 * x
	} else if x<2 {
		x -= 1
		ret = 2 - 3*x
	} else {
		ret = 0
	}
	return ret
}

func ddctmrApply(ww []float64, xa float64) {
	ww[0] = 2 - 3*xa;
	ww[1] = -5 + 9*xa;
	ww[2] = 4 - 9*xa;
	ww[3] = -1 + 3*xa;
}

func dddctmrEval(_x float64) float64 {
	var ret float64
	x := math.Abs(_x)
	if x<1 {
		ret = 9
	} else if x<2 {
		ret = -3
	} else {
		ret = 0
	}
	if _x<0 {
		return -ret
	} else {
		return ret
	}
}

func dddctmrApply(ww []float64, xa float64) {
	ww[0] = -3
	ww[1] = 9
	ww[2] = -9
	ww[3] = 3
}



var zero1Kernel = rndKernel{
	Name: "Zero",
	Desc: "Returns zero everywhere (support 1)",
	Support: 1,
	Eval: ZeroEval,
	Apply: Zero1Apply,
}



var zero2Kernel = rndKernel{
	Name: "Zero",
	Desc: "Returns zero everywhere (support 2)",
	Support: 2,
	Eval: ZeroEval,
	Apply: Zero2Apply,
}

var zero3Kernel = rndKernel{
	Name: "Zero",
	Desc: "Returns zero everywhere (support 3)",
	Support: 3,
	Eval: ZeroEval,
	Apply: Zero3Apply,
}

var zero4Kernel = rndKernel{
	Name: "Zero",
	Desc: "Returns zero everywhere (support 4)",
	Support: 4,
	Eval: ZeroEval,
	Apply: Zero4Apply,
}

var zero5Kernel = rndKernel{
	Name: "Zero",
	Desc: "Returns zero everywhere (support 5)",
	Support: 5,
	Eval: ZeroEval,
	Apply: Zero5Apply,
}

var zero6Kernel = rndKernel{
	Name: "Zero",
	Desc: "Returns zero everywhere (support 6)",
	Support: 6,
	Eval: ZeroEval,
	Apply: Zero6Apply,
}

var dddCtmr = rndKernel{
	Name: "dddCtmr",
	Desc: "3rd deriv of Catmull-Rom",
	Support: 4,
	Eval: dddctmrEval,
	Apply: dddctmrApply,
	Deriv: &zero4Kernel,
}

var ddCtmr = rndKernel{
	Name: "dddCtmr",
	Desc: "2nd deriv of Catmull-Rom",
	Support: 4,
	Eval: ddctmrEval,
	Apply: ddctmrApply,
	Deriv: &dddCtmr,
}

var dCtmr = rndKernel{
	Name: "dddCtmr",
	Desc: "1st deriv of Catmull-Rom",
	Support: 4,
	Eval: dctmrEval,
	Apply: dctmrApply,
	Deriv: &ddCtmr,
}


var ctmr = rndKernel{
	Name: "Ctmr",
	Desc: "Catmull-Rom spline (C1, reconstructs quadratic)",
	Support: 4,
	Eval: ctmrEval,
	Apply: dctmrApply,
	Deriv: &dCtmr,
}

var rndKernelZero *rndKernel = &zero1Kernel
var rndKernelCtmr *rndKernel = &ctmr
// C way, we can use Map in Go
//var RndKernelAll = [](*rndKernel){
//	&zero1Kernel,
//	&ctmr,
//}
var rndKernelMap = map[string]*rndKernel{
	"ctmr": rndKernelCtmr,
	"zero": rndKernelZero,
}




func InitKernels() {
	zero1Kernel.Deriv = &zero1Kernel
	zero2Kernel.Deriv = &zero2Kernel
	zero3Kernel.Deriv = &zero3Kernel
	zero4Kernel.Deriv = &zero4Kernel
	zero5Kernel.Deriv = &zero5Kernel
	zero6Kernel.Deriv = &zero6Kernel
}

func rndKernelParse(_kstr string) *rndKernel{
	kstr := strings.ToLower(_kstr)
	if ret, ok := rndKernelMap[kstr]; ok {
		return ret
	}
	fmt.Println("The kernel [%s] does not exist", _kstr)
	return nil
}

