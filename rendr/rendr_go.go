package rendr

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type floatArrFlags [3]float64

func (f *floatArrFlags) String() string{
	return fmt.Sprintf("(%f, %f, %f)", (*f)[0], (*f)[1], (*f)[2])
}

func (f *floatArrFlags) Set(value string) error {
	s := strings.Fields(value)
	for i, _ := range s {
		(*f)[i], _ = strconv.ParseFloat(s[i], 64)
	}
	return nil
}

func goMain(goCmd *flag.FlagSet, args []string){
	var volFilename, metaVolFilename string
	goCmd.StringVar(&volFilename,"id", "cube.npy", "The input data file, cube.npy")
	goCmd.StringVar(&metaVolFilename,"im", "cube-meta.npy", "The ItoW data file for volume, cube-meta.npy")
	var fr, at, up floatArrFlags
	goCmd.Var(&fr, "fr", "Camera fr")
	goCmd.Var(&at, "at", "Camera at")
	goCmd.Var(&up, "up", "Camera at")
	var nc, fc, fov, us, s float64
	goCmd.Float64Var(&nc, "nc", 0, "Camera nc")
	goCmd.Float64Var(&fc, "fc", 0, "Camera fc")
	goCmd.Float64Var(&fov, "fov", 0, "Camera fov")
	goCmd.Float64Var(&us, "us", 0, "stepsize to use as the unit length")
	goCmd.Float64Var(&s, "s", 0, "spacing>0 between slice planes that contain ray samples.")
	var lutFilename, litFileName string
	goCmd.StringVar(&lutFilename, "lutd", "lut.npy", "The Lut data file")
	//goCmd.StringVar(&metaLutFilename, "lutm", "lut-meta.npy", "The Lut meta data file")
	goCmd.StringVar(&litFileName, "lit", "lit.txt", "The light, not file just string")
	var size floatArrFlags
	goCmd.Var(&size, "sz","Camera size")
	var outputDataFile string
	goCmd.StringVar(&outputDataFile, "od", "out.png", "output data")
	//goCmd.StringVar(&outputMetaDataFile, "om", "out-meta.npy", "output meta data")
	var ortho int
	goCmd.IntVar(&ortho, "ortho", 0, "Whether using orthogonal")
	var nt int
	var k, p, b string // kernel, probe, blend
	goCmd.IntVar(&nt, "nt", 0, "Multithreading: the number of threads")
	goCmd.StringVar(&k, "k", "ctmr", "kernel, default to be ctmr")
	goCmd.StringVar(&p, "p", "rgbalit", "probe, default to be rgbalit")
	goCmd.StringVar(&b, "b", "over", "blend mode, default to over")
	var dcn, dcf floatArrFlags
	goCmd.Var(&dcn, "dcn", "depth cueing dcn")
	goCmd.Var(&dcf, "dcf", "depth cueing dcf")
	if err:=goCmd.Parse(args); err!= nil {
		fmt.Printf("Failed to parse commandline arguments in rndGo")
		os.Exit(1)
	}
	fmt.Println("Finished parsing arguments")
	// Start the real stuff
	var vol *rndVolume
	var ctx *rndCtx
	vol = rndVolumeNew()
	vol.loadVolume(volFilename, metaVolFilename)
	var lit [7]float64
	litStr :=strings.Fields(litFileName)
	for i, value:= range litStr{
		lit[i], _ = strconv.ParseFloat(value, 64)
	}
	ctx = rndCtxNew(vol, rndKernelParse(k), fr, at, up, dcn,dcf,size, nc, fc, fov, us, s,p, b,ortho, nt, lit, lutFilename)
	fmt.Println("Start rendering!")
	start := time.Now()
	img := rndImageNew()
	rndRender(img, ctx, -1, -1)
	duration := time.Since(start)
	fmt.Printf("The duration of rendering is %d\n", duration.Milliseconds())
	img.Content = "Gorendr beta (partial implementation)"
	rndImageSave(outputDataFile, img)
}