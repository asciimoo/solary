package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/asciimoo/solary/arena"
	"github.com/asciimoo/solary/arena/board"
)

func gameLoop(address string, players_num int) {
	game_server, _ := net.Listen("tcp", address)
	fmt.Println("Server started on " + address)

	for {
		a := arena.Create()
		for i := 0; i < players_num; i++ {
			conn, _ := game_server.Accept()
			a.AddPlayer(conn)
		}
		fmt.Println("Game", a.Id, "started")
		go a.Play()
	}
}

func main() {

	address := flag.String("listen", "127.0.0.1:6666", "server listen address")
	players_per_game := flag.Int("players", 2, "number of players in one arena")

	flag.Parse()
	if len(flag.Args()) != 0 {
		b := board.Create()
		fmt.Println(b)
	} else {
		gameLoop(*address, *players_per_game)
	}

}
