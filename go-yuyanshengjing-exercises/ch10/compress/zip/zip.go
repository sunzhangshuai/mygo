package zip

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"exercises/ch10/compress"
)

func DecodeZip(filename string) error {
	var zr *zip.ReadCloser
	var err error

	// 写文件
	if zr, err = zip.OpenReader(filename); err != nil {
		return err
	}
	defer zr.Close()

	for _, file := range zr.File {
		var nr int64
		var fileReader io.ReadCloser

		fileName := path.Join(compress.BasePath, file.Name)

		switch {
		case file.Mode().IsDir():
			if !compress.ExistDir(fileName) {
				if err = os.MkdirAll(fileName, 0775); err != nil {
					return err
				}
			}
		case file.Mode().IsRegular():

			if err = func() error {
				var f *os.File
				if fileReader, err = file.Open(); err != nil {
					return err
				}
				defer fileReader.Close()

				// 生成文件
				if f, err = os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, file.Mode()); err != nil {
					return err
				}
				defer f.Close()

				if nr, err = io.Copy(f, fileReader); err != nil {
					return err
				}
				fmt.Printf("unzip file %s success， %d size\n", file.Name, nr)
				return nil
			}(); err != nil {
				return err
			}
		}
	}

	return nil
}

// EncodeZip 压缩文档
func EncodeZip(dirOrFile string) error {
	var file *os.File
	var err error

	fileName := path.Join(compress.BasePath, path.Base(dirOrFile)+".zip")

	os.RemoveAll(fileName)

	// 创建文件
	if file, err = os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0755); err != nil {
		return err
	}
	defer file.Close()

	// 创建写入器

	zw := zip.NewWriter(file)
	defer zw.Close()

	// 递归处理
	return filepath.Walk(dirOrFile, func(path string, info fs.FileInfo, err error) error {
		var fhd *zip.FileHeader
		var wirter io.Writer
		var file *os.File
		var nr int64

		// 如果是源路径，提前进行下一个遍历
		if path == dirOrFile {
			return nil
		}

		if err != nil {
			return err
		}

		// 写文件头
		if fhd, err = zip.FileInfoHeader(info); err != nil {
			return err
		}
		fhd.Name = strings.TrimPrefix(path, filepath.Dir(dirOrFile))
		fhd.Name = strings.TrimPrefix(fhd.Name, string(filepath.Separator))
		if info.IsDir() {
			fhd.Name += "/"
		} else {
			fhd.Method = zip.Deflate
		}

		if wirter, err = zw.CreateHeader(fhd); err != nil {
			return err
		}

		// 目录等不写内容
		if !info.Mode().IsRegular() {
			return nil
		}

		// 写内容
		if file, err = os.Open(path); err != nil {
			return err
		}
		defer file.Close()

		if nr, err = io.Copy(wirter, file); err != nil {
			return err
		}

		fmt.Printf("zip file %s success，write %d size\n", fhd.Name, nr)
		return nil
	})
}

// init 注册加解压方法
func init() {
	compress.Register("zip", EncodeZip, DecodeZip)
}
