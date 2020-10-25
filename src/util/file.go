package util

import (
	"bufio"
	"io"
	"os"
	"strings"
)

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
