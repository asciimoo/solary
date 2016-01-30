package main

import (
	"fmt"
	"net"

	"github.com/asciimoo/solary/arena"
)

const (
	PLAYERS_NUM int = 2
)

func main() {

	game_server, _ := net.Listen("tcp", ":6666")
	fmt.Println("Server started")

	for {
		a := arena.Create()
		for i := 0; i < PLAYERS_NUM; i++ {
			conn, _ := game_server.Accept()
			a.AddPlayer(conn)
		}
		fmt.Println("Game", a.Id, "started")
		go a.Play()
	}
}
