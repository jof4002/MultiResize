package main // import github.com/jof4002/MultiResize

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/nfnt/resize"
)

type OutputInfo struct {
	OutPath string `json:"path"`
	Size    string `json:"size"`
}

type Outputs []OutputInfo

func main() {
	jsonPathPtr := flag.String("config", "", "config json path")
	outputPathPtr := flag.String("output", "", "output base path")
	imagePathPtr := flag.String("image", "", "original image path")

	flag.Parse()

	jsonPath := *jsonPathPtr
	outputPath := *outputPathPtr
	imagePath := *imagePathPtr

	if jsonPath == "" || outputPath == "" || imagePath == "" {
		// jsonPath = "d:/Projects/Go/src/jof4002/MultiResize/example.json"
		// imagePath = "f:/Perforce/Project/Project.png"
		// outputPath = "f:/Perforce/Project/"
		flag.CommandLine.Usage()
		return
	}

	// read image
	imageFile, err := os.Open(imagePath)
	if err != nil {
		log.Fatal(err)
		time.Sleep(10 * time.Second)
		return
	}
	defer imageFile.Close()

	// decode image
	image, _, err := image.Decode(imageFile)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// read json
	jsonContent, err := os.ReadFile(jsonPath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	outputs := Outputs{}
	err = json.Unmarshal([]byte(jsonContent), &outputs)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, output := range outputs {

		// resize size
		wh := strings.Split(output.Size, "x")
		w, _ := strconv.Atoi(wh[0])
		h := w
		if len(wh) > 1 {
			h, _ = strconv.Atoi(wh[1])
		}

		// resize
		m := resize.Resize(uint(w), uint(h), image, resize.Lanczos3)

		// output file's path and ext
		outfilePath := path.Join(outputPath, output.OutPath)
		os.MkdirAll(filepath.Dir(outfilePath), os.ModePerm)
		ext := path.Ext(outfilePath)

		// create output file
		outFile, err := os.Create(outfilePath)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer outFile.Close()

		// encode
		if ext == ".png" {
			png.Encode(outFile, m)
		} else { // force jpg
			jpeg.Encode(outFile, m, nil)
		}

		fmt.Println("Saved : " + output.OutPath)
	}

}
