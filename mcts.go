package gmcts

import (
	"math/rand"
	"reflect"
	"sync"
)

//NewMCTS returns a new MCTS wrapper
//
//If either the Game or its Action types are not comparable,
//this function panics
func NewMCTS(initial Game) *MCTS {
	//Check if Game type if comparable
	if !reflect.TypeOf(initial).Comparable() {
		panic("gmcts: game type is not comparable")
	}

	//Check if Action type is comparable
	//We only need to check the actions that can affect the initial gamestate
	//as those are the only actions that need to be compared.
	actions := initial.GetActions()
	for i := range actions {
		if !reflect.TypeOf(actions[i]).Comparable() {
			panic("gmcts: action type is not comparable")
		}
	}

	return &MCTS{
		init:  initial,
		trees: make([]*Tree, 0),
		mutex: new(sync.RWMutex),
	}
}

//SpawnTree creates a new search tree. The tree returned uses Sqrt(2) as the
//exploration constant.
func (m *MCTS) SpawnTree() *Tree {
	return m.SpawnCustomTree(DefaultExplorationConst)
}

//SpawnCustomTree creates a new search tree with a given exploration constant.
func (m *MCTS) SpawnCustomTree(explorationConst float64) *Tree {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	t := &Tree{
		gameStates:       make(map[gameState]*node),
		explorationConst: explorationConst,
		randSource:       rand.New(rand.NewSource(m.seed)),
	}
	t.current = initializeNode(gameState{m.init, 0}, t)

	m.seed++
	return t
}

//AddTree adds a searched tree to its list of trees to consider
//when deciding upon an action to take.
func (m *MCTS) AddTree(t *Tree) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.trees = append(m.trees, t)
}

//BestAction takes all of the searched trees and returns
//the best action based on the highest win percentage
//of each action.
//
//BestAction returns nil if it has received no trees
//to search through or if the current state
//it's considering has no legal actions or is terminal.
func (m *MCTS) BestAction() Action {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if len(m.trees) == 0 {
		return nil
	}

	//Safe guard set in place in case we're dealing
	//with a terminal state
	if m.init.IsTerminal() {
		return nil
	}

	//Democracy Section: each tree votes for an action
	actionScore := make(map[Action]int)
	for _, t := range m.trees {
		actionScore[t.bestAction()]++
	}

	//Democracy Section: the action with the most votes wins
	var bestAction Action
	var mostVotes int
	for a, s := range actionScore {
		if s > mostVotes {
			bestAction = a
			mostVotes = s
		}
	}
	return bestAction
}
