// 画像の一部分を切り出すサンプル

package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	usageFmt     = "Usage: %s infile_pattern [ out_dir ]\n"
	confFileName = "conf.json"
)

const (
	_ = iota
	argError
	confFileError
	runtimeError
)

func main() {
	os.Exit(run())
}

func run() (exitCode int) {

	// 引数チェック
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, usageFmt, os.Args[0])
		exitCode = argError
		return
	}

	// 設定ファイルパスの取得
	confFilePath, err := resolveConfFile()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = confFileError
		return
	}

	// 設定ファイルの読み出し
	conf, err := loadJSON(confFilePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = confFileError
		return
	}

	// 引数で指定された入力ファイルのパターン（ワイルドカード）を展開し
	// 入力ファイルリストを取得
	inFiles, err := filepath.Glob(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = argError
		return
	}

	//fmt.Println("inFiles =", inFiles)

	// 出力先フォルダの取得

	outDir := conf.OutDir // デフォルトは設定ファイルの値

	if len(os.Args) > 2 { // 引数で出力先が指定された場合
		outDir, err = resolveOutDir(os.Args[2])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			exitCode = argError
			return
		}
	}

	// 出力先が指定されていない場合
	if outDir == "" {
		fmt.Fprintln(os.Stderr, "out_dir is empty in config or argument")
		fmt.Fprintf(os.Stderr, usageFmt, os.Args[0])
		exitCode = confFileError
		return
	}

	// 入力ファイルを順番に処理
	for _, inFile := range inFiles {
		// 出力ファイルパスの作成
		outFile := filepath.Join(outDir, filepath.Base(inFile))
		err = process(conf, inFile, outFile)
		if err != nil {
			break
		}
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = runtimeError
	}

	return
}

// resolveConfFile は設定ファイルパスを決定する関数
// 実行ファイルと同じディレクトリ配下の "conf.json"
func resolveConfFile() (r string, err error) {
	var exePath string
	exePath, err = os.Executable()
	if err != nil {
		return
	}
	r = filepath.Join(filepath.Dir(exePath), confFileName)
	return
}

// resolveOutDir は出力先ディレクトリのパスを決定する関数
func resolveOutDir(outDir string) (r string, err error) {
	// 絶対パスの時
	if filepath.IsAbs(outDir) {
		r = outDir
		return
	}

	// 以下、相対パスの時＝実行ファイルと同じディレクトリ配下
	var exePath string
	exePath, err = os.Executable()
	if err != nil {
		return
	}
	r = filepath.Join(filepath.Dir(exePath), outDir)

	// 出力先ディレクトリの存在確認
	f, e := os.Stat(r)    // [XXX] e の値で分岐しなくて大丈夫かな？
	if os.IsNotExist(e) { // 存在しない時
		//err = fmt.Errorf("out_dir: %s does not exist", r)
		err = os.MkdirAll(r, 0750) // ないなら作成しちゃった方が使いやすいのでは
		if err != nil {
			return
		}
	} else { // 存在する時にはディレクトリかどうか確認
		if f != nil && !f.IsDir() { // [MEMO] f!=nil の条件は不要かも？
			err = fmt.Errorf("out_dir: %s is not directory", r)
		}
	}

	return
}
