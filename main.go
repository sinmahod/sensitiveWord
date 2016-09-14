package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var wordmap map[rune]interface{}
var count int

func init() {
	wordmap = make(map[rune]interface{})
	count = 0
	if f, err := os.Open("dic.txt"); err == nil {
		input := bufio.NewScanner(f)
		for input.Scan() {
			if input.Err() != nil {
				fmt.Println(input.Err())
				continue
			}
			addString(input.Text())
			count++
		}
	} else {
		fmt.Fprintln(os.Stdout, err)
	}
}

func form(w http.ResponseWriter, r *http.Request) {
	if data, err := ioutil.ReadFile("form.html"); err == nil {
		html := string(data)
		html = strings.Replace(html, "${DicTotal}", strconv.Itoa(count), 1)
		fmt.Fprintln(w, html)
	}
}

func addString(word string) {
	//声明一个map
	var wm map[rune]interface{} = wordmap
	//把得到的敏感词汇分割单个的unicode字符
	srune := []rune(word)
	//遍历
	for _, r := range srune {
		//去map中寻找有没有这个key
		if mp, ok := wm[r]; !ok { //如果没有则存储这个key
			mp = make(map[rune]interface{})
			wm[r] = mp
			wm = mp.(map[rune]interface{}) //用子map继续存储下一个字符
		} else {
			//如果存在则直接返回子map
			wm = mp.(map[rune]interface{})
		}
	}
}

/*
 * ## return ##
 *  []string	匹配到的敏感词
 *  string	替换后的内容
 *  int		敏感词命中次数
 */
func findWord(content string) ([]string, string, int) {
	count, strs, srune, wmp := 0, []string{}, []rune(content), make(map[string]int)
	for i := 0; i < len(srune); i++ {
		strs2, wm := []rune{}, wordmap
		for j := i; j < len(srune); j++ {
			if mp, ok := wm[srune[j]]; ok {
				strs2 = append(strs2, srune[j])
				if isnil(mp.(map[rune]interface{})) { //如果没有子元素
					for x := i; x <= j; x++ {
						//替换敏感词
						srune[x] = rune('*')
					}
					//判断这个词是否已经被匹配到了，如果被匹配到了则次数+1
					wmp[string(strs2)]++
					//匹配成功后从词汇之后继续匹配（去掉i=j则从词汇的第二个字继续匹配）
					i = j
				}
				wm = mp.(map[rune]interface{})
			} else {
				break
			}
		}
	}
	for dic, c := range wmp {
		count += c
		strs = append(strs, dic+"("+strconv.Itoa(c)+")")
	}
	return strs, string(srune), count
}

//判断map是否没有元素 没有true
func isnil(mp map[rune]interface{}) bool {
	for range mp {
		return false
	}
	return true
}

type data struct {
	Content string   `json:"content"`
	Count   int      `json:"count"`
	Dic     []string `json:"dic"`
	Time    string   `json:"time"`
}

func handlers(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	for k, v := range r.Form {
		//读取参数
		if k == `content` {
			if len(v) > 0 {
				s := time.Now()
				dics, content, count := findWord(v[0])
				e := time.Since(s).Nanoseconds()
				jsondata := &data{content, count, dics, fmt.Sprintf("%f", float64(e)/1000000)}
				js, _ := json.Marshal(jsondata)
				fmt.Fprintf(w, "%s", js)
			}
		}
	}
}

func main() {
	http.HandleFunc("/", form)
	http.HandleFunc("/sensitive", handlers)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
