package sys

import (
	"fmt"
	"github.com/nwidger/jsoncolor"
	"github.com/sirupsen/logrus"
)

func PrintJSON(appInfo interface{}) {

	if jsonByte, err := jsoncolor.MarshalIndent(appInfo, "", "    "); err == nil {
		fmt.Println(string(jsonByte))
	} else {
		logrus.Error("解析json字符传异常")
	}
}
