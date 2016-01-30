package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/asciimoo/solary/arena"
	"github.com/asciimoo/solary/player"
)

var Directions = []string{"up", "down", "left", "right"}

type Client struct {
	Me     *player.Player
	Arena  *arena.Arena
	Reader *bufio.Reader
	Writer *bufio.Writer
}

func CreateClient(address string) (*Client, bool) {
	var me *player.Player

	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	b, err := reader.ReadBytes('\n')
	if err != nil {
		fmt.Println("disconnected")
		return &Client{}, false
	}

	err = json.Unmarshal(b, &me)
	if err != nil {
		fmt.Println(string(b))
		fmt.Println(err)
		fmt.Println("cannot deecode init json")
	} else {
		fmt.Println("my id:", me.Id)
	}

	return &Client{
		me,
		&arena.Arena{},
		reader,
		writer,
	}, true

}

func (c *Client) ReceiveNextRound() bool {
	var a *arena.Arena
	b, err := c.Reader.ReadBytes('\n')
	if err != nil {
		fmt.Println("disconnected")
		return false
	}
	err = json.Unmarshal(b, &a)
	if err != nil {
		fmt.Println("cannot deecode json")
		return false
	}
	c.Arena = a
	return true
}

func (c *Client) Write(b []byte) {
	c.Writer.Write(b)
	c.Writer.WriteRune('\n')
	c.Writer.Flush()
}

func PrintBoard(c *Client) {
	for _, p := range c.Arena.Players {
		if p.Id == c.Me.Id {
			fmt.Println(p.Life, p.Score, p.Inventory, p.Position)
		}
	}
	for y, row := range c.Arena.Board.Fields {
		row1 := ""
		for x, field := range row {
			if field.Type == "rock" {
				row1 += " \033[41mXX\033[0m"
				continue
			}
			players_num := 0
			player_color := 0
			for _, p := range c.Arena.Players {
				if p.Position.X == uint(x) && p.Position.Y == uint(y) {
					if p.Id == c.Me.Id {
						player_color = 42
					} else {
						player_color = 43 + (int(p.Id) % len(c.Arena.Players))
					}
					players_num += 1
				} else if p.SpawnPosition.X == uint(x) && p.SpawnPosition.Y == uint(y) {
					player_color = 44
				}
			}
			if player_color > 0 {
				row1 += fmt.Sprintf(" \033[%vm%v%v\033[0m", player_color, len(field.Loot), field.Traps)
			} else {
				row1 += fmt.Sprintf(" %v%v", len(field.Loot), field.Traps)
			}

		}
		fmt.Println(row1)
		//fmt.Println(row2)
	}
	time.Sleep(100 * time.Millisecond)
}
