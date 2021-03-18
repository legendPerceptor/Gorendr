package rendr

import (
	"fmt"
	"testing"
)

func TestLoadVolume(t *testing.T) {
	v := rndVolumeNew()
	v.loadVolume("cube.npy", "cube-meta.npy")
	s1:="("
	for _, value:= range v.Size {
		//var t string
		t:=fmt.Sprintf( "%d  ", value)
		s1+=t
	}
	s1+=")"
	t.Logf(s1)
	s:= fmt.Sprintf("ItoW=[%.5f, %.5f, %.5f, %.5f\n    %.5f, %.5f, %.5f, %.5f\n    %.5f, %.5f, %.5f, %.5f\n    %.5f, %.5f, %.5f, %.5f]\n",
		v.ItoW[0], v.ItoW[1], v.ItoW[2], v.ItoW[3],
		v.ItoW[4], v.ItoW[5], v.ItoW[6], v.ItoW[7],
		v.ItoW[8], v.ItoW[9], v.ItoW[10], v.ItoW[11],
		v.ItoW[12], v.ItoW[13], v.ItoW[14], v.ItoW[15])
	t.Logf(s)
}