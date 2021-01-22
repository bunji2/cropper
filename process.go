package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
)

// SubImager は部分画像のインターフェース
type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

func process(conf Conf, inFile, outFile string) (err error) {
	// 入力・出力ファイルの拡張子のチェック
	ext, ok := getExt(inFile)
	if !ok {
		err = fmt.Errorf("unknown extension: %s", ext)
		return
	}
	ext, ok = getExt(outFile)
	if !ok {
		err = fmt.Errorf("unknown extension: %s", ext)
		return
	}

	// 入力元ファイルの読み出しオープン
	f, err := os.Open(filepath.Clean(inFile))
	if err != nil {
		return
	}

	defer func() {
		e := f.Close()
		if err == nil {
			err = e
		}
	}()

	// 入力ファイルを画像オブジェクトにデコード
	img, _, err := image.Decode(f)
	if err != nil {
		return
	}

	// [TODO]
	// config で設定された矩形領域が画像オブジェクトの中に入ってるかチェック

	// Crop & Save
	err = cropSave(img, outFile, conf.P1, conf.P2)

	return
}
