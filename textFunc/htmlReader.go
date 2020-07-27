package textFunc

import (
	"fmt"
	"regexp"
	"strings"
)

//配置变量
var (
	regexp_first  = `<a href="/song\?id=\d+"><b title="[^"]+">`        //用于找到音乐信息字符串的正则表达式
	regexp_second = `<a href="/song\?id=(\d+)"><b title="([^"]+)">`    //用于对音乐信息字符串进行子匹配的正则表达式
	url_template  = `https://music.163.com/song/media/outer/url?id=%s` //播放音乐的链接模板，唯一的参数为音乐的id
)

//与音乐列表中的有用信息
type musicList struct {
	Name string
	Url  string
}

//从html文本中提取出音乐列表
func GetMusicList(html string) (err error, Resultlist []musicList) {
	if len(html) == 0 {
		return fmt.Errorf("nothing in html"), nil
	}

	//粗提取
	reg1, err := regexp.Compile(regexp_first)
	if err != nil {
		return fmt.Errorf("reg1 worng regexp: %v", err), nil
	}
	tmpList := reg1.FindAllString(html, 900)
	if len(tmpList) == 0 {
		return fmt.Errorf("can't find music in it text"), nil
	}
	fmt.Printf("The numbers of music found is %d \n", len(tmpList))

	//从上次提取的结果中分离出有用信息
	reg2, err := regexp.Compile(regexp_second)
	if err != nil {
		return fmt.Errorf("reg2 worng regexp: %v", err), nil
	}
	Resultlist = make([]musicList, len(tmpList))
	index := 0
	for _, tmpStr := range tmpList {
		params := reg2.FindStringSubmatch(tmpStr)
		if len(params) != 3 {
			fmt.Printf("This line maybe worng: %s \n", tmpStr)
			continue
		}
		Resultlist[index].Url = fmt.Sprintf(url_template, params[1])
		Resultlist[index].Name = changeName(params[2])
		index++
	}
	return nil, Resultlist
}

//纠正非法的音乐文件名
func changeName(name string) string {
	name = strings.ReplaceAll(name, "&nbsp;", " ")
	name = strings.ReplaceAll(name, "&amp;", "&")
	name = strings.ReplaceAll(name, "\\", "-")
	name = strings.ReplaceAll(name, "/", "-")
	return name
}
