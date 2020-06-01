package gmcts

import (
	"fmt"
	"sync"
	"testing"

	tictactoe "github.com/0xhexnumbers/go-tic-tac-toe"
)

func TestTicTacToeDraw(t *testing.T) {
	game := tttGame{tictactoe.NewGame()}
	concurrentSearches := 1 //runtime.NumCPU()

	for !game.IsTerminal() {
		mcts := NewMCTS(game)

		var wait sync.WaitGroup
		wait.Add(concurrentSearches)
		for i := 0; i < concurrentSearches; i++ {
			go func() {
				tree := mcts.SpawnTree()
				tree.SearchRounds(10000)
				mcts.AddTree(tree)
				wait.Done()
			}()
		}
		wait.Wait()

		bestAction := mcts.BestAction()
		_, ok := bestAction.(tictactoe.Move)
		if !ok {
			t.Errorf("gmcts: type of best action is not a move: %T", bestAction)
			t.FailNow()
		} else {
			nextState, _ := game.ApplyAction(bestAction)
			game = nextState.(tttGame)
			fmt.Println(game.game)
		}
	}

	//Fail if there's a winner. Because tic-tac-toe is a simple game,
	//this algorithm should've finished in a draw.
	if len(game.Winners()) != 2 {
		t.Errorf("gmcts: tic-tac-toe game did not end in a draw")
		t.FailNow()
	}
}

func TestTicTacToeMiddle(t *testing.T) {
	mcts := NewMCTS(tttGame{tictactoe.NewGame()})
	concurrentSearches := 1 //runtime.NumCPU()

	var wait sync.WaitGroup
	wait.Add(concurrentSearches)
	for i := 0; i < concurrentSearches; i++ {
		go func() {
			tree := mcts.SpawnTree()
			tree.SearchRounds(10000)
			mcts.AddTree(tree)
			wait.Done()
		}()
	}
	wait.Wait()

	bestAction := mcts.BestAction()
	action, ok := bestAction.(tictactoe.Move)
	if !ok {
		t.Errorf("gmcts: type of best action is not a move: %T", bestAction)
		t.FailNow()
	} else {
		if fmt.Sprintf("%v", action) != "{1 1}" {
			t.Errorf("gmcts: first action is not to take the middle spot: %v", action)
			t.FailNow()
		}
	}
}
