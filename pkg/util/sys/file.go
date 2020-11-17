package sys

import (
	"archive/tar"
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// 计算目录的大小
func CalDirSize(dirpath string) (dirsize int64, err error) {
	err = os.Chdir(dirpath)
	if err != nil {
		return
	}
	files, err := ioutil.ReadDir(".")
	if err != nil {
		return
	}

	for _, file := range files {
		dirsize += file.Size()
		if file.Mode().IsDir() {
			subDirSize, err2 := CalDirSize(path.Join(dirpath, file.Name()))
			if err2 == nil {
				dirsize += subDirSize
			}
		}
	}
	return
}

// 按行读取文件
func ReadLine(filePth string) []string {
	f, err := os.Open(filePth)
	lines := make([]string, 0)
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()

	bfRd := bufio.NewReader(f)
	for {
		b, _, err := bfRd.ReadLine()
		if err != nil { //遇到任何错误立即返回，并忽略 EOF 错误信息
			if err == io.EOF {
				break
			} else {
				panic(err.Error())
			}
		}
		lineStr := strings.TrimSpace(string(b))
		if !strings.EqualFold(lineStr, "") {
			lines = append(lines, lineStr)
		}

	}
	return lines
}

// 将指定文件目录及其子目录的内容压缩成tar，写入到文件流中
func MakeTar(srcPath, destPath string, writer io.Writer) error {
	// TODO: use compression here?
	tarWriter := tar.NewWriter(writer)
	defer tarWriter.Close()

	srcPath = path.Clean(srcPath)
	destPath = path.Clean(destPath)
	return recursiveTar(path.Dir(srcPath), path.Base(srcPath), path.Dir(destPath), path.Base(destPath), tarWriter)
}

func recursiveTar(srcBase, srcFile, destBase, destFile string, tw *tar.Writer) error {
	srcPath := path.Join(srcBase, srcFile)
	matchedPaths, err := filepath.Glob(srcPath)
	if err != nil {
		return err
	}
	for _, fpath := range matchedPaths {
		stat, err := os.Lstat(fpath)
		if err != nil {
			return err
		}
		if stat.IsDir() {
			files, err := ioutil.ReadDir(fpath)
			if err != nil {
				return err
			}
			if len(files) == 0 {
				//case empty directory
				hdr, _ := tar.FileInfoHeader(stat, fpath)
				hdr.Name = destFile
				if err := tw.WriteHeader(hdr); err != nil {
					return err
				}
			}
			for _, f := range files {
				if err := recursiveTar(srcBase, path.Join(srcFile, f.Name()), destBase, path.Join(destFile, f.Name()), tw); err != nil {
					return err
				}
			}
			return nil
		} else if stat.Mode()&os.ModeSymlink != 0 {
			//case soft link
			hdr, _ := tar.FileInfoHeader(stat, fpath)
			target, err := os.Readlink(fpath)
			if err != nil {
				return err
			}

			hdr.Linkname = target
			hdr.Name = destFile
			if err := tw.WriteHeader(hdr); err != nil {
				return err
			}
		} else {
			//case regular file or other file type like pipe
			hdr, err := tar.FileInfoHeader(stat, fpath)
			if err != nil {
				return err
			}
			hdr.Name = destFile

			if err := tw.WriteHeader(hdr); err != nil {
				return err
			}

			f, err := os.Open(fpath)
			if err != nil {
				return err
			}
			defer f.Close()

			if _, err := io.Copy(tw, f); err != nil {
				return err
			}
			return f.Close()
		}
	}
	return nil
}
