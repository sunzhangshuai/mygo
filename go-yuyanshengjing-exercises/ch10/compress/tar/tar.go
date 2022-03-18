package tar

import (
	"archive/tar"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"exercises/ch10/compress"
)

// DecodeTar 解压文档
func DecodeTar(fileName string) error {
	var file *os.File
	var hfd *tar.Header
	var err error

	// 打开文件
	if file, err = os.Open(fileName); err != nil {
		return err
	}
	defer file.Close()

	// 写文件
	tr := tar.NewReader(file)

	for {
		var file *os.File
		var nr int64

		hfd, err = tr.Next()

		// 判断异常情况
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case hfd == nil:
			continue
		}

		fileName := path.Join(compress.BasePath, hfd.Name)

		// 创建文件
		switch hfd.Typeflag {
		case tar.TypeDir: // 文件
			if !compress.ExistDir(fileName) {
				if err = os.MkdirAll(fileName, 0775); err != nil {
					return err
				}
			}
			fmt.Printf("tar dir %s success\n", hfd.Name)
		case tar.TypeReg: // 文件入磁盘
			if err = func() error {
				if file, err = os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, os.FileMode(hfd.Mode)); err != nil {
					return err
				}
				defer file.Close()
				if nr, err = io.Copy(file, tr); err != nil {
					return err
				}
				fmt.Printf("untar file %s success：%d size\n", hfd.Name, nr)
				return nil
			}(); err != nil {
				return err
			}
		}
	}
}

// EncodeTar 压缩文档
func EncodeTar(dirOrFile string) error {
	var file *os.File
	var err error

	fileName := path.Join(compress.BasePath, path.Base(dirOrFile)+".tar.gz")
	if file, err = os.Create(fileName); err != nil {
		if err != os.ErrExist {
			return err
		}
		if err = file.Truncate(0); err != nil {
			return err
		}
	}

	defer file.Close()

	tw := tar.NewWriter(file)
	defer tw.Close()

	return filepath.Walk(dirOrFile, func(path string, info fs.FileInfo, err error) error {
		var fhd *tar.Header
		var file *os.File
		var nr int64

		if err != nil {
			return err
		}

		// 获取文件头信息
		if fhd, err = tar.FileInfoHeader(info, ""); err != nil {
			return err
		}

		// 去掉前面的"/"
		fhd.Name = strings.TrimPrefix(path, filepath.Dir(dirOrFile))
		fhd.Name = strings.TrimPrefix(fhd.Name, string(filepath.Separator))

		// 写入头
		if err = tw.WriteHeader(fhd); err != nil {
			return err
		}

		// 判断是不是标准文件，目录不写内容
		if !info.Mode().IsRegular() {
			return nil
		}

		// 写文件
		if file, err = os.Open(path); err != nil {
			return err
		}
		defer file.Close()
		if nr, err = io.Copy(tw, file); err != nil {
			return err
		}

		fmt.Printf("tar file %s success，write %d size\n", fhd.Name, nr)
		return nil
	})
}

// init 注册加解压方法
func init() {
	compress.Register("tar", EncodeTar, DecodeTar)
}
