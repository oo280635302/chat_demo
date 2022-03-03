package tool

type TrieNode struct {
	Node map[rune]*TrieNode
	End  bool
}

func NewTrieNode() *TrieNode {
	return &TrieNode{Node: map[rune]*TrieNode{}}
}

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

func (t *TrieNode) Replace(str string) string {
	n := len(str)
	s := []rune(str)
	ans := s

	for idx, v := range str {
		// 首字母匹配
		if _, ok := t.Node[v]; !ok {
			continue
		}

		// 匹配成功
		node := t
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

		if end != -1 {
			for i := idx; i < end+1; i++ {
				ans[i] = '*'
			}
		}

	}

	return string(ans)
}
