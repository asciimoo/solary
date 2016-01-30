package board

import (
	"math/rand"
	"time"

	"github.com/asciimoo/solary/arena/coord"
	"github.com/asciimoo/solary/arena/field"
)

const (
	BOARD_SIZE_MIN = 10
	BOARD_SIZE_MAX = 30
)

type Board struct {
	height int
	width  int
	Fields [][]*field.Field
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func Create() *Board {
	size := (rand.Intn(BOARD_SIZE_MAX/2-BOARD_SIZE_MIN/2)+BOARD_SIZE_MIN/2)*2 + 1
	b := &Board{
		size,
		size,
		make([][]*field.Field, size),
	}
	for y := 0; y < size; y++ {
		b.Fields[y] = make([]*field.Field, size)
		for x := 0; x < size; x++ {
			b.Fields[y][x] = field.Create("rock")
		}
	}
	pos_x := rand.Intn(b.width)
	pos_y := rand.Intn(b.height)
	carve_count := (b.width * b.height) / 6
	for carve_count > 0 {
		b.mirrorCarve(pos_x, pos_y)
		if rand.Intn(2)%2 == 0 {
			pos_x = (pos_x + 1) % (b.width)
		} else {
			pos_y = (pos_y + 1) % (b.height)
		}
		carve_count -= 1
	}
	b.PopulateRandomLoot()
	return b
}

func (b *Board) mirrorCarve(x, y int) {
	b.Fields[y][x].Type = "ground"
	b.Fields[b.height-y-1][x].Type = "ground"
	b.Fields[y][b.width-x-1].Type = "ground"
	b.Fields[b.height-y-1][b.width-x-1].Type = "ground"
}

func (b *Board) String() string {
	board_str := ""
	for _, row := range b.Fields {
		for _, field := range row {
			if field.Type == "rock" {
				board_str += " X"
			} else {
				board_str += " ."
			}
		}
		board_str += "\n"
	}
	return board_str
}

func (b *Board) PopulateRandomLoot() {
	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			if b.Fields[y][x].Type != "rock" && rand.Int()%2 != 0 {
				b.Fields[y][x].AddRandomLoot()
			}
		}
	}
}

func (b *Board) IsValidLocation(x, y int) bool {
	return x >= 0 && y >= 0 && x < b.width && y < b.height && b.Fields[y][x].Type != "rock"
}

func (b *Board) FieldByCoord(c coord.Coord) *field.Field {
	return b.Fields[c.Y][c.X]
}
