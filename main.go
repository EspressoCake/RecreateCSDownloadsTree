package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"beacon_download_sync/data"
)

var (
	relSyncPath = flag.String("dldir", "downloads", "Relative directory name in this current location storing synced Beacon download files")
	tomlFile    = flag.String("toml", "downloads.toml", "TOML file containing metadata associated with Beacon downloads")
)

func checkPathing(path string) bool {
	cwd, _ := os.Getwd()
	fullPath := fmt.Sprintf("%s/%s", cwd, path)

	if _, err := os.Stat(fullPath); err != nil {
		return false
	}

	return true
}

func main() {
	// Parse the relative download path flag
	flag.Parse()

	// Sanity check the location of synced files before footguns ensue
	if !checkPathing(*relSyncPath) {
		log.Fatal(fmt.Errorf("local directory name supplied was invalid: %s", *relSyncPath))
	}

	// Sanity check user input for a TOML file to be read
	if !checkPathing(*tomlFile) {
		log.Fatal(fmt.Errorf("file provided as TOML input is not present in the current directory: %s", *tomlFile))
	}

	marshal := data.New(*tomlFile, *relSyncPath)

	fmt.Printf("Attempting to integrity check: %d entities\n", len(marshal.Download))

	for index, _ := range marshal.Download {
		marshal.Download[index].ModifyNetPath()
		marshal.Download[index].PartitionFinalLPath()
		marshal.Download[index].PopulateDownloadSHA256(*relSyncPath)

		if marshal.Download[index].SHA256 != "" {
			marshal.Download[index].CheckLocalPaths()
			marshal.Download[index].CheckLocalCopySHA256()
		}
	}
}
