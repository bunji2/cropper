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

// process はファイル単位での処理を行う関数
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

	// 入力元ファイルと出力先ファイルのディレクトリが同じかチェック
	if isSameFolder(inFile, outFile) {
		err = fmt.Errorf("input and output are the same folder")
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

	// config で設定された矩形領域が画像オブジェクトの中に入ってるかチェック
	bounds := img.Bounds()
	if conf.P1[0] < bounds.Min.X || conf.P1[1] < bounds.Min.Y ||
		conf.P2[0] > bounds.Max.X || conf.P2[1] > bounds.Max.Y {
		fmt.Fprintf(os.Stderr, "#WARN cropping area (%d,%d)-(%d,%d) is not included in image bounds (%d,%d)-(%d,%d) of %s\n",
			conf.P1[0], conf.P1[1], conf.P2[0], conf.P2[1],
			bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Max.Y, inFile)
	}
	// Crop & Save
	err = cropSave(img, outFile, conf.P1, conf.P2)

	return
}

// isSameFolder 二つのパスのディレクトリが同じかどうか調べる関数
func isSameFolder(path1, path2 string) (ok bool) {
	dir1, err := filepath.Abs(filepath.Dir(path1))
	if err != nil {
		return
	}
	dir2, err := filepath.Abs(filepath.Dir(path2))
	if err != nil {
		return
	}
	if dir1 == dir2 {
		ok = true
	}
	return
}
