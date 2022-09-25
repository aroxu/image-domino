package main

import (
	"bufio"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/noelyahan/impexp"
	"github.com/noelyahan/mergi"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	worker(pwd)
}

func worker(folderPath string) {
	fmt.Printf("Reading %s...\n", folderPath)
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		panic(err)
	}

	images := []string{}
	_folderName := folderPath

	for _, file := range files {
		if file.IsDir() {
			_folderName := folderPath
			if runtime.GOOS == "windows" {
				_folderName = _folderName + "\\" + file.Name() + "\\"
			} else {
				_folderName = _folderName + "/" + file.Name() + "/"
			}
			worker(_folderName)
		} else {
			_fileExtension := filepath.Ext(file.Name())
			if file.Name() != "out"+_fileExtension && (_fileExtension == ".jpeg" || _fileExtension == ".jpg" || _fileExtension == ".png") {
				images = append(images, _folderName+file.Name())
			}
		}
	}
	if len(images) > 1 {
		processImage(_folderName, images)
	}
}

func processImage(folder string, images []string) {
	imageFiles := []image.Image{}

	align := "T"

	for _, image := range images {
		asdf, err := mergi.Import(impexp.NewFileImporter(image))
		fmt.Println("Loading image: " + image)
		if err != nil {
			panic(err)
		}
		imageFiles = append(imageFiles, asdf)
	}
	for i := 1; i < len(images); i++ {
		align += "B"
	}

	_path := folder
	if runtime.GOOS == "windows" {
		_path = _path + "\\"
	} else {
		_path = _path + "/"
	}

	fmt.Printf("Found %d Image(s)\n", len(images))

	// fileExtension := filepath.Ext(images[0])
	outFilePath := fmt.Sprintf("%sout%s", _path, ".png")
	fmt.Printf("Generating: %s\n", outFilePath)

	outFile, err := os.Create(outFilePath)
	if err != nil {
		panic(err)
	}
	outFileExt := filepath.Ext(outFilePath)
	defer outFile.Close()
	filewriter := bufio.NewWriter(outFile)
	mergedImageFile, err := mergi.Merge(align, imageFiles)

	if outFileExt == ".jpg" || outFileExt == ".jpeg" {
		if err := jpeg.Encode(filewriter, mergedImageFile, nil); err != nil {
			panic(err)
		}
	} else {
		if err := png.Encode(filewriter, mergedImageFile); err != nil {
			panic(err)
		}
	}

	if err != nil {
		panic(err)
	}
	runtime.GC()
}
