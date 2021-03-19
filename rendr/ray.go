package rendr

import (
	"gonum.org/v1/gonum/mat"
	"math"
)

func rndQuantize(min, val, max float64, num uint) uint {
	ret:=uint(0)
	if val<=min {
		ret = uint(0)
	} else if val>=max {
		ret = num-1
	} else {
		tmp:= float64(num)*(val-min)/(max-min)
		ret = uint(tmp)
		if ret==num{
			ret-=1
		}
	}
	return ret
}

func rndLerp(omin, omax, imin, xx, imax float64) float64 {
	alpha:= (xx-imin)/(imax-imin)
	return (1-alpha)*omin + alpha*omax
}

func rndRayStart(hi, vi uint, ray *rndRay, cnv *rndConvo, ctx *rndCtx)  {
	ray.hi = hi
	ray.vi = vi
	var tmp float64
	if ctx.Cam.Ortho>0 {
		ray.r_img = [3]float64{0, 0, -ctx.Cam.D}
		ray.r0[0] = (ctx.Cam.Wdth/2)*rndLerp(-1, 1, -0.5, float64(hi), float64(ctx.Cam.Size[0])-0.5)
		ray.r0[1] = (ctx.Cam.Hght/2)*rndLerp(1, -1, -0.5, float64(vi), float64(ctx.Cam.Size[1])-0.5)
	} else {
		ray.r_img = [3]float64{(ctx.Cam.Wdth/2)*rndLerp(-1, 1, -0.5, float64(hi), float64(ctx.Cam.Size[0])-0.5),
			(ctx.Cam.Hght/2)*rndLerp(1, -1, -0.5, float64(vi), float64(ctx.Cam.Size[1])-0.5), -ctx.Cam.D}
		tmp = ctx.Cam.Ncv/ctx.Cam.D
		ray.r0 = [3]float64{tmp*ray.r_img[0], tmp*ray.r_img[1], tmp*ray.r_img[2]}
	}
	ray.k = -1
	ray.p[3]=1 //homogeneous
	tmp = ctx.PlaneSep/ctx.Cam.D
	ray.r_step = [3]float64{tmp*ray.r_img[0], tmp*ray.r_img[1], tmp*ray.r_img[2]}
	ray.delta = math.Sqrt(ray.r_step[0]*ray.r_step[0]+ray.r_step[1]*ray.r_step[1]+ray.r_step[2]*ray.r_step[2])
	ray.T = 1
	ray.rgb = [3]float64{ctx.OutsideValue, ctx.OutsideValue, ctx.OutsideValue}
	ray.result = [4]float64{ctx.OutsideValue,0,0,0}
	if ctx.Blend == rndBlendSum {
		ray.mid_result = [4]float64{0,0,0,0}
		ray.litresult = [4]float64{0,0,0,0}
	} else if ctx.Blend == rndBlendMax {
		ray.mid_result = [4]float64{-32768,-32768,-32768,-32768}
		ray.litresult = [4]float64{-32768,-32768,-32768,-32768}
	}
	dist:= -1/math.Sqrt(ray.r_img[0]*ray.r_img[0]+ray.r_img[1]*ray.r_img[1]+ray.r_img[2]*ray.r_img[2])
	ray.VdirView = [4]float64{ray.r_img[0]*dist, ray.r_img[1]*dist, ray.r_img[2]*dist, 0}
	vtow:=mat.NewDense(4,4,ctx.Cam.VtoW[:])
	vdir:=mat.NewVecDense(4, ray.Vdir[:])
	vdirView := mat.NewVecDense(4, ray.VdirView[:])
	vdir.MulVec(vtow, vdirView)
	dist = 1/math.Sqrt(ray.Vdir[0]*ray.Vdir[0]+ray.Vdir[1]*ray.Vdir[1]+ray.Vdir[2]*ray.Vdir[2])
	ray.Vdir = [4]float64{ray.Vdir[0]*dist, ray.Vdir[1]*dist, ray.Vdir[2]*dist, 0}
	ray.set = 0
	if ctx.SampleOnceStop>0 {
		rndRayStep(ray, cnv, ctx)
	}
}

