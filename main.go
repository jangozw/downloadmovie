package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"
)

func main() {
	//url := "https://www.rebooo.com/archives/2467"
	flagInit := func() string {
		url := flag.String("url", "http://www.test.com", "下载页面的地址")
		flag.Parse()
		r:=regexp.MustCompile(`^http`)
		if ! r.MatchString(*url) {
			panic("请输入正确的url!")
		}
		return *url
	}
	url := flagInit()
	fmt.Println("您输入的页面地址：", url)

	ch := make(chan string)
	go getUrlContent(url, ch)
	content := loading(ch)
	reg := regexp.MustCompile(`<a href="(ed2k://.*?)"`)
	m := reg.FindAllStringSubmatch(content, -1)
	urls := make([]string, len(m))
	for index, value := range m {
		urls[index] = value[1]
	}
	urls = RemoveRepeatedElement(urls)
	fmt.Println("下载链接", len(urls), "个:")
	for _, url := range urls {
		fmt.Println(url)
	}
}

func getUrlContent(url string, ch chan string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("\n 获取资源失败: ", err)
	}
	defer resp.Body.Close()
	shtml, _ := ioutil.ReadAll(resp.Body)
	ch <- string(shtml)
}
func loading(ch chan string) string {
	fmt.Print("正在获取页面内容 ")
	var c int
	for {
		select {
		case content := <-ch:
			fmt.Println("\n即将解析下载链接")
			return content
		default:
			fmt.Print("#")
			if c > 0 && c % 100 == 0 {
				fmt.Print("\n")
			}
			time.Sleep(50 * time.Millisecond)
		}
		c++
		if c > 1000{
			panic("没有接收到数据, 异常")
		}
	}
}

func RemoveRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}


