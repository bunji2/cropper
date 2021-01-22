package main

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

func cropSave(img image.Image, outFile string, p1, p2 [2]int) (err error) {
	cimg := img.(SubImager).SubImage(image.Rect(p1[0], p1[1], p2[0], p2[1]))

	var fso *os.File
	// 出力先ファイルの書き込みオープン
	fso, err = os.Create(filepath.Clean(outFile))
	if err != nil {
		return
	}
	defer func() {
		e := fso.Close()
		if err == nil {
			err = e
		}
	}()

	ext, _ := getExt(outFile)
	switch ext { // 出力画像の拡張子によってエンコードを変える
	case "jpg", "jpeg":
		err = jpeg.Encode(fso, cimg, &jpeg.Options{Quality: 100})
	case "png":
		err = png.Encode(fso, cimg)
	case "gif":
		err = gif.Encode(fso, cimg, nil)
	default:
		err = fmt.Errorf("unknown extension: %s", ext)
	}

	return
}

// 受付可能な画像の拡張子
var acceptableExts = []string{
	"jpg",
	"jpeg",
	"png",
	"gif",
}

func getExt(filePath string) (ext string, ok bool) {
	slice := strings.Split(filepath.Base(filePath), ".")
	ext = strings.ToLower(slice[len(slice)-1]) // スライスの最後が拡張子のはず
	for _, accebtableExt := range acceptableExts {
		if ext == accebtableExt {
			ok = true
			break
		}
	}
	return
}
