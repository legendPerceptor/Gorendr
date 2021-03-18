package rendr

import (
	"flag"
	"fmt"
	"os"
)

func RendrMain() {
	cevalCmd := flag.NewFlagSet("ceval", flag.ExitOnError)
	goCmd := flag.NewFlagSet("go", flag.ExitOnError)
	m4aiCmd :=flag.NewFlagSet("m4ai", flag.ExitOnError)
	if len(os.Args) < 2 {
		// TODO: Helper Information
		fmt.Println("Expected subcommand 'ceval' or 'go' ")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "m4ai":
		m4aiMain(m4aiCmd, os.Args[2:])
	case "ceval":
		cevalMain(cevalCmd, os.Args[2:])
	case "go":
		goMain(goCmd, os.Args[2:])
	default:
		fmt.Printf("The subcommand %s is not supported\n", os.Args[1])
		os.Exit(1)
	}
}