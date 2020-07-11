package board

import (
	"testing"
	"fmt"
)

func printPlayers(b *Board) {
	fmt.Println("---");
	for i := 1; i<=len(b.Seats); i++ {
		fmt.Printf("Player %d: %v\n", i, b.Player(i))
	}
}

func TestNewBoard(t *testing.T) {
	b := new(Board)
	meta := map[string]string {"version": "0.1",}
	roles := []string{"预言家","女巫","猎人","村民","高级村民","白狼王","狼人","机械狼","混血儿"}
	b.New(123, roles, meta)

}