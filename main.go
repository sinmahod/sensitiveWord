package main

import "fmt"

func init() {
	Wordmap = make(map[rune]interface{})
	str := "毛泽东"
	str2 := "毛ze东1"
	str3 := "ze东"
	addString(str)
	addString(str2)
	addString(str3)
}

var Wordmap map[rune]interface{}

func addString(word string) {
	//声明一个map
	var wm map[rune]interface{} = Wordmap
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

func findWord(content string) []string {
	var wm map[rune]interface{}
	srune := []rune(content)
	strs := []string{}
	for i := 0; i < len(srune); i++ {
		wm = Wordmap
		strs2 := []rune{}
		for j := i; j < len(srune); j++ {
			if mp, ok := wm[srune[j]]; ok {
				strs2 = append(strs2, srune[j])
				if isnil(mp.(map[rune]interface{})) { //如果没有子元素
					strs = append(strs, string(strs2))
				}
				wm = mp.(map[rune]interface{})
			} else {
				break
			}
		}
	}
	return strs
}

//判断map是否没有元素 没有true
func isnil(mp map[rune]interface{}) bool {
	for range mp {
		return false
	}
	return true
}

func main() {
	content := "中国最伟大的领袖是毛泽东，就是毛ze东。"
	s := findWord(content)
	fmt.Println(s)
}
