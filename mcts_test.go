package gmcts

import (
	"fmt"
	"sync"
	"testing"

	tictactoe "github.com/0xhexnumbers/go-tic-tac-toe"
)

func getPlayerID(ascii byte) Player {
	if ascii == 'x' || ascii == 'X' {
		return Player(0)
	}
	return Player(1)
}

type tttGame struct {
	game tictactoe.Game
}

func (g tttGame) GetActions() []Action {
	gameActions := g.game.GetActions()

	actions := make([]Action, len(gameActions))

	for i, a := range gameActions {
		actions[i] = a
	}

	return actions
}

func (g tttGame) ApplyAction(a Action) (Game, error) {
	move, ok := a.(tictactoe.Move)
	if !ok {
		return nil, fmt.Errorf("action not correct type")
	}

	game, err := g.game.ApplyAction(move)

	return tttGame{game}, err
}

func (g tttGame) Player() Player {
	return getPlayerID(g.game.Player())
}

func (g tttGame) IsTerminal() bool {
	return g.game.IsTerminal()
}

func (g tttGame) Winners() []Player {
	winner, _ := g.game.Winner()
	if winner == '_' {
		return []Player{Player(0), Player(1)}
	}

	return []Player{getPlayerID(winner)}
}

//Global vars to be checked by other tests
var finishedGame tttGame
var firstMove tictactoe.Move
var treeToTest *Tree

//TestMain runs through a tictactoe game, saving the first move made and
//the resulting terminal game state into global variables to be used by
//other tests.
func TestMain(m *testing.M) {
	game := tttGame{tictactoe.NewGame()}
	concurrentSearches := 1 //runtime.NumCPU()

	var setFirstMove sync.Once
	var setTestingTree sync.Once

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

				//Set the tree to perform benchmarks on
				setTestingTree.Do(func() {
					treeToTest = tree
				})
			}()
		}
		wait.Wait()

		bestAction := mcts.BestAction()
		nextState, _ := game.ApplyAction(bestAction)
		game = nextState.(tttGame)
		fmt.Println(game.game)

		//Save the first action taken
		setFirstMove.Do(func() {
			firstMove = bestAction.(tictactoe.Move)
		})
	}
	//Save the terminal game state
	finishedGame = game

	m.Run()
}

func TestTicTacToeDraw(t *testing.T) {
	//Fail if there's a winner. Because tic-tac-toe is a simple game,
	//this algorithm should've finished in a draw.
	if len(finishedGame.Winners()) != 2 {
		t.Errorf("gmcts: tic-tac-toe game did not end in a draw")
		t.FailNow()
	}
}

func TestTicTacToeMiddle(t *testing.T) {
	//Fail if the first move doesn't pick the middle square. Because tic-tac-toe
	//is a simple game, this algorithm should've picked the middle square.
	if fmt.Sprintf("%v", firstMove) != "{1 1}" {
		t.Errorf("gmcts: first action is not to take the middle spot: %v", firstMove)
		t.FailNow()
	}
}

func TestZeroTrees(t *testing.T) {
	mcts := NewMCTS(finishedGame)
	bestAction := mcts.BestAction()
	if bestAction != nil {
		t.Errorf("gmcts: recieved a best action from no trees: %#v", bestAction)
		t.FailNow()
	}
}

func TestTerminalState(t *testing.T) {
	mcts := NewMCTS(finishedGame)
	mcts.AddTree(mcts.SpawnTree())
	bestAction := mcts.BestAction()
	if bestAction != nil {
		t.Errorf("gmcts: recieved a best action from a terminal state: %#v", bestAction)
		t.FailNow()
	}
}

func BenchmarkTicTacToe1KRounds(b *testing.B) {
	mcts := NewMCTS(tttGame{tictactoe.NewGame()})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mcts.SpawnTree().SearchRounds(1000)
	}
}

func BenchmarkTicTacToe10KRounds(b *testing.B) {
	mcts := NewMCTS(tttGame{tictactoe.NewGame()})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mcts.SpawnTree().SearchRounds(10000)
	}
}

func BenchmarkTicTacToe100KRounds(b *testing.B) {
	mcts := NewMCTS(tttGame{tictactoe.NewGame()})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mcts.SpawnTree().SearchRounds(100000)
	}
}
