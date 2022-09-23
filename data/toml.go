package data

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml"
)

// TomlWrapper is a key-value pair of embedded structs representing Beacon download metadata
type TomlWrapper struct {
	Download []DownloadMetadata `toml:"download"`
}

// DownloadMetadata is a struct representing Beacon download metadata
type DownloadMetadata struct {
	Bname     string `toml:"bname"`
	Fname     string `toml:"fname"`
	Fpath     string `toml:"fpath"`
	Host      string `toml:"host" omitempty:"true"`
	ID        string `toml:"id"`
	SHA256    string
	DLPath    string
	FinalPath string
}

func New(tomlpath string, dlpath string) *TomlWrapper {
	var data TomlWrapper

	fbytes, _ := os.ReadFile(tomlpath)
	_ = toml.Unmarshal(fbytes, &data)

	for index, _ := range data.Download {
		data.Download[index].DLPath = fmt.Sprintf("./%s", dlpath)
	}

	return &data
}
