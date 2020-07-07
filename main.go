package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/otiai10/gosseract"
	"github.com/pkg/errors"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var ocrClient = gosseract.NewClient()
var ocrLock sync.Mutex

var mapping sync.Map

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	// prepare folders

	defer ocrClient.Close()

	// init sdl ttf
	must(ttf.Init())

	if len(os.Args) == 1 {
		fmt.Println(`usage: ocr_woff <font_file>`)
		return
	}

	fontPath := os.Args[1]
	imgPath := "tmp.png"

	must(renderFont(fontPath, imgPath))

	result := ocrText(imgPath)

	if isOCRValid(result) {
		fmt.Printf("ocr result: %s\n", result)
	} else {
		fmt.Printf("ocr fail")
	}
}

func renderFont(fontPath, imgPath string) error {
	font, err := ttf.OpenFont(fontPath, 50)

	if err != nil {
		return errors.Wrap(err, "could not open font")
	}

	surface, err := font.RenderUTF8Solid("0123456789", sdl.Color{0, 0, 0, 0})
	if err != nil {
		return errors.Wrap(err, "could not render text")
	}

	return img.SavePNG(surface, imgPath)
}

func ocrText(imgPath string) string {
	ocrLock.Lock()
	defer ocrLock.Unlock()

	ocrClient.SetWhitelist("0123456789")
	ocrClient.SetImage(imgPath)
	text, _ := ocrClient.Text()
	return text
}

// should a string of length 10 with 0 ~ 9
func isOCRValid(str string) bool {
	return len(str) == 10 &&
		strings.ContainsRune(str, '0') &&
		strings.ContainsRune(str, '1') &&
		strings.ContainsRune(str, '2') &&
		strings.ContainsRune(str, '3') &&
		strings.ContainsRune(str, '4') &&
		strings.ContainsRune(str, '5') &&
		strings.ContainsRune(str, '6') &&
		strings.ContainsRune(str, '7') &&
		strings.ContainsRune(str, '8') &&
		strings.ContainsRune(str, '9')
}
