package gmcts

import (
	"context"
	"time"
)

//Search searches the tree for a specified time
//
//Search will panic if the Game's ApplyAction
//method returns an error
func (t *Tree) Search(duration time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	t.SearchContext(ctx)
}

//SearchContext searches the tree using a given context
//
//SearchContext will panic if the Game's ApplyAction
//method returns an error
func (t *Tree) SearchContext(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			t.search()
		}
	}
}

//SearchRounds searches the tree for a specified number of rounds
//
//SearchRounds will panic if the Game's ApplyAction
//method returns an error
func (t *Tree) SearchRounds(rounds int) {
	for i := 0; i < rounds; i++ {
		t.search()
	}
}

//search performs 1 round of the MCTS algorithm
func (t *Tree) search() {
	t.current.selectNode()
}
