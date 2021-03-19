package rendr

import (
	"fmt"
	"github.com/sbinet/npyio"
	"os"
	"testing"
)

//func TestLoadVolume(t *testing.T) {
//	v := rndVolumeNew()
//	v.loadVolume("cube.npy", "cube-meta.npy")
//	s1:="("
//	for _, value:= range v.Size {
//		//var t string
//		t:=fmt.Sprintf( "%d  ", value)
//		s1+=t
//	}
//	s1+=")"
//	t.Logf(s1)
//	s:= fmt.Sprintf("ItoW=[%.5f, %.5f, %.5f, %.5f\n    %.5f, %.5f, %.5f, %.5f\n    %.5f, %.5f, %.5f, %.5f\n    %.5f, %.5f, %.5f, %.5f]\n",
//		v.ItoW[0], v.ItoW[1], v.ItoW[2], v.ItoW[3],
//		v.ItoW[4], v.ItoW[5], v.ItoW[6], v.ItoW[7],
//		v.ItoW[8], v.ItoW[9], v.ItoW[10], v.ItoW[11],
//		v.ItoW[12], v.ItoW[13], v.ItoW[14], v.ItoW[15])
//	t.Logf(s)
//}

func loadCube(datafile string, t *testing.T) {
	f, _ := os.Open(datafile)
	r, err := npyio.NewReader(f)
	if err!=nil || r==nil{
		fmt.Printf("Failed to open file %s\n", datafile)
		os.Exit(1)
	}
	fmt.Println(r.Header)
	shape := r.Header.Descr.Shape
	size := [3]uint{uint(shape[0]),uint(shape[1]),uint(shape[2])}
	raw := make([]float64, shape[0]*shape[1]*shape[2])
	err = r.Read(&raw)
	if err!=nil {
		fmt.Printf(err.Error())
	}
	t.Logf("(%d %d %d)", size[0], size[1], size[2])

}

func TestCube(t *testing.T) {
	loadCube("test-cube.npy", t)
}