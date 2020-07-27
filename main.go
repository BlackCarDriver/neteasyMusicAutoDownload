package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"./netFunc"
	"./textFunc"
)

//配置变量
var (
	defaultTargetPath = "./target.html" //显示歌单的网页文本位置
	SavePath          = "./tmp"         //文件保存的路径
	SavePathTemplate  = "./tmp/%s.mp3"  //保存文件的路径模板
)

func main() {
	//初始化检查
	if err := initCheck(); err != nil {
		fmt.Printf("Init fail: %v\n", err)
	}

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
	countSuc, countErr := 0, 0
	for idx, row := range musicList {
		savePath := fmt.Sprintf(SavePathTemplate, row.Name)
		err := netFunc.DownloadSource(row.Url, savePath)
		if err != nil {
			fmt.Printf("%d/%d---FAIL---%s---%v \n", idx+1, len(musicList), row.Name, err)
			countErr++
		} else {
			fmt.Printf("%d/%d---OK---%s\n", idx+1, len(musicList), row.Name)
			countSuc++
		}
	}
	fmt.Printf("download complete! %d success, %d fail \n", countSuc, countErr)

	//倒数结束
	second := 10
	for _ = range time.NewTicker(time.Second).C {
		fmt.Printf("Exist after %d second...\n", second)
		second--
		if second <= 0 {
			break
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

//检查目录是否存在
func checkDirExist(dir string) bool {
	info, err := os.Stat(dir)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Printf("error with file exist: %v", err)
		}
		return false
	}
	return info.IsDir()
}

//检查文件是否存在
func CheckFileExist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}

//初始化检查
func initCheck() error {
	if !checkDirExist(SavePath) {
		if err := os.Mkdir(SavePath, os.ModeDir); err != nil {
			return fmt.Errorf("SavePath create fail: %s", err)
		}
	}
	if !CheckFileExist(defaultTargetPath) {
		return fmt.Errorf("target file not exist: %s", defaultTargetPath)
	}
	return nil
}
