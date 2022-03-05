package common

import (
	"strings"
	"testing"
)

func Test_Trie(t *testing.T) {
	tn := NewTrieNode()
	tn.Insert("shit")
	tn.Insert("fuck")
	tn.Insert("fucker")

	s := tn.Replace("oh shit fuck!", '*')
	t.Log(s)

	s = tn.Replace("oh shit fucker!", '*')
	t.Log(s)

	s = tn.Replace("oh shit fucker fucccker shit fuckee fuc s shit shiit  shhit !", '*')
	t.Log(s)
}

// 16核 15000ns/op
func Benchmark_Trie(b *testing.B) {
	tn := NewTrieNode()
	Libs := strings.Split(TRIELIBRARY, "\n")
	for _, v := range Libs {
		tn.Insert(v)
	}
	for i := 0; i < b.N; i++ {
		tn.Replace("oh shit fucker fucccker shit fuckee fuc s shit shiit  shhit !", '*')
	}
}
