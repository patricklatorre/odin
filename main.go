package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	e "github.com/patricklatorre/odin/error"
	"github.com/patricklatorre/odin/path"
)

const Version = "0.1.0"

func main() {
	// Setup Odin directories if it doesn't exist
	mustSetupDirs()

	// Odin flags
	flagHelp := flag.Bool("h", false, "You're looking at it")
	flagVersion := flag.Bool("v", false, "Prints the version")

	// Create flags
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)

	// Start flags
	startCmd := flag.NewFlagSet("start", flag.ExitOnError)
	startPort := startCmd.Int("port", 55555, "")
	startPassword := startCmd.String("password", "1234", "")

	// Open flags
	openCmd := flag.NewFlagSet("open", flag.ExitOnError)

	// Help text helper func
	printUsage := func() {
		fmt.Printf("ODIN %s\n\n", Version)
		fmt.Printf("Usage: odin <command> <server_name> [...params]\n\n")

		fmt.Printf("-------- Commands --------\n\n")

		fmt.Printf("CREATE <SERVER_NAME>\n")
		createCmd.PrintDefaults()

		fmt.Printf("\nSTART <SERVER_NAME>\n")
		startCmd.PrintDefaults()

		fmt.Printf("\nOPEN <SERVER_NAME>\n")
		openCmd.PrintDefaults()

		fmt.Printf("\nHELP\n")
	}

	// Parse odin flags
	flag.Usage = printUsage
	flag.Parse()
	args := flag.Args()

	// Print help as default behavior
	if len(args) == 0 || *flagVersion || *flagHelp {
		help()
	}

	// Validate command
	cmd := args[0]
	isValidCmd := false
	validCmds := []string{"create", "start", "open", "help"}
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
		help()
	}

	// Subcommands expect at least 2 non-flag args
	if len(args) < 2 {
		fmt.Println("You must provide a server name for this command")
		os.Exit(1)
	}

	// Handle subcommands
	switch cmd {
	case "create":
		createCmd.Parse(args[1:])
		name := createCmd.Arg(0)
		err := create(name)
		e.Must(err)

	case "start":
		startCmd.Parse(args[1:])
		name := startCmd.Arg(0)
		err := start(name, *startPort, *startPassword)
		e.Must(err)

	case "open":
		openCmd.Parse(args[1:])
		name := openCmd.Arg(0)
		open(name)

	default:
		fmt.Println("Invalid command:", cmd)
		os.Exit(1)
	}
}

// Creates the Odin dirs if it doesn't exist
func mustSetupDirs() {
	// Required dirs
	dirs := []string{
		path.Relative("servers"),
		path.Relative("worlds"),
		path.Relative("steamcmd"),
	}

	// Create required dirs
	for _, dir := range dirs {
		exists, err := path.Exists(dir)
		e.Must(err)
		if !exists {
			err := os.Mkdir(dir, 0755)
			e.Must(err)
			fmt.Printf("Created directory: %s\n", dir)
		}
	}

	{ // Check if steamcmd.exe exists
		exe := path.Relative("steamcmd", "steamcmd.exe")
		exists, err := path.Exists(exe)
		e.Must(err)
		if !exists {
			fmt.Printf("SteamCMD not found, downloading... ")
			err := setupSteamcmd()
			e.Must(err)
			fmt.Printf("done\n")
		}
	}
}

// Downloads and unpacks steamcmd
func setupSteamcmd() error {
	var (
		dlUrl  = "https://steamcdn-a.akamaihd.net/client/installer/steamcmd.zip"
		dlPath = path.Relative("steamcmd.zip")
	)

	// Download the zip file
	response, err := http.Get(dlUrl)
	if err != nil {
		fmt.Printf("Could not download the steamcmd.zip")
		return err
	}

	// Create a new file for writing the zip file
	zipFile, err := os.Create(dlPath)
	if err != nil {
		fmt.Printf("Could not download the steamcmd.zip")
		return err
	}

	// Copy the downloaded zip file to the created file
	_, err = io.Copy(zipFile, response.Body)
	if err != nil {
		fmt.Printf("Could not download the steamcmd.zip")
		return err
	}

	// Extract the zip file
	zipReader, err := zip.OpenReader(dlPath)
	if err != nil {
		fmt.Printf("Could not extract the files from steamcmd.zip")
		return err
	}

	// Extract each file in the zip archive
	for _, file := range zipReader.File {
		// Open the file inside the zip
		zipFile, err := file.Open()
		if err != nil {
			fmt.Printf("Could not extract the files from steamcmd.zip")
			return err
		}

		// Create the file in the steamcmd dir
		extractedFilePath := path.Relative("steamcmd", file.Name)
		extractedFile, err := os.Create(extractedFilePath)
		if err != nil {
			fmt.Printf("Could not extract the files from steamcmd.zip")
			zipFile.Close() // Close the zip file handle
			return err
		}

		// Copy the file contents from the zip to the extracted file
		_, err = io.Copy(extractedFile, zipFile)
		if err != nil {
			fmt.Printf("Could not extract the files from steamcmd.zip")
			extractedFile.Close() // Close the extracted file handle
			zipFile.Close()       // Close the zip file handle
			return err
		}

		extractedFile.Close() // Close the extracted file handle
		zipFile.Close()       // Close the zip file handle
	}

	// Close handles before deleting the zip file
	response.Body.Close()
	zipFile.Close()
	zipReader.Close()

	// Delete the zip file
	err = os.Remove(dlPath)
	if err != nil {
		fmt.Println(
			"Could not delete steamcmd.zip.",
			"Please delete it manually at %s",
			dlPath)
		return err
	}

	return nil
}