func rndRayStep(ray *rndRay, cnv *rndConvo, ctx *rndCtx) bool{
	var world_p, index_p [4]float64
	var tmp,color [3]float64
	ray.k = ray.k + 1
	tmp = [3]float64{ray.r_step[0]*float64(ray.k),ray.r_step[1]*float64(ray.k),ray.r_step[2]*float64(ray.k)}
	ray.p = [4]float64{ray.r0[0]+tmp[0], ray.r0[1]+tmp[1], ray.r0[2]+tmp[2], 1}
	if -ray.p[2] - ctx.Cam.Fcv > 0 {
		return false
	}
	vtow:=mat.NewDense(4,4,ctx.Cam.VtoW[:])
	wtoi:=mat.NewDense(4,4,ctx.WtoI[:])
	rayp:= mat.NewVecDense(4, ray.p[:])
	worldp := mat.NewVecDense(4, world_p[:])
	indexp := mat.NewVecDense(4, index_p[:])
	worldp.MulVec(vtow, rayp)
	indexp.MulVec(wtoi, worldp)
	rndConvoEval(world_p[0], world_p[1], world_p[2], cnv, ctx)
	var m_alpha float64
	// TODO: Other than RGBA
	if cnv.Inside==1 && ctx.Txf.len>0 {
		t:= rndQuantize(ctx.Txf.vmin, cnv.Value, ctx.Txf.vmax, ctx.Txf.len)
		rgb:= ctx.Txf.rgba[4*t:4*t+3]
		m_alpha = ctx.Txf.rgba[4*t+3]
		if m_alpha < 0 {
			m_alpha = 0
		} else if m_alpha>1 {
			m_alpha = 1
		}
		m_alpha = 1 - math.Pow(1-m_alpha, ray.delta/ctx.Txf.unitStep)
		color = [3]float64{rgb[0], rgb[1], rgb[2]}
		if ctx.Probe==rndProbeRgbaLit {
			blin:=rndBlinnPhong([3]float64{rgb[0], rgb[1], rgb[2]}, cnv.Gradient, ray.Vdir, ctx)
			copy(color[:], blin[:])
			// TODO: depth cueing
		}
		ray.set = 1;
		if ctx.Blend==rndBlendOver {
			ray.rgb = [3]float64{ray.rgb[0]+ray.T*m_alpha*color[0],ray.rgb[1]+ray.T*m_alpha*color[1],ray.rgb[2]+ray.T*m_alpha*color[2]}
			ray.T *= 1-m_alpha
			if 1-ray.T - ctx.Txf.alphaNear1 >0 {
				ray.T = 0
				return false
			}
		}
		//TODO: BlendSum and BlendMax
	}
	if ctx.SampleOnceStop==1 {
		return false
	}
	return true
}

func rndRayFinish(ray *rndRay, ctx *rndCtx) {
	var alpha float64
	switch ctx.Probe {
	case rndProbeRgba:
		fallthrough
	case rndProbeRgbaLit:
		if ray.set==0 {
			ray.result = [4]float64{ctx.OutsideValue, ctx.OutsideValue, ctx.OutsideValue, ctx.OutsideValue}
		}
		if ctx.Blend == rndBlendOver {
			alpha = 1-ray.T
			if alpha<0 {
				alpha=0
			}else if alpha>1 {
				alpha=1
			}
			if alpha>0 {
				ray.result= [4]float64{ray.rgb[0], ray.rgb[1], ray.rgb[2], alpha}
			}else{
				ray.result = [4]float64{ctx.OutsideValue, ctx.OutsideValue, ctx.OutsideValue, ctx.OutsideValue}
			}
		}
	default:
		if ray.mid_result[0]==-32768{
			ray.result = [4]float64{ctx.OutsideValue, ctx.OutsideValue, ctx.OutsideValue, ctx.OutsideValue}
		}
		copy(ray.result[:], ray.mid_result[:])
	}
}

func rndRayGo(out []float64, ii, jj uint, ray *rndRay, cnv *rndConvo, ctx *rndCtx) {
	rndRayStart(ii, jj, ray, cnv, ctx)
	if ctx.SampleOnceStop == 0 {
		var keepGoing bool
		for {
			keepGoing = rndRayStep(ray, cnv, ctx)
			if !keepGoing {
				break
			}
		}
	}
	rndRayFinish(ray, ctx)
	plen := rndProbeLen(ctx.Probe)
	for pi:=uint(0); pi<plen; pi++ {
		out[pi] = ray.result[pi]
	}
	//if ctx.Timing>0 {
	//	out[plen] = ray.millisecs
	//}
}

func rndBlinnPhong(rgbIn [3]float64, grad [3]float64, Vdir [4]float64, ctx *rndCtx) [3]float64{
	var ambient, diffuse, specular, N, H, result [3]float64
	ambient = [3]float64{ctx.Lparam.ka * rgbIn[0], ctx.Lparam.ka * rgbIn[1], ctx.Lparam.ka * rgbIn[2]}
	copy(result[:], ambient[:])
	dist:= math.Sqrt(grad[0]*grad[0]+grad[1]*grad[1]+grad[2]*grad[2])
	if ctx.Light.num==0 || (ctx.Lparam.kd==0 && ctx.Lparam.ks==0) || dist==0  {
		return result
	}
	N = [3]float64{-grad[0]/dist, -grad[1]/dist, -grad[2]/dist}
	for i:=uint(0); i<ctx.Light.num; i++ {
		lcol := ctx.Light.rgb[3*i:3*i+3]
		ldir := ctx.Light.xyz[3*i:3*i+3]
		dotTmp:= N[0]*ldir[0] + N[1]*ldir[1] + N[2]*ldir[2]
		tmp:= ctx.Lparam.kd * math.Max(0, dotTmp)
		diffuse = [3]float64{tmp*rgbIn[0]*lcol[0],
			tmp*rgbIn[1]*lcol[1], tmp*rgbIn[2]*lcol[2]}
		H = [3]float64{Vdir[0]+ldir[0], Vdir[1]+ldir[1], Vdir[2]+ldir[2]}
		dist = 1/math.Sqrt(H[0]*H[0]+H[1]*H[1]+H[2]*H[2])
		H = [3]float64{H[0]*dist, H[1]*dist, H[2]*dist}
		dotTmp = N[0]*H[0]+N[1]*H[1]+N[2]*H[2]
		tmp = ctx.Lparam.ks * math.Pow(math.Max(0, dotTmp), ctx.Lparam.p)
		specular = [3]float64{tmp*lcol[0], tmp*lcol[1], tmp*lcol[2]}
		for j:=0;j<3;j++ {
			result[j] = ambient[j] + specular[j] + diffuse[j]
		}
	}
	return result
}