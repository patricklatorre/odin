package main

import (
	"flag"
	"fmt"
	"os"
)

const Version = "0.2.0"

// Path of the odin executable. Other paths will be relative to this.
var OdinExePath string

func main() {
	mustSetup()

	var (
		flagHelp      = flag.Bool("h", false, "You're looking at it")
		flagVersion   = flag.Bool("v", false, "Prints the version")
		createCmd     = flag.NewFlagSet("create", flag.ExitOnError)
		openCmd       = flag.NewFlagSet("open", flag.ExitOnError)
		startCmd      = flag.NewFlagSet("start", flag.ExitOnError)
		startPort     = startCmd.Int("port", 2456, "")
		startPassword = startCmd.String("password", "123456", "")
	)

	// Change func for printing usage
	flag.Usage = func() {
		fmt.Printf("Usage: odin <COMMAND> <WORLDNAME> [..OPTIONS]\n\n" +
			"<COMMAND>\n" +
			" help                        You're looking at it\n" +
			" open   <world>              Opens the server in explorer\n" +
			" create <world>              Creates a new server\n" +
			" start  <world>              Starts a server\n" +
			"        [-port 2456]\n" +
			"        [-password 123456]\n" +
			"\n")
	}

	// Parse args
	flag.Parse()
	args := flag.Args()

	// Handle version and help flags
	if *flagVersion {
		PrintVersion()
		os.Exit(0)
	} else if len(args) == 0 || *flagHelp {
		PrintHelp()
		os.Exit(0)
	}

	// Validate command
	var (
		validCmds  = []string{"create", "start", "open", "help"}
		isValidCmd = false
		cmd        = args[0]
	)

	for _, validCmd := range validCmds {
		if cmd == validCmd {
			isValidCmd = true
		}
	}

	if !isValidCmd {
		fmt.Println("Invalid command:", cmd)
		os.Exit(1)
	}

	// Handle no-arg subcommands
	switch cmd {
	case "help":
		PrintHelp()
		os.Exit(0)
	}

	// Other subcommands expect at least 2 non-flag args
	if len(args) < 2 {
		fmt.Println("You must provide a server name for this command")
		os.Exit(1)
	}

	// Handle subcommands
	switch cmd {
	case "create":
		createCmd.Parse(args[1:])
		name := createCmd.Arg(0)
		err := Create(name)
		Must(err)

	case "start":
		startCmd.Parse(args[1:])
		name := startCmd.Arg(0)
		err := Start(name, *startPort, *startPassword)
		Must(err)

	case "open":
		openCmd.Parse(args[1:])
		name := openCmd.Arg(0)
		Open(name)

	default:
		fmt.Println("Invalid command:", cmd)
		os.Exit(1)
	}
}

// Panics if error is not null
func Must(err error) {
	if err != nil {
		panic(err)
	}
}
