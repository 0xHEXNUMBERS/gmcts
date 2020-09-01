//Package gmcts is a generic implementation of the
//Monte-Carlo Tree Search (mcts) algorithm.
//
//This package attempts to save memory and time by caching states as to not
//have duplicate nodes in the search tree. This optimization is efficient for
//games like tic-tac-toe, checkers, and go among others.
//
//This package also allows support for tree parallelization. Trees may
//be spawned and ran in their own goroutine. After searching, they may be
//compiled together to produce a more informed action than just searching
//through one tree.
package gmcts
