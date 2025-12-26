package motion

type Result struct {
	Motion   Motion
	Count    int
	HasCount bool
	Done     bool
	Valid    bool
}

type Parser struct {
	root    *TrieNode
	current *TrieNode
	count   int
}

func NewParser(root *TrieNode) *Parser {
	return &Parser{
		root:    root,
		current: root,
	}
}

func (p *Parser) Reset() {
	p.current = p.root
	p.count = 0
}

func (p *Parser) ResetWith(newRoot *TrieNode) {
	p.current = newRoot
	p.count = 0
}

func (p *Parser) Feed(k Key) Result {
	if k.Code == KeyRune && k.Rune >= '0' && k.Rune <= '9' {
		p.count = p.count*10 + int(k.Rune-'0')
		return Result{Valid: true}
	}

	next := p.current.Children[k]
	if next == nil {
		p.Reset()
		return Result{Done: true, Valid: false}
	}

	p.current = next

	if next.Motion != nil {
		c := p.count
		hasCount := true
		if c == 0 {
			c = 1
			hasCount = false
		}
		m := next.Motion
		p.Reset()
		return Result{
			Motion:   m,
			Count:    c,
			HasCount: hasCount,
			Done:     true,
			Valid:    true,
		}
	}

	return Result{Valid: true}
}
