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
		fmt.Println("It is necessary at least a subcommand to continue. Type 'cleaner --help' so you can see the available subcommands and flags.")
		os.Exit(1)
	}

	isValidOp := slices.Contains(cliOps, cmd)
	if !isValidOp {
		usage()
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

func usage() {
	fmt.Print(`Cleanup is a tool for cloud providers' resources cleaning. You can quickly list, validate or delete resources from your current cloud provider vendor. The tool is being incrementally built, so the only provider currently supported is AWS with few resources.

Usage:
	cleanup <command> (--service | -service) <service_name>

Commands Usage:
	cleanup list (--service | -service) <service_name>

		Options:
			--service STRING  (required) the service name you're trying to list (these service names are available in the documentation).

	cleanup validate (--service | -service) <service_name>

		Options:
			--service STRING  (required) the service name you're trying to validate (these service names are available in the documentation). Each resource has its own rules to be considered empty, so check docs to understand these rules.

	cleanup delete (--service | -service) <service_name>

		Options:
			--service STRING  (required) the service name you're trying to delete (these service names are available in the documentation). A resource can only be deleted if empty, so check it first using the "validate" operation.
	`)
}

// Considerations about this package:
// -> I don't think this is well written nor something extensible. The Start() function has some limitations that could cause some non-expected errors. E.g.: It does not have the flag <-> flag value relationship, which would be something hard to deal if we had more than one flag;
// -> Since this is starting simple and meant to be simple, I won't rewrite this package to prevent non-existent (at least for now) scenarios;
// -> The first value inside the slice returned by the Start() function should always be the operation being done (list, validate or delete);
