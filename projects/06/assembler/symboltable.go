package main

type Address = uint16

type SymbolTable struct {
	table          map[string]Address
	nextRAMAddress Address
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		table: map[string]Address{
			"SP":     0x0000,
			"LCL":    0x0001,
			"ARG":    0x0002,
			"THIS":   0x0003,
			"THAT":   0x0004,
			"R0":     0x0000,
			"R1":     0x0001,
			"R2":     0x0002,
			"R3":     0x0003,
			"R4":     0x0004,
			"R5":     0x0005,
			"R6":     0x0006,
			"R7":     0x0007,
			"R8":     0x0008,
			"R9":     0x0009,
			"R10":    0x000a,
			"R11":    0x000b,
			"R12":    0x000c,
			"R13":    0x000d,
			"R14":    0x000e,
			"R15":    0x000f,
			"SCREEN": 0x4000,
			"KBD":    0x6000,
		},
		nextRAMAddress: 0x0010,
	}
}

func (s *SymbolTable) GetAddress(key string) (value Address, ok bool) {
	value, ok = s.table[key]
	return
}

func (s *SymbolTable) AddEntry(key string, value Address) {
	s.table[key] = value
}

// AddAutoEntry adds new entry to table and return its address.
// The address of new entry is auto incremented value.
func (s *SymbolTable) AddAutoEntry(key string) Address {
	value := s.nextRAMAddress
	s.table[key] = value
	s.nextRAMAddress++
	return value
}

func (s *SymbolTable) LoadLabelAddress(p *Parser) {
	var romAddr Address
	for p.Parse() {
		if cmd, ok := p.CurrentCommand().(*LCommand); ok {
			s.AddEntry(cmd.Symbol, romAddr)
		}
		romAddr++
	}
}
