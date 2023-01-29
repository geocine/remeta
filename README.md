# ReMeta

ReMeta is a powerful tool for removing and reading metadata from images. With ReMeta, you can easily remove metadata from images using the Windows context menu, making it simple and convenient to use. Additionally, ReMeta can also read AI metadata from images using the Windows context menu, providing you with valuable information about your images.

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