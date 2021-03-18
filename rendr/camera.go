package rendr

import (
	"math"
)

func rndCameraNew() *rndCamera  {
	return &rndCamera{
		Fr:    [3]float64{},
		At:    [3]float64{},
		Up:    [3]float64{},
		Nc:    0,
		Fc:    0,
		FOV:   0,
		Size:  [2]uint{},
		Ortho: 0,
		Ar:    0,
		D:     0,
		Ncv:   0,
		Fcv:   0,
		Hght:  0,
		Wdth:  0,
		U:     [3]float64{},
		V:     [3]float64{},
		N:     [3]float64{},
		VtoW:  [16]float64{},
	}
}

func cross(a,b [3]float64) [3]float64 {
	return [3]float64{a[1]*b[2] - a[2]*b[1], a[2]*b[0] - a[0]*b[2], a[0]*b[1] - a[1]*b[0]}
}

func dot(a, b [3]float64) float64 {
	return a[0]*b[0] + a[1]*b[1] + a[2]*b[2]
}

func (cam *rndCamera) rndCameraUpdate() {
	dist:=float64(0)
	var d_v, u_v [3]float64
	var i int
	for i=0; i<3;i++ {
		d_v[i] = cam.Fr[i] - cam.At[i]
		dist += d_v[i] * d_v[i]
	}
	dist = math.Sqrt(dist)
	cam.D = dist
	for i=0; i<3;i++ {
		cam.N[i] = d_v[i]/dist
	}
	u_v = cross(cam.Up, cam.N)
	dist = 0
	for i=0; i<3;i++ {
		dist += u_v[i]*u_v[i]
	}
	dist = math.Sqrt(dist)
	for i=0;i<3;i++ {
		cam.U[i] = u_v[i]/dist
	}
	cam.V = cross(cam.N, cam.U)
	M_view:= [16]float64{
		cam.U[0], cam.U[1], cam.U[2], -dot(cam.Fr, cam.U),
		cam.V[0], cam.V[1], cam.V[2], -dot(cam.Fr, cam.V),
		cam.N[0], cam.N[1], cam.N[2], -dot(cam.Fr, cam.N),
		0,0,0,1,
	}
	cam.VtoW, _ = rnd_m4_affine_inv(M_view)
	cam.Ncv = cam.Nc + cam.D
	cam.Fcv = cam.Fc + cam.D
	cam.Hght = 2* cam.D * math.Tan(cam.FOV/2 * math.Pi/180)
}

func (cam *rndCamera) rndCameraSet(fr,at,up [3]float64, nc,fc,FOV float64, size0, size1 uint, ortho int){
	cam.Fr = fr
	cam.At = at
	cam.Up = up
	cam.Nc = nc
	cam.Fc = fc
	cam.FOV = FOV
	cam.Size[0] = size0
	cam.Size[1] = size1
	cam.Ortho = ortho
}