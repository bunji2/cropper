package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

// Conf は設定データ
type Conf struct {
	P1     [2]int `json:"p1"`      // 左上の座標
	P2     [2]int `json:"p2"`      // 右下の座標
	OutDir string `json:"out_dir"` // 出力先ディレクトリ
}

func loadJSON(filePath string) (r Conf, err error) {
	var bb []byte
	bb, err = ioutil.ReadFile(filepath.Clean(filePath)) // for CWE-22
	if err != nil {
		return
	}
	err = json.Unmarshal(bb, &r)
	if err != nil {
		return
	}

	// 設定データの値のチェック
	if r.P2[0] < 0 || r.P2[1] < 0 || r.P1[0] < 0 || r.P1[1] < 0 {
		err = fmt.Errorf("coodinate of config is negative value")
		return
	}

	return
}
