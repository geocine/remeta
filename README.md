# ReMeta

ReMeta is a powerful tool for removing and reading metadata from PNG images generated by the Stable Diffusion Web UI.

### Reading Metadata
https://user-images.githubusercontent.com/507464/215328728-ba06be72-a0e0-4252-8e98-11b00f526ab6.mp4

### Removing Metadata
https://user-images.githubusercontent.com/507464/215329121-06af4bf1-341a-4970-8277-103b44e10a8b.mp4


> Note: This tool was specifically built for the Stable Diffusion Web UI. It may not work with other PNG images. This tool is only tested in Windows.

## Building

You need to have `go1.20rc1` 
```
go install golang.org/dl/go1.20rc1@latest
go1.20rc1 download
```
Then build using the following command:
```sh
.\build.bat
```

## Usage

ReMeta is a CLI tool so you can use it from the command line 
Removing metadata from images. This will create a new image with the same name but with the suffix `_raw`:
```sh
remeta remove <image>
```
Reading metadata from images. This will open a window with the metadata:
```sh
remeta read <image>
```

## Installation

Build or download ReMeta from the Releases page and place it in any folder.

Modify path in to `remeta.exe` in `add.reg`. Replace `D:\\SW\\bin\\` with your path to `remeta.exe`. Then run `add.reg` to add ReMeta to the Windows context menu.

## Uninstallation

Run `remove.reg` to remove ReMeta from the Windows context menu.


