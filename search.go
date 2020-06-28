package gmcts

import (
	"fmt"
	"math"
)

const (
	//DefaultExplorationConst is the default exploration constant of UCB1 Formula
	//Sqrt(2) is a frequent choice for this constant as specified by
	//https://en.wikipedia.org/wiki/Monte_Carlo_tree_search
	DefaultExplorationConst = math.Sqrt2
)

func initializeNode(g gameState, tree *Tree) *node {
	return &node{
		state:     g,
		tree:      tree,
		nodeScore: make(map[Player]float64),
	}
}

func (n *node) UCT2(i int, p Player) float64 {
	exploit := n.children[i].nodeScore[p] / n.children[i].nodeVisits

	explore := math.Sqrt(
		math.Log(n.nodeVisits) / n.childVisits[i],
	)

	return exploit + n.tree.explorationConst*explore
}

func (n *node) selectNode() ([]Player, float64) {
	var selectedChildIndex int
	var winners []Player
	var scoreToAdd float64
	if n.actionCount == 0 {
		n.expand()
	}

	if n.state.Game.IsTerminal() {
		//Get the result of the game
		winners = n.simulate()
		scoreToAdd = 1.0 / float64(len(winners))
	} else if len(n.unvisitedChildren) > 0 {
		//Grab the first unvisited child and run a simulation from that point
		selectedChildIndex = n.actionCount - len(n.unvisitedChildren)
		n.children[selectedChildIndex].nodeVisits++
		n.unvisitedChildren = n.unvisitedChildren[1:]

		winners = n.children[selectedChildIndex].simulate()
		scoreToAdd = 1.0 / float64(len(winners))
	} else {
		//Select the child with the max UCT2 score with the current player
		//and get the results to add from its selection
		maxScore := -1.0
		thisPlayer := n.state.Player()
		for i := 0; i < n.actionCount; i++ {
			score := n.UCT2(i, thisPlayer)
			if score > maxScore {
				maxScore = score
				selectedChildIndex = i
			}
		}
		winners, scoreToAdd = n.children[selectedChildIndex].selectNode()
	}

	//Update this node along with each parent in this path recursively
	n.nodeVisits++
	if n.actionCount != 0 {
		n.childVisits[selectedChildIndex]++
	}

	for _, p := range winners {
		n.nodeScore[p] += scoreToAdd
	}
	return winners, scoreToAdd
}

func (n *node) isParentOf(potentialChild *node) bool {
	for _, child := range n.children {
		if child != nil && child == potentialChild {
			return true
		}
	}
	return false
}

func (n *node) expand() {
	n.actions = n.state.GetActions()
	n.actionCount = len(n.actions)
	n.unvisitedChildren = make([]*node, n.actionCount)
	n.children = n.unvisitedChildren
	n.childVisits = make([]float64, n.actionCount)
	for i, a := range n.actions {
		newGame, err := n.state.ApplyAction(a)
		if err != nil {
			panic(fmt.Sprintf("gmcts: Game returned an error when exploring the tree: %s", err))
		}

		newState := gameState{newGame, n.state.turn + 1}

		//If we already have a copy in cache, use that and update
		//this node and its parents
		if cachedNode, made := n.tree.gameStates[newState]; made {
			if n.isParentOf(cachedNode) {
				continue
			}

			n.unvisitedChildren[i] = cachedNode
		} else {
			newNode := initializeNode(newState, n.tree)
			n.unvisitedChildren[i] = newNode

			//Save node for reuse
			n.tree.gameStates[newState] = newNode
		}
	}
}

func (n *node) simulate() []Player {
	game := n.state.Game
	for !game.IsTerminal() {
		var err error

		actions := game.GetActions()
		panicIfNoActions(game, actions)

		randomIndex := n.tree.randSource.Intn(len(actions))
		game, err = game.ApplyAction(actions[randomIndex])
		if err != nil {
			panic(fmt.Sprintf("gmcts: game returned an error while searching the tree: %s", err))
		}
	}
	return game.Winners()
}
