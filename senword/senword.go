package senword

import "fmt"

var Wordmap map[string]interface{}

func mian() {
	str := "毛泽东"
	add(str, Wordmap)
	fmt.Println(Wordmap)
}

func add(str string, wm map[string]interface{}) {
	rs := []rune(str)
	for i := 0; i < len(rs); i++ {
		s := string(rs[i])
		if _, ok := wm[s]; !ok {
			//如果不存在这个key
			mp := make(map[string]interface{})
			mp["isEnd"] = false
			wm[s] = mp
		} else {
			//如果存在这个key
			mp := wm[s]
			add(s, mp)
		}
	}
}

func addRune(s string, mp map[string]interface{}) {
	if _, ok := mp[s]; !ok {
		mp[s] = 0
	}
}
