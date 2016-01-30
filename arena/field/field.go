package field

import (
	"github.com/asciimoo/solary/arena/loot"
)

type Field struct {
	Loot  []string `json:",omitempty"`
	Traps int      `json:",omitempty"`
	Type  string
}

func Create(field_type string) *Field {
	return &Field{make([]string, 0), 0, field_type}
}

func (f *Field) ClearTraps() {
	f.Traps = 0
}

func (f *Field) AddTrap() {
	f.Traps += 1
}

func (f *Field) Passable() bool {
	return true
}

func (f *Field) AddRandomLoot() {
	f.Loot = append(f.Loot, loot.GetRandomLootName())
}

func (f *Field) ClearLoot() {
	f.Loot = make([]string, 0)
}
