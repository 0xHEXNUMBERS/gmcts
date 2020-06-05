package gmcts

import (
	"fmt"

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
