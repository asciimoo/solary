package board

import (
	"math/rand"

	"github.com/asciimoo/solary/arena/coord"
	"github.com/asciimoo/solary/arena/field"
)

type Board struct {
	height uint
	width  uint
	Fields [][]*field.Field
}

func Create(height, width uint) *Board {
	b := &Board{
		height,
		width,
		make([][]*field.Field, height),
	}
	for y := uint(0); y < height; y++ {
		b.Fields[y] = make([]*field.Field, width)
		for x := uint(0); x < width; x++ {
			b.Fields[y][x] = field.Create()
		}
	}
	b.PopulateRandomLoot()
	return b
}

func (b *Board) PopulateRandomLoot() {
	for y := uint(0); y < b.height; y++ {
		for x := uint(0); x < b.width; x++ {
			if rand.Int()%10 != 0 {
				b.Fields[y][x].AddRandomLoot()
			}
		}
	}
}

func (b *Board) IsValidLocation(x, y uint) bool {
	return x < b.width && y < b.height
}

func (b *Board) FieldByCoord(c coord.Coord) *field.Field {
	return b.Fields[c.Y][c.X]
}
