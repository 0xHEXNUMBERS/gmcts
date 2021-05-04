[![Documentation](https://img.shields.io/badge/Documentation-GoDoc-green.svg)](https://pkg.go.dev/github.com/0xhexnumbers/gmcts/v2)

GMCTS - Monte-Carlo Tree Search (the g stands for whatever you want it to mean :^) )
====================================================================================

GMCTS is an implementation of the Monte-Carlo Tree Search algorithm
with support for any deterministic game.

How To Install
==============

This project requires Go 1.7+ to run. To install, use `go get`:

```bash
go get github.com/0xhexnumbers/gmcts
```

Alternatively, you can clone it yourself into your $GOPATH/src/github.com/0xhexnumbers/ folder to get the latest dev build:

```bash
git clone https://github.com/0xhexnumbers/gmcts
```

How To Use
==========

```go
package pkg

import (
    "github.com/0xhexnumbers/gmcts/v2"
)

func NewGame() gmcts.Game {
    var game gmcts.Game
    //...
    //Setup a new game
    //...
    return game
}

func runGame() {
    gameState := NewGame()

    //MCTS algorithm will play against itself
    //until a terminal state has been reached
    for !gameState.IsTerminal() {
        mcts := gmcts.NewMCTS(gameState)

        //Spawn a new tree and play 1000 game simulations
        tree := mcts.SpawnTree()
        tree.SearchRounds(1000)

        //Add the searched tree into the mcts tree collection
        mcts.AddTree(tree)

        //Get the best action based off of the trees collected from mcts.AddTree()
        bestAction, err := mcts.BestAction()
        if err != nil {
            //...
            //handle error
            //...
        }

        //Update the game state using the tree's best action
        gameState, _ = gameState.ApplyAction(bestAction)
    }
}
```

If you choose to, you can run multiple trees concurrently.

```go
concurrentTrees := 4

mcts := gmcts.NewMCTS(gameState)

//Run 4 trees concurrently
var wait sync.WaitGroup
wait.Add(concurrentTrees)
for i := 0; i < concurrentTrees; i++ {
    go func(){
        tree := mcts.SpawnTree()
        tree.SearchRounds(1000)
        mcts.AddTree(tree)
        wait.Done()
    }()
}
//Wait for the 4 trees to finish searching
wait.Wait()

bestAction, err := mcts.BestAction()
if err != nil {
    //...
    //handle error
    //...
}

gameState, _ = gameState.ApplyAction(bestAction)
```

Testing
=======

You can test this package with `go test`. The test plays a game of tic-tac-toe against itself. The test should:

1. Start the game by placing an x piece in the middle, and
2. Finish in a draw.

If either of these fail, the test fails. It's a rather neat way to make sure everything works as intended!

Documentation
=============

Documentation for this package can be found at [pkg.go.dev](https://pkg.go.dev/github.com/0xhexnumbers/gmcts/v2)

Bug Reports
===========

Email me at 0xhexnumbers@gmail.com :D
