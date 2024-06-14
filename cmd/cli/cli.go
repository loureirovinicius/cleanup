package cli

import (
	"flag"
	"fmt"
	"os"
	"slices"
)

// Starts the CLI and return the argument and the subcommands' values
func Start() []string {
	var args []string
	var cmd = os.Args[1]
	var service string
	var cliOps = []string{"list", "validate", "delete"}

	if len(os.Args) < 2 {
		fmt.Println("It is necessary at least a subcommand to continue. Type 'passgen --help' so you can see the available subcommands and flags.")
		os.Exit(1)
	}

	isValidOp := slices.Contains(cliOps, cmd)
	if !isValidOp {
		flag.Usage = func() {
			usage()
		}
	}

	// These blocks are meant to be incremental, so according to new options are created for each operation, they can be incremented
	// List operation flag set
	list := flag.NewFlagSet("list", flag.ExitOnError)
	list.StringVar(&service, "service", "", "cloud provider service")
	list.Parse(os.Args[2:])

	// Validate operation flag set
	validate := flag.NewFlagSet("validate", flag.ExitOnError)
	validate.StringVar(&service, "service", "", "cloud provider service")
	validate.Parse(os.Args[2:])

	// Delete operation flag set
	delete := flag.NewFlagSet("delete", flag.ExitOnError)
	delete.StringVar(&service, "service", "", "cloud provider service")
	delete.Parse(os.Args[2:])

	args = []string{cmd, service}

	return args
}

func usage() {}

// Considerations about this package:
// -> I don't think this is well written nor something extensible. The Start() function has some limitations that could cause some non-expected errors. E.g.: It does not have the flag <-> flag value relationship, which would be something hard to deal if we had more than one flag;
// -> Since this is starting simple and meant to be simple, I won't rewrite this package to prevent non-existent (at least for now) scenarios;
// -> The first value inside the slice returned by the Start() function should always be the operation being done (list, validate or delete);
