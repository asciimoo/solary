package arena

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/asciimoo/solary/arena/board"
	"github.com/asciimoo/solary/arena/coord"
	"github.com/asciimoo/solary/player"
)

const ROUND_TIMEOUT = 2

var battle_id uint = 0
var player_id uint = 0

type Arena struct {
	Id      uint
	Round   uint
	Players []*player.Player
	Board   *board.Board
}

func Create() *Arena {
	a := &Arena{
		battle_id,
		0,
		make([]*player.Player, 0),
		board.Create(),
	}
	battle_id += 1
	return a
}

func (a *Arena) Play() {
	fmt.Println("Game", a.Id, "started")
	ch := make(chan *player.Move)
	for _, p := range a.Players {
		defer p.Conn.Close()
		go p.Read(ch)
	}
	for {
		a.Round += 1
		a.broadcastStatus()
		if a.Round == 500 || a.getActivePlayersNum() == 0 {
			fmt.Println("Game", a.Id, "finished")
			return
		}
		if a.Round%100 == 0 {
			a.Board.PopulateRandomLoot()
		}
		fmt.Println("Round", a.Round, "started")
		// collect moves
		moves := a.getMoves(ch)
		// activate laser beams and traps
		for _, move := range moves {
			if move.Item == "" {
				continue
			}
			c, ok := move.Player.Inventory[move.Item]
			if !ok || c <= 0 {
				continue
			}
			move.Player.Inventory[move.Item] -= 1
			switch move.Item {
			case "trap":
				a.Board.FieldByCoord(move.Player.Position).AddTrap()
			case "laser beam":
				laser_x := 0
				laser_y := 0
				switch move.Direction {
				case "up":
					laser_y -= 1
				case "down":
					laser_y += 1
				case "left":
					laser_x -= 1
				case "right":
					laser_x += 1
				}
				for i := 1; i <= 2; i++ {
					laser_coord := coord.Coord{
						uint(laser_x*i) + move.Player.Position.X,
						uint(laser_y*i) + move.Player.Position.Y,
					}
					for _, p := range a.Players {
						if p.Position.X == laser_coord.X && p.Position.Y == laser_coord.Y {
							p.Life -= 25
						}
					}
				}
			case "oil":
				if move.Player.Life > 80 {
					move.Player.Life = 100
				} else {
					move.Player.Life += 20
				}
			}
		}

		a.movePlayers(&moves)

		a.checkDeath()

		// helper coord->user map
		player_coords := make(map[coord.Coord][]*player.Player)
		for _, p := range a.Players {
			player_coords[p.Position] = append(player_coords[p.Position], p)
		}

		// activate destination traps && collect loot
		for c, players := range player_coords {
			field := a.Board.FieldByCoord(c)
			if len(players) == 1 {
				for _, l := range a.Board.FieldByCoord(c).Loot {
					players[0].Inventory[l] += 1
				}
				field.ClearLoot()
			}
			if field.Traps > 0 {
				for _, p := range players {
					p.Life -= 50 * a.Board.FieldByCoord(c).Traps
				}
				field.ClearTraps()
			}
		}

		a.checkDeath()

		// trigger inventory items
		for _, p := range a.Players {
			p.Score += p.Inventory["solar panel"]
		}
	}
}

func (a *Arena) movePlayers(moves *[]*player.Move) {
	for _, move := range *moves {
		a.movePlayer(move)
	}
}

func (a *Arena) checkDeath() {
	for _, p := range a.Players {
		if p.Life <= 0 {
			field := a.Board.FieldByCoord(p.Position)
			for loot, loot_count := range p.Inventory {
				p.Inventory[loot] = 0
				for i := 0; i < loot_count; i++ {
					field.Loot = append(field.Loot, loot)
				}
			}
			p.Life = 100
			p.Position.X = 0
			p.Position.Y = 0
		}
	}
}

func (a *Arena) movePlayer(m *player.Move) {
	new_y := m.Player.Position.Y
	new_x := m.Player.Position.X
	distance := uint(1)
	if m.Item == "pogo stick" {
		distance = 2
	}
	switch m.Direction {
	case "up":
		new_y -= distance
	case "down":
		new_y += distance
	case "left":
		new_x -= distance
	case "right":
		new_x += distance
	}
	if a.Board.IsValidLocation(int(new_x), int(new_y)) {
		m.Player.Position.X = new_x
		m.Player.Position.Y = new_y
	}
	fmt.Println("player", m.Player.Id, "moves to", m.Direction)
}

func (a *Arena) getMoves(ch chan *player.Move) []*player.Move {
	moves := []*player.Move{}
	timeout := make(chan bool)
	go func() {
		time.Sleep(ROUND_TIMEOUT * time.Second)
		timeout <- true
	}()
	recv_break := false
	for !recv_break && len(moves) < a.getActivePlayersNum() {
		select {
		case move := <-ch:
			if move.Error != nil {
				fmt.Println("user error:", move.Error)
				break
			}
			can_move := true
			for _, m := range moves {
				if m.Player == move.Player {
					can_move = false
					break
				}
			}
			if can_move {
				moves = append(moves, move)
			} else {
				fmt.Println("Error: already moved")
			}
		case <-timeout:
			fmt.Println("Timeout")
			recv_break = true
		}
	}
	return moves
}

func (a *Arena) getActivePlayersNum() int {
	active_players := 0
	for _, player := range a.Players {
		if !player.Disconnected {
			active_players += 1
		}
	}
	return active_players
}

func (a *Arena) broadcastStatus() {
	b, err := json.Marshal(a)
	if err != nil {
		fmt.Println(err)
	}
	for _, p := range a.Players {
		p.Write(b)
	}
}

func (a *Arena) AddPlayer(conn io.ReadWriteCloser) {
	p := player.Create(player_id, conn)
	msg, _ := json.Marshal(p)
	p.Write(msg)
	player_id += 1
	a.Players = append(a.Players, p)
}
