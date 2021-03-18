package rendr

import (
	"flag"
	"fmt"
	"os"
)

func cevalMain(cevalCmd *flag.FlagSet,args []string) {
	fname := cevalCmd.String("input", "", "filename of input volume")
	cevalCmd.StringVar(fname, "i", "", "alias to input")
	wpos := cevalCmd.Float64("position", 0,"the single position (in world space) at which to evaluate the convolution between imageand kernel")
	cevalCmd.Float64Var(wpos, "w", 0,"alias to position")
	wtoi := cevalCmd.Int("bool", 0,"instead of printing convolution results, just print the index-space position ctx->ipos; computing this is the first step of rndConvoEval")
	cevalCmd.IntVar(wtoi, "wtoi", 0, "alias to bool")
	kernStr := cevalCmd.String("kernel", "ctmr", "kernel to use for convolution")
	cevalCmd.StringVar(kernStr, "k","ctmr", "alias to kernel")
	findgrad := cevalCmd.Int("findgrad", 0, "find not just value, but gradient (in world-space)")
	cevalCmd.IntVar(findgrad, "g", 0, "alias to findgrad")
	if err:=cevalCmd.Parse(args); err!= nil {
		fmt.Printf("Failed to parse commandline arguments in ceval")
		os.Exit(1)
	}
	//kernel to use
	//kern := rndKernelParse(*kernStr)
	//var ctx *rndCtx

}