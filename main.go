package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

var ext string

func main() {
	if err := GetAliyunCliPkg(); err != nil {
		log.Fatal(err)
	}
}

func GetAliyunCliPkg() error {
	version := os.Getenv("VERSION")
	if version == "" {
		version = "3.0.55"
	}
	var system string
	switch runtime.GOOS {
	case "windows":
		system = "windows"
		ext = "zip"
	case "darwin":
		system = "macosx"
		ext = "tgz"
	case "linux":
		system = "linux"
		ext = "tgz"
	}
	var url = "https://github.com/aliyun/aliyun-cli/releases/download/v" + version + "/aliyun-cli-" + system + "-" + version + "-amd64." + ext
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if system == "linux" || system == "macosx" {
		file, err := os.Create("aliyun.tgz")
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(file, resp.Body)
		if err != nil {
			return err
		}
		absPath, err := filepath.Abs(file.Name())
		if err != nil {
			return err
		}
		if err := Decompression(absPath); err != nil {
			return err
		}
		cmd := exec.Command("mv", "./aliyun", "/usr/local/bin")
		if err := cmd.Run(); err != nil {
			return err
		}
		return nil
	}
	file, err := os.Create("aliyun.zip")
	if err != nil {
		return err
	}
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	absPath, err := filepath.Abs("aliyun.zip")
	if err != nil {
		return err
	}
	fmt.Println("test", absPath)
	destFile, err := DecompressZip(absPath, "")
	if err != nil {
		return err
	}
	pa, err := filepath.Abs(destFile)
	if err != nil {
		return err
	}
	cmd := exec.Command("dir")
	if err := cmd.Run(); err != nil {
		return err
	}
	out := cmd.Stdout
	fmt.Println(out)
	s := fmt.Sprintf("PATH=%PATH%;" + pa + "\\aliyun.exe")
	cmd = exec.Command("set", s)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func Decompression(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	gr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	tr := tar.NewReader(gr)
	for {
		th, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fi, err := os.OpenFile(th.Name, os.O_CREATE|os.O_WRONLY, 0777)
		if err != nil {
			return err
		}
		defer fi.Close()
		_, err = io.Copy(fi, tr)
		if err != nil {
			return err
		}
	}
	return nil
}

func DecompressZip(fileName, destDir string) (string, error) {
	fmt.Println("zip")
	zr, err := zip.OpenReader(fileName)
	if err != nil {
		fmt.Println("world")
		return "", err
	}
	defer zr.Close()
	for _, f := range zr.File {
		fpath := filepath.Join(destDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return "", err
			}
			inFile, err := f.Open()
			if err != nil {
				return "", err
			}
			defer inFile.Close()
			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return "", err
			}
			defer outFile.Close()
			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return "", err
			}
		}
	}
	return destDir, nil
}
