package main

import (
	//"image"
	//"image/color"
	//"log"
	//"github.com/disintegration/imaging"
	"fmt"
	"github.com/disintegration/imaging"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {

	// static will be pushed into arguments
	prefix := "cornelius_car_show"

	// config
	incomingPath := "incoming"

	// start reading folder
	imageList := getImages(incomingPath)

	// image count
	imageCount := 0

	// check images
	imageCheck, imageCount := checkImages(imageList)

	if imageCheck {
		fmt.Printf("We found %d images\n", imageCount)
		// create a new folder with prefix
		makeFolder(prefix)

		for _, image := range imageList {
			stringCounter := strconv.Itoa(imageCount)
			destinationFile := "out/" + prefix + "/" + prefix + "_" + stringCounter + ".jpg"
			_, err := imageCopy("incoming/"+image, destinationFile)
			if err != nil {
				//TODO: add logging
				fmt.Println(err)
			}
			// TODO: find different loop?
			imageCount++
			resizeImage(destinationFile)
		}

	} else {
		fmt.Println("There are no images in the incoming folder")
	}
}

func checkImages(incoming []string) (bool, int) {

	if len(incoming) != 0 {
		return true, len(incoming)
	} else {
		return false, 0
	}
}

func getImages(incomingPath string) []string {

	var incoming []string

	// walk the folder looking for images
	err := filepath.Walk(incomingPath, func(path string, info os.FileInfo, err error) error {
		fileName := strings.Split(path, "/")
		// the folder itself is considered a filename
		if len(fileName) > 1 {
			incoming = append(incoming, fileName[1])
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return incoming
}

func makeFolder(folderName string) {

	outFolder := "out/" + folderName

	if _, err := os.Stat(outFolder); os.IsNotExist(err) {
		err = os.MkdirAll(outFolder, 0755)
		if err != nil {
			//TODO: Logs
			panic(err)
		}
	}
}

func imageCopy(src, dst string) (int64, error) {

	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func resizeImage(fileName string) {

	src, err := imaging.Open(fileName)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	// create the thumbnail with a suffix
	newFileName := strings.TrimRight(fileName, ".jpg") + "_t.jpg"

	// Resize the cropped image to width = 200px preserving the aspect ratio.
	src = imaging.Resize(src, 200, 0, imaging.Lanczos)

	// Save the resulting image as JPEG.
	err = imaging.Save(src, newFileName)
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}
}

func buildHtml(prefix string, imageCount int) {

}