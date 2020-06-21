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

func initializeNode(g gameState, tree *Tree) *node {
	return &node{
		state:     g,
		tree:      tree,
		nodeScore: make(map[Player]float64),
	}
}

func (a *actionNodePair) UCT2(p Player) float64 {
	//Calculate exploitation component
	//Wins counts as a whole point, while draws count as half a point
	exploit := a.child.nodeScore[p] / float64(a.child.nodeVisits)

	explore := math.Sqrt(
		math.Log(float64(a.parent.nodeVisits)) / float64(a.visits),
	)

	return exploit + a.parent.tree.explorationConst*explore
}

func (n *node) selectNode() ([]Player, float64) {
	var selectedChild *actionNodePair
	var winners []Player
	var scoreToAdd float64
	if len(n.children) == 0 {
		n.expand()
	}

	if n.state.Game.IsTerminal() {
		//Get the result of the game
		winners = n.simulate()
		scoreToAdd = 1.0 / float64(len(winners))
	} else if len(n.unvisitedChildren) > 0 {
		//Grab the first unvisited child and run a simulation from that point
		selectedChild = n.unvisitedChildren[0]
		selectedChild.child.nodeVisits++
		n.unvisitedChildren = n.unvisitedChildren[1:]

		winners = selectedChild.child.simulate()
		scoreToAdd = 1.0 / float64(len(winners))
	} else {
		//Select the child with the max UCT2 score with the current player
		//and get the results to add from its selection
		maxScore := -1.0
		thisPlayer := n.state.Player()
		for _, an := range n.children {
			score := an.UCT2(thisPlayer)
			if score > maxScore {
				maxScore = score
				selectedChild = an
			}
		}
		winners, scoreToAdd = selectedChild.child.selectNode()
	}

	//Update this node along with each parent in this path recursively
	n.nodeVisits++
	if selectedChild != nil {
		selectedChild.visits++
	}

	for _, p := range winners {
		n.nodeScore[p] += scoreToAdd
	}
	return winners, scoreToAdd
}

func (n *node) isParentOf(potentialChild *node) bool {
	for _, an := range n.children {
		if an != nil && an.child == potentialChild {
			return true
		}
	}
	return false
}

func (n *node) expand() {
	actions := n.state.GetActions()
	n.unvisitedChildren = make([]*actionNodePair, len(actions))
	n.children = n.unvisitedChildren
	for i, a := range actions {
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

			n.unvisitedChildren[i] = &actionNodePair{
				action: a,
				parent: n,
				child:  cachedNode,
			}
		} else {
			newNode := initializeNode(newState, n.tree)
			n.unvisitedChildren[i] = &actionNodePair{
				action: a,
				parent: n,
				child:  newNode,
			}

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

		randomIndex := rand.Intn(len(actions))
		game, err = game.ApplyAction(actions[randomIndex])
		if err != nil {
			panic(fmt.Sprintf("gmcts: game returned an error while searching the tree: %s", err))
		}
	}
	return game.Winners()
}
