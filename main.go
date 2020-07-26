package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"./netFunc"
	"./textFunc"
)

//配置变量
var (
	defaultTargetPath = "./target.html" //显示歌单的网页文本位置
	SavePathTemplate  = "./tmp/%s.mp3"  //保存文件的路径模板
)

func main() {
	//读取html文本内容
	err, html := ReadHtmlFromFile(defaultTargetPath)
	if err != nil {
		fmt.Printf("get target file fail: %v\n", err)
		return
	}
	//分析html文本信息
	err, musicList := textFunc.GetMusicList(html)
	if err != nil {
		fmt.Printf("get music list fail: %v\n", err)
		return
	}
	//开始下载
	for idx, row := range musicList {
		savePath := fmt.Sprintf(SavePathTemplate, row.Name)
		err := netFunc.DownloadSource(row.Url, savePath)
		if err != nil {
			fmt.Printf("%d/%d---FAIL---%s---%v \n", idx+1, len(musicList), row.Name, err)
			continue
		} else {
			fmt.Printf("%d/%d---OK---%s\n", idx+1, len(musicList), row.Name)
		}
	}

	return
}

//读取指定路径的文件
func ReadHtmlFromFile(path string) (err error, html string) {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Open %s fall: %v", path, err), ""
	}
	defer file.Close()
	buf := bufio.NewReader(file)
	bytes, err := ioutil.ReadAll(buf)
	if err != nil {
		return fmt.Errorf("ioutil.ReadAll fall : %v", err), ""
	}
	return nil, fmt.Sprintf("%s", bytes)
}
