package rendr

import (
	"flag"
)

func cevalMain(cevalCmd *flag.FlagSet) {
	fname := cevalCmd.String("input", "", "filename of input volume")
	cevalCmd.StringVar(fname, "i", "", "alias to input")
	wpos := cevalCmd.Float64("position", 0,"the single position (in world space) at which to evaluate the convolution between imageand kernel")
	cevalCmd.Float64Var(wpos, "w", 0,"alias to position")

}