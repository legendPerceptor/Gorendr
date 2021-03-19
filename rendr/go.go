package rendr

import (
	"fmt"
	"sync"
)

type threadArg struct {
	ray *rndRay
	cnv *rndConvo
	ctx *rndCtx
	odata []float64
	plen uint
}

func Worker(wg*sync.WaitGroup, arg *threadArg) {
	//defer wg.Done()
	size0:= arg.ctx.Cam.Size[0]
	for {
		arg.ctx.rWlock.Lock()
		jj:= arg.ctx.nextTask
		arg.ctx.nextTask+=1
		arg.ctx.rWlock.Unlock()
		if jj>=arg.ctx.Cam.Size[1] {
			// all tasks are finished
			break
		}
		plen:=arg.plen
		index:=jj*plen*size0
		for ii:=uint(0); ii<size0; ii++ {
			odata:= arg.odata[index:index+plen]
			rndRayGo(odata, ii, jj, arg.ray, arg.cnv, arg.ctx)
			index += plen
		}
	}
	wg.Done()
}

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
		args := make([]threadArg,ctx.ThreadNum)
		var wg sync.WaitGroup

		for i:=uint(0); i<ctx.ThreadNum; i++ {
			args[i].plen=plen
			args[i].ctx=ctx
			args[i].cnv=rndConvoNew(ctx)
			args[i].odata = img.Data
			args[i].ray=rndNewRay()
			wg.Add(1)
			go Worker(&wg,&args[i])
		}
		wg.Wait()
	}
}