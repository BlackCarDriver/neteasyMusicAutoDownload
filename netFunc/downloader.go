package netFunc

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

//配置变量
var (
	ignoreSize = 1024 //小于多kb的响应主体放弃保存
)

//将url的主体数据保存到指定的文件之中
func DownloadSource(imgUrl, fileName string) error {
	resp, err := http.Get(imgUrl)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	bodySize := len(body)
	if bodySize == 0 {
		return fmt.Errorf("body size = 0")
	}

	imgPath := fmt.Sprintf("./%s", fileName)

	//响应主体非常小,大概率不为mp3文件,改变其后缀以容易识别
	if bodySize>>10 < ignoreSize {
		imgPath += ".err"
	}

	out, err := os.Create(imgPath)
	defer out.Close()
	if err != nil {
		return fmt.Errorf("create file fail: %v", err)
	}

	_, err = io.Copy(out, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("save file fail: %v", err)
	}

	if bodySize>>10 < ignoreSize {
		return fmt.Errorf("it file maybe not a music")
	}

	return nil
}
