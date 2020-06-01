package gmcts

import (
	"fmt"
	"math"
	"math/rand"
)

const (
	//ExplorationConst is the exploration constant of UCB1 Formula
	//Sqrt(2) is a frequent choice for this constant as specified by
	//https://en.wikipedia.org/wiki/Monte_Carlo_tree_search
	ExplorationConst = math.Sqrt2
)

func initializeNode(g gameState, parent []*node, tree *Tree) *node {
	return &node{
		state:     g,
		tree:      tree,
		parents:   parent,
		children:  make(map[Action]*node),
		nodeScore: make(map[Player]float64),
	}
}

func (n *node) UCB1(p Player) float64 {
	if n.nodeVisits == 0 {
		return math.MaxFloat64
	}
	//Calculate exploitation component
	//Wins counts as a whole point, while draws count as half a point
	exploit := n.nodeScore[p] / float64(n.nodeVisits)

	//Calculate exploration component
	//Because a node may have multiple parents,
	//we need to sum the visits to the parent nodes
	var parentVisits float64 = 0
	for _, p := range n.parents {
		parentVisits += float64(p.nodeVisits)
	}
	explore := math.Sqrt(
		math.Log(parentVisits) / float64(n.nodeVisits),
	)

	return exploit + explore
}

func (n *node) selectNode() *node {
	if n == nil {
		return n
	}
	if n.children == nil || len(n.children) == 0 {
		if n.nodeVisits == 0 || n.state.Game.IsTerminal() {
			return n
		}
		n.expand()
	}

	//Select the child with the max UCB score
	var maxScore float64 = -1
	var maxChild *node = nil
	for _, c := range n.children {
		score := c.UCB1(n.state.Player())
		if score > maxScore {
			maxScore = score
			maxChild = c
		}
	}
	return maxChild.selectNode()
}

func (n *node) expand() {
	for _, action := range n.state.GetActions() {
		newGame, err := n.state.ApplyAction(action)
		if err != nil {
			panic(fmt.Sprintf("gmcts: Game returned an error when exploring the tree: %s", err))
		}

		newState := gameState{newGame, n.state.turn + 1}

		//If we already have a copy in cache, use that and update
		//this node and its parents
		if _, made := n.tree.gameStates[newState]; made {
			n.children[action] = n.tree.gameStates[newState]
			n.children[action].parents = append(
				n.children[action].parents, n,
			)

			//Update this node and each parent with the
			//scores of this already existing child node
			n.updateScoresWithExistingChild(n.children[action])
		} else {
			n.children[action] = initializeNode(newState, []*node{n}, n.tree)

			//Save node for reuse
			n.tree.gameStates[newState] = n.children[action]
		}
	}
}

func (n *node) updateScoresWithExistingChild(child *node) {
	for p, s := range child.nodeScore {
		n.nodeScore[p] += s
	}
	n.nodeVisits += child.nodeVisits

	for _, p := range n.parents {
		p.updateScoresWithExistingChild(child)
	}
}

func (n *node) simulate() []Player {
	game := n.state.Game
	for !game.IsTerminal() {
		var err error

		actions := game.GetActions()
		randomIndex := rand.Intn(len(actions))
		game, err = game.ApplyAction(actions[randomIndex])
		if err != nil {
			panic(fmt.Sprintf("gmcts: game returned an error while searching the tree: %s", err))
		}
	}
	return game.Winners()
}

func (n *node) backpropagation(winners []Player, scoreToAdd float64) {
	n.nodeVisits++

	for _, p := range winners {
		n.nodeScore[p] += scoreToAdd
	}

	for _, p := range n.parents {
		p.backpropagation(winners, scoreToAdd)
	}
}
