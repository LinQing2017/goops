package sys

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
)

func PrintJSON(appInfo interface{}) {
	if jsonByte, err := json.Marshal(appInfo); err == nil {
		var jsonFormate bytes.Buffer
		if err = json.Indent(&jsonFormate, jsonByte, "", "    "); err != nil {
			logrus.Error("解析json字符传异常")
		} else {
			fmt.Println(jsonFormate.String())
		}
	} else {
		logrus.Error("解析json字符传异常")
	}
}
