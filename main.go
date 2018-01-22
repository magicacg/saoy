package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/headzoo/surf"
)

func regexp0(s1, s2 string) []string {
	r1 := regexp.MustCompile(s2)
	r2 := r1.FindAllString(s1, -1)
	return r2
}
func regexp1(s1, s2 string) string {
	r1 := regexp.MustCompile(s2)
	r2 := r1.ReplaceAllString(s1, "")
	return r2
}
func Txtopen(s1, s2 string) {
	f0, _ := os.Create(s1 + ".txt") //写入到date1时间的TXT文件里
	f0.WriteString(s2)              //写入内容
	defer f0.Close()                //关闭文件

}

func main() {
	var url1 string
	var html, html2 string
	var pagego, pageend string
	var page string
	var k, l int
	bow := surf.NewBrowser()
	bow.AddRequestHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36")
	fi, err := os.Open("setting.txt") //读取文件夹里的setting.txt
	if err != nil {
		log.Println("Without this file:	setting.txt!")
		os.Exit(0)

	}
	defer fi.Close()              //关闭文件
	fd, err := ioutil.ReadAll(fi) //IO转换

	pagego = regexp1(regexp0(string(fd), `pagego=(\w+|)`)[0], `\w+=`)   //找到setting.txt里的pagego	为开始页数
	pageend = regexp1(regexp0(string(fd), `pageend=(\w+|)`)[0], `\w+=`) //找到setting.txt里的pageend 为结束页数

	if pagego == "" { //判断，如果pagego没有填入，则默认从第一页开始
		url1 = "https://yande.re/post?page="
		page = "1"

	} else { //如果pagego有值，则从值的页数开始
		url1 = "https://yande.re/post?page=" + pagego
		page = pagego
	}
	l, _ = strconv.Atoi(page) //把当前页数找出来赋值到l身上
	if pageend == "" {        //判断，如果pageend没有填入，则默认为当前页
		k = 1
	} else { //判断，如果有，则从为结束页
		k, _ = strconv.Atoi(pageend)
	}
	fmt.Println(l, k)

	for j := l; j < k+1; j++ { //进行循环，l为开始页数，K为结束页数
		for {
			err = bow.Open(url1) //打开网页
			if err != nil {
				time.Sleep(2 * time.Second)
				continue
			} else {
				break
			}
		}
		body := regexp0(bow.Body(), `class="thumb" href="/post/show/\w+`) //找到当前页数的所有帖子ID
		t := len(body)                                                    //找到的帖子ID
		for i := 0; i < t; i++ {                                          //打开当前页数的所有帖子
			for {
				err = bow.Open("https://yande.re/" + regexp1(body[i], `class="thumb" href="/`)) //打开找到的帖子ID
				if err != nil {                                                                 //如果错误则输出一个! ，并且5秒后跳转到LA2继续重新打开
					fmt.Printf("!")
					time.Sleep(2 * time.Second)
					continue
				} else {
					break
				}
			}
			body1 := regexp1(regexp0(bow.Body(), `class="(original-file-changed|original-file-unchanged)" id="highres" href=".+">`)[0], `(.+href=|"|>|<)`) //找到帖子里的图片链接，可能不齐全，可能会出现错误

			fmt.Printf("*")      //打印一个*代表成功
			html += body1 + "\n" //把所有链接赋值
			html2 += body1 + "\n"
		}
		fmt.Printf("\n")
		Txtopen("page"+page, html) //执行把html的内容保存到txt文件里

		html = "" //清空html的内容

		page = strconv.Itoa(j + 1)                  //整形转换字符串						//执行j+1操作，如果循环不结束，则当前页数+1页
		url1 = "https://yande.re/post?page=" + page //进行赋值，把下一页的链接赋值到url1里,执行下一次循环后会使用本url
	}
	Txtopen("setting", html2) //执行把html的内容保存到txt文件里

}
