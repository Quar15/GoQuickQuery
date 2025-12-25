package motion

type TrieNode struct {
	Children map[Key]*TrieNode
	Motion   Motion
}

func NewTrie() *TrieNode {
	return &TrieNode{
		Children: make(map[Key]*TrieNode),
	}
}

func (t *TrieNode) Insert(keys []Key, m Motion) {
	node := t
	for _, k := range keys {
		if node.Children[k] == nil {
			node.Children[k] = NewTrie()
		}
		node = node.Children[k]
	}
	node.Motion = m
}
