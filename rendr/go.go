package rendr

import "fmt"

func rndRender(img *rndImage, ctx *rndCtx, _dhi, _dvi int) {
	plen:= rndProbeLen(ctx.Probe)
	//timing := uint(ctx.Timing)
	img.rndImageAlloc(plen, ctx.Cam.Size[0], ctx.Cam.Size[1], rndTypefloat64)
	cnv := rndConvoNew(ctx)
	if ctx.ThreadNum==0 {
		ray := rndNewRay()
		index:= uint(0)
		odata:=img.Data
		//nipx := ctx.Cam.Size[0] * ctx.Cam.Size[1]
		//fmt.Fprintf(os.Stderr, "Rendering")
		for jj:= uint(0); jj<ctx.Cam.Size[1]; jj++ {
			for ii:=uint(0); ii<ctx.Cam.Size[0]; ii++ {
				rndRayGo(odata[index:index+plen], ii, jj, ray, cnv, ctx)
				index = index+plen
			}
		}
	} else {
		fmt.Println("Multithreading version of Gorendr")
	}
}