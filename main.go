package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"github.com/dsoprea/go-exif"
	pngstructure "github.com/dsoprea/go-png-image-structure"
	"golang.design/x/clipboard"
)

type IfdEntry struct {
	IfdPath     string                `json:"ifd_path"`
	FqIfdPath   string                `json:"fq_ifd_path"`
	IfdIndex    int                   `json:"ifd_index"`
	TagId       uint16                `json:"tag_id"`
	TagName     string                `json:"tag_name"`
	TagTypeId   exif.TagTypePrimitive `json:"tag_type_id"`
	TagTypeName string                `json:"tag_type_name"`
	UnitCount   uint32                `json:"unit_count"`
	Value       interface{}           `json:"value"`
	ValueString string                `json:"value_string"`
}
type IfdEntries map[string]IfdEntry

var (
	user32                  = syscall.NewLazyDLL("user32.dll")
	messageBoxW             = user32.NewProc("MessageBoxW")
	MB_OK                   = 0x00000000
	MB_ICONINFORMATION      = 0x00000040
	MB_SETFOREGROUND        = 0x00010000
	MB_TOPMOST              = 0x00040000
	MB_SERVICE_NOTIFICATION = 0x00200000
)

func RemoveMetadata(filePath string) error {

	// Open the image file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	// Create a new file path with "_raw" appended to the base name
	// Extract the file name and extension from the file path
	elements := strings.Split(filePath, string(filepath.Separator))
	fileName := elements[len(elements)-1]
	ext := filepath.Ext(fileName)

	// Create a new file path with "_raw" appended to the base name
	dir := filepath.Dir(filePath)
	rawFilePath := filepath.Join(dir, fileName[:len(fileName)-len(ext)]+"_raw"+ext)

	// Create the output file
	outFile, err := os.Create(rawFilePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Encode the image with no metadata
	return png.Encode(outFile, img)
}

func ReadExif(filePath string) (exifData IfdEntries, err error) {

	pmp := pngstructure.NewPngMediaParser()
	cs, _ := pmp.ParseFile(filePath)
	_, rawExif, exifErr := cs.Exif()
	if exifErr != nil {
		return nil, exifErr
	}

	print(rawExif)

	// Run the parse.
	// im := exif.NewIfdMappingWithStandard()
	// ti := exif.NewTagIndex()

	// entries := make(IfdEntries, 0)
	// visitor := func(fqIfdPath string, ifdIndex int, tagId uint16,
	// 	tagType exif.TagType, valueContext exif.ValueContext) (err error) {
	// 	defer func() {
	// 		if state := recover(); state != nil {
	// 			err = log.Wrap(state.(error))
	// 			log.Panic(err)
	// 		}
	// 	}()

	// 	ifdPath, pathErr := im.StripPathPhraseIndices(fqIfdPath)
	// 	if pathErr != nil {
	// 		return pathErr
	// 	}

	// 	it, tagErr := ti.Get(ifdPath, tagId)
	// 	if tagErr != nil {
	// 		if log.Is(tagErr, exif.ErrTagNotFound) {
	// 			fmt.Printf("WARNING: Unknown tag: [%s] (%04x)\n",
	// 				ifdPath, tagId)
	// 			return nil
	// 		} else {
	// 			return tagErr
	// 		}
	// 	}

	// 	valueString := ""
	// 	var value interface{}
	// 	if tagType.Type() == exif.TypeUndefined {
	// 		var undefErr error
	// 		value, undefErr = valueContext.Undefined()
	// 		if undefErr != nil {
	// 			if undefErr == exif.ErrUnhandledUnknownTypedTag {
	// 				value = nil
	// 			} else {
	// 				return nil
	// 			}
	// 		}
	// 		valueString = fmt.Sprintf("%v", value)
	// 	} else if tagType.Type() == exif.TypeByte {
	// 		var byteErr error
	// 		value, byteErr = valueContext.ReadBytes()
	// 		if byteErr != nil {
	// 			return byteErr
	// 		}
	// 	}
	// 	var formatErr error
	// 	valueString, formatErr = valueContext.FormatFirst()
	// 	if formatErr != nil {
	// 		return formatErr
	// 	}
	// 	if tagType.Type() == exif.TypeAscii {
	// 		value = valueString
	// 	}

	// 	entry := IfdEntry{
	// 		IfdPath:     ifdPath,
	// 		FqIfdPath:   fqIfdPath,
	// 		IfdIndex:    ifdIndex,
	// 		TagId:       tagId,
	// 		TagName:     it.Name,
	// 		TagTypeId:   tagType.Type(),
	// 		TagTypeName: tagType.Name(),
	// 		UnitCount:   valueContext.UnitCount(),
	// 		Value:       value,
	// 		ValueString: valueString,
	// 	}
	// 	entries[it.Name] = entry
	// 	return nil
	// }

	// _, visitErr := exif.Visit(exif.IfdStandard, im, ti, rawExif, visitor)
	// if visitErr != nil {
	// 	return nil, visitErr
	// }
	// return entries, nil
	return make(IfdEntries, 0), nil
}

func GetMetadata(filePath string) error {

	// Read the EXIF data from the PNG data
	exifData, err := ReadExif(filePath)
	if err != nil {
		return err
	}

	// Print the EXIF data
	fmt.Println(exifData)

	// fmt.Sprintf("%v", focalLength)
	clipboard.Write(clipboard.FmtText, []byte("Hello"))

	return nil

}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: remeta <mode> <image_file>")
		return
	}

	mode := os.Args[1]
	filePath := os.Args[2]

	if mode == "clear" {
		err := RemoveMetadata(filePath)
		if err != nil {
			fmt.Println(err)
			return
		}

		title, _ := syscall.UTF16PtrFromString("ReMeta - geocine")
		message, _ := syscall.UTF16PtrFromString("Metadata successfully removed from image.")
		flags := MB_OK | MB_ICONINFORMATION | MB_SETFOREGROUND | MB_TOPMOST | MB_SERVICE_NOTIFICATION
		messageBoxW.Call(0, uintptr(unsafe.Pointer(message)), uintptr(unsafe.Pointer(title)), uintptr(flags))
	} else if mode == "get" {
		err := GetMetadata(filePath)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
