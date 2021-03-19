package rendr

var blendmap = map[string] rndBlend {
	"over": rndBlendOver,
	"max": rndBlendMax,
	"sum": rndBlendSum,
}
var probemap = map[string] rndProbe {
	"rgbalit":rndProbeRgbaLit,
	"rgba":rndProbeRgba,
	"value":rndProbeValue,
}

func rndCtxNew(vol *rndVolume, kern *rndKernel,
	fr, at, up, dcn,dcf,size [3]float64,
	nc, fc, fov, us, s float64,
	p, b string,
	orth, nt int,
	lit [7]float64,
	lutfile string) *rndCtx{
	ctx := rndCtx{
		Vol:            vol,
		Kern:           kern,
		SampleOnceStop: 0,
		PlaneSep:       s,
		Probe:          probemap[p],
		Blend:          blendmap[b],
		ThreadNum:      uint(nt),
		OutsideValue:   0,
		Timing:         0,
		Cam:            rndCameraNew(),
		Txf:            _txf{
			len:        0,
			vmin:       0,
			vmax:       1,
			rgba:       nil,
			unitStep:   us,
			alphaNear1: 0.98,
		},
		Levoy:          _levoy{
			num: 0,
			vra: nil,
		},
		Light:          _light{
			num: 1,
			rgb: [3]float64{lit[0], lit[1], lit[2]},
			dir: [3]float64{lit[3], lit[4], lit[5]},
			vsp: int(lit[6]),
			xyz: [3]float64{},
		},
		Lparam:         _lparam{
			ka:  0.2,
			kd:  0.8,
			ks:  0.1,
			p:   150,
			dcn: dcn,
			dcf: dcf,
		},
		WtoI:           [16]float64{},
		MT:             [9]float64{},
	}
	ctx.Cam.rndCameraSet(fr, at, up, nc, fc, fov, uint(size[0]), uint(size[1]), orth)
	ctx.Cam.rndCameraUpdate()
	ctx.Txf.loadLut(lutfile)
	ctx.rndCtxLightUpdate()
	ctx.WtoI, ctx.MT = rnd_m4_affine_inv(ctx.Vol.ItoW)
	return &ctx
}

func rndProbeLen(probe rndProbe) uint {
	l:=uint(0)
	switch  probe {
	case rndProbeRgba:
		l=4
	case rndProbeRgbaLit:
		l=4
	default:
		l=1
	}
	return l
}