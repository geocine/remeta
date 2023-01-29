# ReMeta

ReMeta is a powerful tool for removing and reading metadata from PNG images generated by the Stable Diffusion Web UI. 

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

## Installation

Modify path in to `remeta.exe` in `add.reg`. Replace `D:\\SW\\bin\\` with your path to `remeta.exe`. Then run `add.reg` to add ReMeta to the Windows context menu.

## Uninstallation

Run `remove.reg` to remove ReMeta from the Windows context menu.

## Usage

You can use ReMeta as a stanalone CLI tool if you don't want to use it from the context menu.
Removing metadata from images:
```sh
remeta clear <image>
```
Reading metadata from images:
```sh
remeta get <image>
```
