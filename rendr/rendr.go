package rendr

import (
	"flag"
	"fmt"
	"os"
)

func RendrMain() {
	cevalCmd := flag.NewFlagSet("ceval", flag.ExitOnError)
	goCmd := flag.NewFlagSet("go", flag.ExitOnError)
	if len(os.Args) < 2 {
		// TODO: Helper Information
		fmt.Println("Expected subcommand 'ceval' or 'go' ")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "ceval":
		cevalMain(cevalCmd)
	case "go":
		goMain(goCmd)
	default:
		fmt.Println("The subcommand %s is not supported", os.Args[1])
		os.Exit(1)
	}
}