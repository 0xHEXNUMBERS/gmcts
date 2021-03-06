package gmcts

import (
	"testing"
	"time"

	"github.com/0xhexnumbers/go-tic-tac-toe"
)

func TestRounds(t *testing.T) {
	rounds := treeToTest.Rounds()
	if rounds != 10000 {
		t.Errorf("Tree performed %d rounds: wanted 1", rounds)
		t.FailNow()
	}
}

func TestNodes(t *testing.T) {
	//The amount of nodes in the tree should not exceed the
	//amount of mcts rounds performed on the tree.
	rounds := treeToTest.Rounds()
	nodes := treeToTest.Nodes()
	if nodes > rounds {
		t.Errorf("Tree has %d nodes: wanted <= %d", nodes, rounds)
		t.FailNow()
	}
}

func TestDepth(t *testing.T) {
	//Because tictactoe is a simple game, the
	//tree should have looked 9 moves ahead.
	depth := treeToTest.MaxDepth()
	if depth != 9 {
		t.Errorf("Tree has depth %d: wanted 0", depth)
		t.FailNow()
	}
}

func TestSearch(t *testing.T) {
	mcts := NewMCTS(tttGame{tictactoe.NewGame()})
	tree := mcts.SpawnTree()

	timeToSearch := 1 * time.Millisecond
	t0 := time.Now()
	tree.Search(timeToSearch)
	td := time.Now().Sub(t0)

	if td < timeToSearch {
		t.Errorf("Tree was searched for %s: wanted >= %s", td, timeToSearch)
		t.FailNow()
	}
}
