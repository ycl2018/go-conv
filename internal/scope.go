package internal

import "strconv"

type Scope struct {
	Name           string
	Symbols        map[string]*Symbol
	EnclosingScope *Scope
}

func NewScope(name string, enclosingScope *Scope) *Scope {
	return &Scope{
		Name:           name,
		Symbols:        make(map[string]*Symbol),
		EnclosingScope: enclosingScope,
	}
}

func (c *Scope) Define(name string, s *Symbol) {
	c.Symbols[name] = s
	s.scope = c
}

func (c *Scope) Remove(name string) {
	delete(c.Symbols, name)
}

func (c *Scope) NextSymbol(prefix string) *Symbol {
	for i := 0; ; i++ {
		var name = prefix
		if i != 0 {
			name += strconv.Itoa(i)
		}
		_, ok := c.Resolve(name)
		if !ok {
			ret := NewSymbol(name)
			c.Define(name, ret)
			return ret
		}
	}
}

func (c *Scope) NextPair(kPrefix, vPrefix string) (k, v *Symbol) {
	for i := 0; ; i++ {
		var kName, vName = kPrefix, vPrefix
		if i != 0 {
			kName += strconv.Itoa(i)
			vName += strconv.Itoa(i)
		}
		_, ok1 := c.Resolve(kName)
		_, ok2 := c.Resolve(vName)
		if !ok1 && !ok2 {
			k, v = NewSymbol(kName), NewSymbol(vName)
			c.Define(kName, k)
			c.Define(vName, v)
			return k, v
		}
	}
}

func (c *Scope) Resolve(name string) (*Symbol, bool) {
	s, ok := c.Symbols[name]
	if ok {
		return s, ok
	}
	if c.EnclosingScope != nil {
		return c.EnclosingScope.Resolve(name)
	}
	return nil, false
}

type Symbol struct {
	Name  string
	scope *Scope
}

func NewSymbol(name string) *Symbol {
	return &Symbol{
		Name: name,
	}
}
