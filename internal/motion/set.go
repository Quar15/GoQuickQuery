package motion

type Set struct {
	root *TrieNode
}

func NewSet() *Set {
	return &Set{
		root: NewTrie(),
	}
}

func (s *Set) Root() *TrieNode {
	return s.root
}

func (s *Set) Add(keys []Key, m Motion) {
	s.root.Insert(keys, m)
}

func (s *Set) AddRune(r rune, m Motion) {
	s.Add([]Key{{Code: KeyRune, Rune: r}}, m)
}

func (s *Set) AddArrow(r rune, m Motion) {
	s.Add([]Key{{Code: KeyArrow, Rune: r}}, m)
}

func (s *Set) AddCtrl(r rune, m Motion) {
	s.Add([]Key{{
		Code:      KeyRune,
		Rune:      r,
		Modifiers: ModCtrl,
	}}, m)
}
