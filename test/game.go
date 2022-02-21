package main

import (
	"fmt"

	"github.com/maaslalani/gambit/game"
)

func main() {
	gm := game.NewGame()
	b := gm.View()
	fmt.Println(b)
}
