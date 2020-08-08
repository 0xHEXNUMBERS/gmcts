package gmcts

import (
	"math/rand"
	"sync"
)

//NewMCTS returns a new MCTS wrapper
func NewMCTS(initial Game) *MCTS {
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

//SetSeed sets the seed of the next tree to be spawned.
//This value is initially set to 1, and increments on each
//spawned tree.
func (m *MCTS) SetSeed(seed int64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.seed = seed
}

//SpawnCustomTree creates a new search tree with a given exploration constant.
func (m *MCTS) SpawnCustomTree(explorationConst float64) *Tree {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	t := &Tree{
		gameStates:       make(map[gameHash]*node),
		explorationConst: explorationConst,
		randSource:       rand.New(rand.NewSource(m.seed)),
	}
	t.current = initializeNode(gameState{m.init, gameHash{m.init.Hash(), 0}}, t)

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
func (m *MCTS) BestAction() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if len(m.trees) == 0 {
		return -1
	}

	//Safe guard set in place in case we're dealing
	//with a terminal state
	if m.init.IsTerminal() {
		return -1
	}

	//Democracy Section: each tree votes for an action
	actionScore := make([]int, m.init.Len())
	for _, t := range m.trees {
		actionScore[t.bestAction()]++
	}

	//Democracy Section: the action with the most votes wins
	var bestAction int
	var mostVotes int
	for a, s := range actionScore {
		if s > mostVotes {
			bestAction = a
			mostVotes = s
		}
	}
	return bestAction
}
