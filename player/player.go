package player

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"

	"github.com/asciimoo/solary/arena/coord"
	"github.com/asciimoo/solary/arena/loot"
)

const (
	MAX_LIFE int = 100
)

type Player struct {
	Id           uint
	Life         int
	Position     coord.Coord
	Score        int
	Inventory    map[string]int
	Disconnected bool `json:"-"`
	reader       *bufio.Reader
	writer       *bufio.Writer
}

func Create(id uint, conn io.ReadWriter) *Player {
	return &Player{
		id,
		MAX_LIFE,
		coord.Coord{},
		0,
		make(map[string]int),
		false,
		bufio.NewReader(conn),
		bufio.NewWriter(conn),
	}
}

type Move struct {
	Player    *Player `json:",omitempty"`
	Direction string
	Item      string `json:",omitempty"`
	Error     error  `json:"-"`
}

func (p *Player) Read(c chan *Move) {
	for {
		var m Move
		m.Player = p
		//s, err := p.reader.ReadString('\n')
		b, err := p.reader.ReadBytes('\n')
		if err != nil {
			p.Disconnected = true
			m.Error = err
		}
		err = json.Unmarshal(b, &m)
		if err != nil {
			m.Error = err
		}
		c <- &m
		if p.Disconnected {
			fmt.Println("player", p.Id, "disconnected")
			return
		}
	}
}

func (p *Player) Write(x []byte) {
	if p.Disconnected {
		return
	}
	_, err := p.writer.Write(x)
	if err != nil {
		p.Disconnected = true
		return
	}
	_, err = p.writer.WriteRune('\n')
	if err != nil {
		p.Disconnected = true
		return
	}
	p.writer.Flush()
}

func (p *Player) AddLoot(l loot.Loot) {
	p.Inventory[l.String()] += 1
}
