package main

type Scope int

const (
	ScopeStatic Scope = iota
	ScopeField
	ScopeArg
	ScopeVar
)

type SymbolTableEntry struct {
	Name  string
	Type  string
	Scope Scope
	Index int
}

type SymbolTable struct {
	classScope map[string]*SymbolTableEntry
	localScope map[string]*SymbolTableEntry
	nextIndex  map[Scope]int
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		classScope: make(map[string]*SymbolTableEntry),
		localScope: make(map[string]*SymbolTableEntry),
	}
}

func (s *SymbolTable) issueIndex(scope Scope) int {
	idx := s.nextIndex[scope]
	s.nextIndex[scope]++
	return idx
}

func (s *SymbolTable) getTableForScope(scope Scope) map[string]*SymbolTableEntry {
	switch scope {
	case ScopeStatic, ScopeField:
		return s.classScope
	case ScopeArg, ScopeVar:
		return s.localScope
	default:
		return nil /* unreachable */
	}
}

func (s *SymbolTable) Define(name, typ string, scope Scope) {
	table := s.getTableForScope(scope)
	table[name] = &SymbolTableEntry{
		Name:  name,
		Type:  typ,
		Scope: scope,
		Index: s.issueIndex(scope),
	}
}

func (s *SymbolTable) Count(scope Scope) int {
	var count int
	for _, e := range s.getTableForScope(scope) {
		if e.Scope == scope {
			count++
		}
	}
	return count
}

func (s *SymbolTable) Find(name string) *SymbolTableEntry {
	if e, ok := s.localScope[name]; ok {
		return e
	}
	if e, ok := s.classScope[name]; ok {
		return e
	}
	return nil
}

func (s *SymbolTable) ResetLocalScope() {
	clear(s.localScope)
}
