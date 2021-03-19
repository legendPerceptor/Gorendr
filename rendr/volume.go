package rendr

import (
	"fmt"
	"github.com/sbinet/npyio"
	"gonum.org/v1/gonum/mat"
	"math"
	"os"
)

func rndVolumeNew() *rndVolume {
	ret := rndVolume{
		Content: "",
		Size:    [3]uint{0},
		ItoW:    [16]float64{0},
		Dtype:   rndTypeUnknown,
		Data:    nil,
	}
	return &ret
}

func (v *rndVolume) loadVolume(datafile string, metadatafile string) {
	f, _ := os.Open(datafile)
	r, err := npyio.NewReader(f)
	if err!=nil || r==nil{
		fmt.Printf("Failed to open file %s\n", datafile)
		os.Exit(1)
	}
	fmt.Println(r.Header)
	shape := r.Header.Descr.Shape
	v.Size = [3]uint{uint(shape[0]),uint(shape[1]),uint(shape[2])}
	raw := make([]float64, shape[0]*shape[1]*shape[2])
	err = r.Read(&raw)
	if err!=nil {
		fmt.Printf(err.Error())
	}
	v.Data = raw
	v.Dtype = rndTypefloat64
	f, _ = os.Open(metadatafile)
	r2, err2 := npyio.NewReader(f)
	if err2!=nil || r2==nil {
		fmt.Printf("failed to open file %s\n", metadatafile)
		os.Exit(1)
	}
	shape2 := r2.Header.Descr.Shape
	fmt.Println(shape2)
	data2:= make([]float64, shape2[0]*shape2[1])
	err2 = r2.Read(&data2)
	if err2!=nil {
		fmt.Printf(err.Error())
	}
	copy(v.ItoW[:], data2)
}

func (txf *_txf) loadLut(datafile string) {
	f, _ := os.Open(datafile)
	r, err := npyio.NewReader(f)
	if err!=nil || r==nil{
		fmt.Printf("Failed to open file %s\n", datafile)
		os.Exit(1)
	}
	fmt.Println(r.Header)
	shape := r.Header.Descr.Shape
	txf.len = uint(shape[1])
	raw := make([]float64, shape[0]*shape[1])
	err = r.Read(&raw)
	if err!=nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}
	txf.rgba = raw
}

func (ctx *rndCtx) rndCtxLightUpdate() {
	var vdir = [4]float64{
		ctx.Light.dir[0], ctx.Light.dir[1], ctx.Light.dir[2], 0,
	}
	vtowMat := mat.NewDense(4,4,ctx.Cam.VtoW[:])
	vdirVec := mat.NewVecDense(4, vdir[:])
	var tmp [4]float64
	result := mat.NewVecDense(4, tmp[:])
	result.MulVec(vtowMat, vdirVec)
	dist:= math.Sqrt(tmp[0]*tmp[0] + tmp[1]*tmp[1] + tmp[2]*tmp[2])
	ctx.Light.xyz = [3]float64{
		tmp[0]/dist, tmp[1]/dist,tmp[2]/dist,
	}
}

func rndImageNew() *rndImage {
	return &rndImage{
		Content: "",
		Channel: 0,
		Size:    [2]uint{},
		Dtype:   0,
		Data:    nil,
	}
}


func rndNewRay() *rndRay {
	return &rndRay{
		hi:           0,
		vi:           0,
		result:       [4]float64{},
		sms0:         [2]int64{},
		sms1:         [2]int64{},
		millisecs:    0,
		r_img:        [3]float64{},
		r0:          [3]float64{},
		r_step:       [3]float64{},
		rgb:          [3]float64{},
		rgb_material: [3]float64{},
		k:            0,
		set:          0,
		T:            0,
		delta:        0,
		p:            [4]float64{},
		mid_result:   [4]float64{},
		litresult:    [4]float64{},
		VdirView:     [4]float64{},
		Vdir:         [4]float64{},
	}
}


