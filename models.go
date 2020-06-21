package gmcts

import (
	"sync"
)

//Action is the interface that represents an action that can be
//performed on a Game.
//
//Any implementation of Action should be comparable (i.e. be a key in a map)
type Action interface{}

//Player is an id for the player
type Player int

//Game is the interface that represents game states.
//
//Any implementation of Game should be comparable (i.e. be a key in a map)
//and immutable (state cannot change as this package calls any function).
type Game interface {
	//GetActions returns a list of actions to consider
	GetActions() []Action

	//ApplyAction applies the given action to the game state,
	//and returns a new game state and an error for invalid actions
	ApplyAction(Action) (Game, error)

	//Player returns the player that can take the next action
	Player() Player

	//IsTerminal returns true if this game state is a terminal state
	IsTerminal() bool

	//Winners returns a list of players that have won the game if
	//IsTerminal() returns true
	Winners() []Player
}

type gameState struct {
	Game

	//This is to separate states that seemingly look the same,
	//but actually occur on different turn orders. Without this,
	//the directed tree that multiple parent nodes will just
	//become a directed graph, which this MCTS implementation
	//cannot handle properly.
	turn int
}

//MCTS contains functionality for the MCTS algorithm
type MCTS struct {
	init  Game
	trees []*Tree
	mutex *sync.RWMutex
}

type actionNodePair struct {
	action Action
	node   *node
}

type node struct {
	state gameState
	tree  *Tree

	parents           []*node
	children          []actionNodePair
	unvisitedChildren []actionNodePair

	nodeScore  map[Player]float64
	nodeVisits int
}

//Tree represents a game state tree
type Tree struct {
	current          *node
	gameStates       map[gameState]*node
	explorationConst float64
}
