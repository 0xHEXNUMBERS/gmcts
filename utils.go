package gmcts

import "fmt"

func panicIfNoActions(game Game, actions []Action) {
	if len(actions) == 0 {
		panic(fmt.Sprintf("gmcts: game returned no actions on a non-terminal state: %#v", game))
	}
}
