package lib


import (
    "strings"
)


type gameManager struct {
    // Use a map to imitate a Python-style set of game states
    states map[string]gameState
}


func NewGame(hand []card, library []card, otp bool) gameManager {
    state := gameState{
        hand: CardMap(hand),
        // Empty string is fine for the initial game state
        hash: "",
        library: CardArray(library),
        onThePlay: otp,
    }
    return GameManager(state)
}


func GameManager(states ...gameState) gameManager {
    manager := gameManager{
        states: make(map[string]gameState),
    }
    for _, state := range states {
        manager.Add(state)
    }
    return manager
}


func (self *gameManager) NextSteps() gameManager {
    gm := GameManager()
    for _, gs := range self.states {
        gm.Update(gs.NextSteps())
    }
    return gm
}


func (self *gameManager) Pretty() string {
    lines := []string{}
    for _, state := range self.states {
        lines = append(lines, state.Pretty()[1:])
    }
    return strings.Join(lines, "\n~~~\n")
}


func (self *gameManager) Add(gs gameState) {
    self.states[gs.Hash()] = gs
}


func (self *gameManager) Update(other gameManager) {
    for hash, state := range other.states {
        self.states[hash] = state
    }
}


func (self *gameManager) Draw(n int) gameManager {
    ret := GameManager()
    for _, state := range self.states {
        ret.Update(state.Draw(n))
    }
    return ret
}
