package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/asciimoo/solary/arena"
	"github.com/asciimoo/solary/arena/board"
)

const (
	PLAYERS_NUM int = 2
)

func gameLoop(address string) {
	game_server, _ := net.Listen("tcp", address)
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

func main() {

	address := flag.String("listen", "127.0.0.1:6666", "server listen address")

	flag.Parse()
	if len(flag.Args()) != 0 {
		b := board.Create()
		fmt.Println(b)
	} else {
		gameLoop(*address)
	}

}
