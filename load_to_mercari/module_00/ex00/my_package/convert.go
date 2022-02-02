//与えられたdirectoryの中を探索して，画像のフォーマットを変換するプログラムのパッケージです．
package my_package

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//コマンドライン引数の確認，pathの存在確認を行い，画像変換の関数に渡す
func Convert() bool {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "error: invalid argument")
		return false
	}

	for _, dirs := range flag.Args() {
		files := Dirwalk(dirs)
		if files == nil {
			return false
		}
		for _, file_path := range files {
			index := strings.LastIndex(file_path, ".")
			if index < 0 || !(".jpg" == file_path[index:] || ".png" == file_path[index:]) {
				fmt.Fprintf(os.Stderr, "error: %s is not a valid file\n", file_path)
				return false
			} else {
				ConvertPngToJpg(file_path)
				ConvertJpgToPng(file_path)
			}
		}
	}
	return true
}

func ConvertJpgToPng(file_path string) {
	index := strings.LastIndex(file_path, ".")
	if index < 0 {
		return
	}
	if ".jpg" == file_path[index:] {
		file, err := os.Open(file_path)
		if err != nil {
			panic("1 : " + err.Error())
		}
		img, _, err := image.Decode(file)
		if err != nil {
			panic("2 : " + err.Error())
		}
		file.Close()
		filename := strings.Split(file_path, ".")
		file_path = filename[0] + ".png"
		out, err := os.Create(file_path)
		if err != nil {
			panic("3 : " + err.Error())
		}
		png.Encode(out, img)
		out.Close()
	}
}

func ConvertPngToJpg(file_path string) {
	index := strings.LastIndex(file_path, ".")
	if index < 0 {
		return
	}
	if ".png" == file_path[index:] {
		file, err := os.Open(file_path)
		if err != nil {
			panic("1 : " + err.Error())
		}
		img, err := png.Decode(file)
		if err != nil {
			panic("2 : " + err.Error())
		}
		file.Close()
		filename := strings.Split(file_path, ".")
		file_path = filename[0] + ".jpg"
		out, err := os.Create(file_path)
		if err != nil {
			panic("3 : " + err.Error())
		}
		jpeg.Encode(out, img, nil)
		out.Close()
	}
}

func Dirwalk(dir string) []string {
	files, error := ioutil.ReadDir(dir)
	if error != nil {
		fmt.Fprintf(os.Stderr, "error: %s: no such file or directory\n", dir)
		return nil
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, Dirwalk(filepath.Join(dir, file.Name()))...)
			continue
		}
		paths = append(paths, filepath.Join(dir, file.Name()))
	}
	return paths
}
