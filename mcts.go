package gmcts

import (
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
	action := initial.GetActions()
	if !reflect.TypeOf(action).Elem().Comparable() {
		panic("gmcts: action type is not comparable")
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
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	t := &Tree{gameStates: make(map[gameState]*node), explorationConst: explorationConst}
	t.current = initializeNode(gameState{m.init, 0}, nil, t)
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
	rootState := m.init
	baseActions := rootState.GetActions()
	if len(baseActions) == 0 || rootState.IsTerminal() {
		return nil
	}

	//Loop through each action and node and calculate the best
	//winrate each action had when searching the trees
	bestAction := baseActions[0]
	bestWinRate := 0.0
	playerTakingAction := rootState.Player()
	for _, a := range baseActions {
		var score float64
		var visits int

		for i := range m.trees {
			node := m.trees[i].current.children[a]
			if node == nil {
				continue
			}
			score += node.nodeScore[playerTakingAction]
			visits += node.nodeVisits
		}

		if visits == 0 {
			continue
		}

		winRate := score / float64(visits)
		if winRate > bestWinRate {
			bestAction = a
			bestWinRate = winRate
		}
	}
	return bestAction
}
