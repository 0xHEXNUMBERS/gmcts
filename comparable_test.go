package gmcts

import "testing"

//comparableState tests the comparable action requirement of gmcts, as
//the GetActions method returns a noncomparable action.
type comparableState struct{}

//nonComparableState tests the comparable state requirement of gmcts.
type nonComparableState struct {
	comparableState
	_ []int
}

func (n comparableState) GetActions() []Action {
	return []Action{nonComparableState{}}
}

func (n comparableState) ApplyAction(a Action) (Game, error) {
	return n, nil
}

func (n comparableState) IsTerminal() bool {
	return true
}

func (n comparableState) Player() Player {
	return 0
}

func (n comparableState) Winners() []Player {
	return nil
}

func TestNonComparableState(t *testing.T) {
	//Calling NewMCTS should panic, as the nonComparableState is, as
	//the name suggests, not comparable.
	defer func() {
		if r := recover(); r == nil {
			t.FailNow()
		}
	}()
	NewMCTS(nonComparableState{})
}

func TestNonComparableAction(t *testing.T) {
	//Calling NewMCTS should panic, as the actions from comparableState
	//are not comparable.
	defer func() {
		if r := recover(); r == nil {
			t.FailNow()
		}
	}()
	NewMCTS(comparableState{})
}
