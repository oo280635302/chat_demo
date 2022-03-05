package common

import "strings"

type TrieNode struct {
	Node map[rune]*TrieNode
	End  bool
}

var Trie *TrieNode

func InitTrieNode() {
	t := NewTrieNode()
	Libs := strings.Split(TRIELIBRARY, "\n")
	for _, v := range Libs {
		t.Insert(v)
	}
	Trie = t
}

func NewTrieNode() *TrieNode {
	return &TrieNode{Node: map[rune]*TrieNode{}}
}

// 填充树
func (t *TrieNode) Insert(str string) {
	if len(str) == 0 {
		return
	}

	cur := t
	for _, v := range str {
		if _, ok := cur.Node[v]; !ok {
			cur.Node[v] = NewTrieNode()
		}
		cur = cur.Node[v]
	}
	cur.End = true
}

// 替换函数
// TODO 优化建议：1.跳过已被替换的符号，2.是否有办法跳过第一次首字母的重复匹配
func (t *TrieNode) Replace(str string) string {
	n := len(str)
	s := []rune(str)
	ans := s

	for idx, v := range str {
		// 首字母匹配
		if _, ok := t.Node[v]; !ok {
			continue
		}

		// 从当前字母开始匹配
		node := t
		// 匹配最长情况使用 end 来保存最后一次结束节点
		j, end := idx, -1
		for j < n {
			// 找不到暂停
			if _, ok := node.Node[s[j]]; !ok {
				break
			}
			node = node.Node[s[j]]
			if node.End {
				end = j
			}
			j++
		}

		// end不是-1说明匹配到了 进行替换
		if end != -1 {
			for i := idx; i < end+1; i++ {
				ans[i] = '*'
			}
		}

	}

	return string(ans)
}
