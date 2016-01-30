package loot

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type Loot interface {
	String() string
	// Trigger(name string, move *player.Move, arena *arena.Arena) //User steps on it
}

var Loots = []string{"pogo stick", "laser beam", "trap", "solar panel", "oil"}

func GetRandomLootName() string {
	return Loots[rand.Intn(len(Loots))]
}
