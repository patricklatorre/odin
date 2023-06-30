package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Ensures the required dirs and binaries are present
func mustSetup() {
	dirs := []string{
		OdinPath("servers"),
		OdinPath("worlds"),
		OdinPath("steamcmd"),
	}

	// Create required dirs
	for _, dir := range dirs {
		exists, err := Exists(dir)
		Must(err)

		if !exists {
			err := os.Mkdir(dir, 0755)
			Must(err)
			fmt.Println("Created directory:", dir)
		}
	}

	exe := OdinPath("steamcmd", "steamcmd.exe")
	exists, err := Exists(exe)
	Must(err)

	if !exists {
		fmt.Printf("SteamCMD not found, downloading... ")
		err := setupSteamcmd()
		Must(err)
		fmt.Println("done")
	}
}

// Downloads and unpacks steamcmd
func setupSteamcmd() error {
	var (
		dlUrl  = "https://steamcdn-a.akamaihd.net/client/installer/steamcmd.zip"
		dlPath = OdinPath("steamcmd.zip")
	)

	// Download the zip file
	response, err := http.Get(dlUrl)
	if err != nil {
		fmt.Println("Could not download the steamcmd.zip")
		return err
	}

	// Create a new file for writing the zip file
	zipFile, err := os.Create(dlPath)
	if err != nil {
		fmt.Println("Could not download the steamcmd.zip")
		return err
	}

	// Copy the downloaded zip file to the created file
	_, err = io.Copy(zipFile, response.Body)
	if err != nil {
		fmt.Println("Could not download the steamcmd.zip")
		return err
	}

	// Extract the zip file
	zipReader, err := zip.OpenReader(dlPath)
	if err != nil {
		fmt.Println("Could not extract the files from steamcmd.zip")
		return err
	}

	// Extract each file in the zip archive
	for _, file := range zipReader.File {
		// Open the file inside the zip
		zipFile, err := file.Open()
		if err != nil {
			fmt.Println("Could not extract the files from steamcmd.zip")
			return err
		}

		// Create the file in the steamcmd dir
		extractedFilePath := OdinPath("steamcmd", file.Name)
		extractedFile, err := os.Create(extractedFilePath)
		if err != nil {
			fmt.Println("Could not extract the files from steamcmd.zip")
			zipFile.Close()
			return err
		}

		// Copy the file contents from the zip to the extracted file
		_, err = io.Copy(extractedFile, zipFile)
		if err != nil {
			fmt.Println("Could not extract the files from steamcmd.zip")
			extractedFile.Close()
			zipFile.Close()
			return err
		}

		extractedFile.Close()
		zipFile.Close()
	}

	// Close handles before deleting the zip file
	response.Body.Close()
	zipFile.Close()
	zipReader.Close()

	// Delete the zip file
	err = os.Remove(dlPath)
	if err != nil {
		fmt.Println("Could not delete steamcmd.zip:", dlPath)
		return err
	}

	return nil
}
