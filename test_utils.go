package gmcts

import (
	"errors"
	"fmt"

	checkers "github.com/0xhexnumbers/go-checkers"
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

type cGame struct {
	game checkers.Game
}

func (g cGame) GetActions() []Action {
	gameActions := g.game.GetActions()

	actions := make([]Action, len(gameActions))

	for i, a := range gameActions {
		actions[i] = a
	}

	return actions
}

func (g cGame) ApplyAction(a Action) (Game, error) {
	action, ok := a.(checkers.Move)
	if !ok {
		return nil, errors.New("gmcts: could not convert action to checkers move")
	}

	game, err := g.game.ApplyAction(action)
	return cGame{game}, err
}

func (g cGame) Player() Player {
	return getPlayerID(g.game.Player())
}

func (g cGame) IsTerminal() bool {
	return g.game.IsTerminalState()
}

func (g cGame) Winners() []Player {
	winner, _ := g.game.Winner()

	if winner == '_' {
		return []Player{Player(0), Player(1)}
	}
	return []Player{getPlayerID(winner)}
}

//Read (*node).equals() to see why we have it here.
func (t *Tree) equals(t2 *Tree) bool {
	if len(t.gameStates) != len(t2.gameStates) {
		return false
	}

	for s1 := range t.gameStates {
		if _, ok := t2.gameStates[s1]; !ok {
			return false
		}
	}

	return t.current.equals(t2.current)
}

//We make our own equals method, because reflect.DeepEqual
//cannot satisfy our needs. Specifically, each node's parents
//may be in a different order from 1 node to another, so we
//need to handle that case specifically.
func (n1 *node) equals(n2 *node) bool {
	if len(n1.nodeScore) != len(n2.nodeScore) {
		return false
	}

	for player := range n1.nodeScore {
		if n2Score, ok := n2.nodeScore[player]; !ok {
			return false
		} else if n1.nodeScore[player] != n2Score {
			return false
		}
	}

	if n1.nodeVisits != n2.nodeVisits {
		return false
	}

	if len(n1.children) != len(n2.children) {
		return false
	}

	if len(n1.parents) != len(n2.parents) {
		return false
	}

	//Parents may be unordered, so we mappify n1's parents' states,
	//and check n2's parents' states against n1. Because n1 and n2
	//have the same number of parents, this will correctly
	//determine whether they have the same parents
	n1Parents := make(map[gameState]bool)
	for _, p1 := range n1.parents {
		n1Parents[p1.state] = true
	}
	for _, p2 := range n2.parents {
		if !n1Parents[p2.state] {
			return false
		}
	}

	for a, c1 := range n1.children {
		if c2, ok := n2.children[a]; !ok {
			return false
		} else if !c1.equals(c2) {
			return false
		}
	}

	return true
}
