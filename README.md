## RecreateCSDownloadsTree
---
### What:
- Given an `rsync`ed directory from a TeamServer, e.g. `/etc/cobaltstrike/downloads`
  - Generate a `TOML` configuration which can be parsed
    - IDs of files and folders are recreated locally
    - The parent directory structure tree is replicated from the current working directory
      - The directory base will start with `SYNCED_DOWNLOADS`
---
#### Why:
- My [existing client-side script](https://github.com/EspressoCake/BeaconDownloadSync) is `async` in nature
- Some performance hits to client-side activity on the `CobaltStrike` application until syncing is done
- No desire to create multiple progress bars to indicate status
---
#### Features:
- [x] `SHA256` integrity checking on every file
  - [x] Only copy if a more recent version of the file from the provided `downloads` directory is present
- [x] Recreate parent directory tree (know where your files came from!)
  - [x] Support for network paths to become the new parent root (start from network share down)
- [x] Only have the most recent version of a disparate file
---
#### Building Instructions:
```sh
cd RecreateCSDownloadsTree
go build
```
---
#### How to Generate the `TOML` Metadata:
- Load the CNA script within the `cna` directory of this repository
- Depending on where you want to execute the underlying logic, consult the chart below:
    | Context | Command  |
    |----------------|-----------------|
    | Beacon         | `toml_download` |
    | Script Console        | `TOMLDownload` |
- The resulting file, `output.toml` will be present within the `cna` directory
    - Copy it to a fresh directory, along with the compiled `golang` application
- Transfer your `downloads` directory from your `TeamServer` however you wish
  - The directory **must** be in the same as wherever your `output.toml` and compiled binary live
---
#### Sample Usage:
```sh
# Display help options
./beacon_download_sync --help

Usage of ./beacon_download_sync:
  -dldir string
        Relative directory name in this current location storing synced Beacon download files (default "downloads")
  -toml string
        TOML file containing metadata associated with Beacon downloads (default "downloads.toml")

# Sample run
./beacon_download_sync -dldir downloads -toml output.toml
```

