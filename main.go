package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type PNG struct {
	Prompt         string
	NegativePrompt string
	Steps          int
	Sampler        string
	CFGScale       int
	Seed           int
	Size           string
	ModelHash      string
}

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

func GetMetadata(r io.ReadSeeker) (result string, err error) {
	// 5.2 PNG signature
	const signature = "\x89PNG\r\n\x1a\n"

	// 5.3 Chunk layout
	const crcSize = 4

	// 8 is the size of both the signature and the chunk
	// id (4 bytes) + chunk length (4 bytes).
	// This is just a coincidence.
	buf := make([]byte, 8)

	var n int
	n, err = r.Read(buf)
	if err != nil {
		print("error: ", err)
		return "", err
	}

	if n != len(signature) || string(buf) != signature {
		print("invalid PNG signature")
		return "", errors.New("invalid PNG signature")
	}

	for {
		n, err = r.Read(buf)
		if err != nil {
			break
		}

		if n != len(buf) {
			break
		}

		length := binary.BigEndian.Uint32(buf[0:4])
		chunkType := string(buf[4:8])
		switch chunkType {
		case "tEXt":
			print("found tEXt chunk\n")

			data := make([]byte, length)
			_, err := r.Read(data)
			if err != nil {
				return "", err
			}

			separator := []byte{0}
			separatorIndex := bytes.Index(data, separator)
			if separatorIndex == -1 {
				return "", errors.New("invalid tEXt chunk")
			}
			return string(data[separatorIndex+1:]), nil

		default:
			// Discard the chunk length + CRC.
			_, err := r.Seek(int64(length+crcSize), io.SeekCurrent)
			if err != nil {
				return "", err
			}
		}
	}

	return "", nil
}

func (png *PNG) populateInfo(input string) {
	// Split the input string into lines
	lines := strings.Split(input, "\n")
	// Split the first line into prompt and negative prompt
	prompt := lines[0]
	negativePrompt := lines[1]

	png.Prompt = prompt
	png.NegativePrompt = strings.Trim(strings.Split(negativePrompt, ":")[1], " ")

	nextLines := strings.Split(lines[2], ",")

	// Iterate through the rest of the lines
	for _, line := range nextLines {
		// Split the line into key and value
		parts := strings.Split(line, ":")
		key := strings.Trim(parts[0], " ")
		value := strings.Trim(parts[1], " ")
		// Use a switch statement to set the appropriate field in the struct
		switch key {
		case "Steps":
			i, _ := strconv.Atoi(value)
			png.Steps = i
		case "Sampler":
			png.Sampler = value
		case "CFG scale":
			i, _ := strconv.Atoi(value)
			png.CFGScale = i
		case "Seed":
			i, _ := strconv.Atoi(value)
			png.Seed = i
		case "Size":
			png.Size = value
		case "Model hash":
			png.ModelHash = value
		}
	}

}

func main() {
	var silent bool
	flag.BoolVar(&silent, "s", false, "whether to run in silent mode")
	flag.BoolVar(&silent, "silent", false, "whether to run in silent mode")

	flag.Parse()

	if flag.NArg() < 2 {
		fmt.Println("Usage: remeta [-s|--silent] <mode> <image_file>")
		return
	}

	mode := flag.Arg(0)
	filePath := flag.Arg(1)

	name := filepath.Base(filePath)

	if mode == "remove" {
		err := RemoveMetadata(filePath)
		if err != nil {
			fmt.Println(err)
			return
		}

		if silent {
			return
		}

		a := app.New()

		windowTitle := name + " - ReMeta by geocine"

		// Create a new window
		w := a.NewWindow(windowTitle)
		w.Resize(fyne.NewSize(350, 100))
		w.SetFixedSize(true)
		w.CenterOnScreen()

		content := container.NewCenter(widget.NewLabel("Metadata successfully removed from image."))

		// Set the window content to the horizontal box layout
		w.SetContent(content)
		w.ShowAndRun()

	} else if mode == "read" {
		imgFile, err := os.Open(filePath)
		if err != nil {
			panic(err)
		}
		defer imgFile.Close()
		text, _ := GetMetadata(imgFile)

		if text == "" {
			a := app.New()

			windowTitle := name + " - ReMeta by geocine"

			w := a.NewWindow(windowTitle)
			w.Resize(fyne.NewSize(350, 100))
			w.SetFixedSize(true)
			w.CenterOnScreen()

			content := container.NewCenter(widget.NewLabel("No metadata found in image."))
			w.SetContent(content)
			w.ShowAndRun()
			return
		}

		png := PNG{}
		png.populateInfo(text)

		// Create a new Fyne app
		a := app.New()

		windowTitle := name + " - ReMeta by geocine"
		// Create a new window
		w := a.NewWindow(windowTitle)
		w.Resize(fyne.NewSize(600, 0))
		w.CenterOnScreen()

		// Create a new grid layout
		// grid := container.New(layout.NewGridLayout(2))

		form := &widget.Form{}

		// Create a new label and textbox for each property in the struct

		promptTextbox := widget.NewMultiLineEntry()
		promptTextbox.SetText(png.Prompt)
		promptTextbox.Wrapping = fyne.TextWrapWord
		form.Append("Prompt", promptTextbox)

		negativePromptTextbox := widget.NewMultiLineEntry()
		negativePromptTextbox.SetText(png.NegativePrompt)
		negativePromptTextbox.Wrapping = fyne.TextWrapWord
		form.Append("Negative Prompt", negativePromptTextbox)

		stepsTextbox := widget.NewEntry()
		stepsTextbox.SetText(fmt.Sprint(png.Steps))
		form.Append("Steps", stepsTextbox)

		samplerTextbox := widget.NewEntry()
		samplerTextbox.SetText(png.Sampler)
		form.Append("Sampler", samplerTextbox)

		cfgScaleTextbox := widget.NewEntry()
		cfgScaleTextbox.SetText(fmt.Sprint(png.CFGScale))
		form.Append("CFG Scale", cfgScaleTextbox)

		seedTextbox := widget.NewEntry()
		seedTextbox.SetText(fmt.Sprint(png.Seed))
		form.Append("Seed", seedTextbox)

		sizeTextbox := widget.NewEntry()
		sizeTextbox.SetText(png.Size)
		form.Append("Size", sizeTextbox)

		modelHashTextbox := widget.NewEntry()
		modelHashTextbox.SetText(png.ModelHash)
		form.Append("Model Hash", modelHashTextbox)

		// Set the window content to the grid layout
		w.SetContent(form)

		// Show and run the window
		w.ShowAndRun()
	}
}
