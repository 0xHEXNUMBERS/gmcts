package gmcts

import (
	"fmt"
	"math"
	"math/rand"
)

const (
	//DefaultExplorationConst is the default exploration constant of UCB1 Formula
	//Sqrt(2) is a frequent choice for this constant as specified by
	//https://en.wikipedia.org/wiki/Monte_Carlo_tree_search
	DefaultExplorationConst = math.Sqrt2
)

func initializeNode(g gameState, parent []*node, tree *Tree) *node {
	return &node{
		state:     g,
		tree:      tree,
		parents:   parent,
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
	var parentVisits float64
	for _, p := range n.parents {
		parentVisits += float64(p.nodeVisits)
	}
	explore := math.Sqrt(
		math.Log(parentVisits) / float64(n.nodeVisits),
	)

	return exploit + n.tree.explorationConst*explore
}

func (n *node) selectNode() *node {
	if n.children == nil || len(n.children) == 0 {
		if n.nodeVisits == 0 || n.state.Game.IsTerminal() {
			return n
		}
		n.expand()
	}

	//Select the child with the max UCB score with the current player
	var maxChild *node
	maxScore := -1.0
	thisPlayer := n.state.Player()
	for _, an := range n.children {
		score := an.node.UCB1(thisPlayer)
		if score > maxScore {
			maxScore = score
			maxChild = an.node
		}
	}
	return maxChild.selectNode()
}

func (n *node) isParentOf(potentialChild *node) bool {
	for _, p := range potentialChild.parents {
		if n == p {
			return true
		}
	}
	return false
}

func (n *node) expand() {
	actions := n.state.GetActions()
	n.children = make([]actionNodePair, len(actions))
	for i, action := range actions {
		newGame, err := n.state.ApplyAction(action)
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

			n.children[i] = actionNodePair{action, cachedNode}
			cachedNode.parents = append(
				cachedNode.parents, n,
			)

			//Update this node and each parent with the
			//scores of this already existing child node
			n.updateScoresWithExistingChild(cachedNode)
		} else {
			newNode := initializeNode(newState, []*node{n}, n.tree)
			n.children[i] = actionNodePair{
				action: action,
				node:   newNode,
			}

			//Save node for reuse
			n.tree.gameStates[newState] = newNode
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
		panicIfNoActions(game, actions)

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
