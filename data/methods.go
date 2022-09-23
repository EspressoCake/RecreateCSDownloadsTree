package data

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

// DisplayDownloadMetadata is for debugging purposes during development
func (data *TomlWrapper) DisplayDownloadMetadata() {
	for _, item := range data.Download {
		fmt.Printf("%+v\n", item)
	}
}

// DisplayIndividualDownloadMetaData is meant for debugging purposes when creation goes wrong
func (data *DownloadMetadata) DisplayIndividualDownloadMetadata() {
	fmt.Printf("%+v\n", data)
}

// CheckLocalPaths will ensure we create local paths for our downloads to go
func (data *DownloadMetadata) CheckLocalPaths() {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	finalDirPath := fmt.Sprintf("%s/%s", path, data.FinalPath)
	err = os.MkdirAll(finalDirPath, 0750)

	if err != nil {
		log.Fatal(err)
	}
}

// CheckLocalCopySHA256 will run to compare the initial run (Beacon download) against an existing path
func (data *DownloadMetadata) CheckLocalCopySHA256() {
	// This will serve as our file path to resolve
	finalizedPath := strings.Replace(fmt.Sprintf("%s/%s", data.FinalPath, data.Fname), "//", "/", -1)

	if _, err := os.Stat(finalizedPath); errors.Is(err, os.ErrNotExist) {
		currentDownloadPath := fmt.Sprintf("%s/%s", data.DLPath, data.ID)

		download, err := os.Open(currentDownloadPath)
		if err != nil {
			log.Fatal(err)
		}
		defer download.Close()

		finalDownloadPath := strings.Replace(fmt.Sprintf("%s/%s", data.FinalPath, data.Fname), "//", "/", -1)
		destinationFile, err := os.Create(finalDownloadPath)
		if err != nil {
			fmt.Println("File failed to be created/written:", finalDownloadPath)
			data.DisplayIndividualDownloadMetadata()

			return
		}
		defer destinationFile.Close()

		_, err = io.Copy(destinationFile, download)
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	_, err := os.Stat(finalizedPath)
	if err == nil {
		finalDownloadPath := fmt.Sprintf("%s/%s", data.FinalPath, data.Fname)

		previousDownload, err := os.Open(finalDownloadPath)
		if err != nil {
			log.Fatal(err)
		}

		h := sha256.New()
		if _, err := io.Copy(h, previousDownload); err != nil {
			log.Fatal(err)
		}

		if !(fmt.Sprintf("%x", h.Sum(nil)) == data.SHA256) {
			currentDownloadPath := fmt.Sprintf("%s/%s", data.DLPath, data.ID)
			download, err := os.Open(currentDownloadPath)
			if err != nil {
				log.Fatal(err)
			}
			defer download.Close()

			destinationFile, err := os.Create(finalDownloadPath)
			if err != nil {
				log.Fatal(err)
			}
			defer destinationFile.Close()

			_, err = io.Copy(destinationFile, download)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

// PopulateDownloadSHA256 will create the string representation of a Beacon download file's SHA256 sum
func (data *DownloadMetadata) PopulateDownloadSHA256(currentDLPath string) {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	filePath := fmt.Sprintf("%s/%s/%s", path, currentDLPath, data.ID)
	file, err := os.Open(filePath)

	if err != nil {
		return
	}
	defer file.Close()

	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		log.Fatal(err)
	}

	data.SHA256 = fmt.Sprintf("%x", h.Sum(nil))
}

// DetermineBeaconDLExistence will determine if the file within a given TOML is actually present within the local "synced downloads" directory
func (data *DownloadMetadata) DetermineBeaconDLExistence(currentDLPath string) bool {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	filePath := fmt.Sprintf("%s/%s/%s", path, currentDLPath, data.Fname)

	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}

// ModifyNetPath is used to create in-place modifications to the DownloadMetadata structure
func (data *DownloadMetadata) ModifyNetPath() {
	if data.IsNetPath() {
		r := regexp.MustCompile(`//\s*(.*?)\s*/`)
		matches := r.FindStringSubmatch(data.Fpath)

		data.Bname = matches[1]
		data.Fpath = strings.Replace(data.Fpath, `//`+matches[1]+`/`, ``, 1)
	}
}

func (data *DownloadMetadata) PartitionFinalLPath() {
	data.FinalPath = fmt.Sprintf("%s/%s/%s", "SYNCED_DOWNLOADS", data.Bname, data.Fpath)
}

// IsNetPath will determine if the current path is obtained from a network share
func (data *DownloadMetadata) IsNetPath() bool {
	return strings.HasPrefix(data.Fpath, `//`)
}
