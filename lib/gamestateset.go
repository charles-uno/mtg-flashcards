package lib


import (
    "errors"
    "log"
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


func (self *gameManager) Next() gameManager {
    manager := GameManager()

    // Passing the turn is always an option

    // Cast spells

    // Play lands

    return manager
}




func (self *gameManager) Pretty() string {
    lines := []string{}
    for _, state := self.states {
        lines = append(lines, state.Pretty())
    }
    return strings.Join(lines, "\n---\n")
}


func (self *gameManager) Pop() (gameState, error) {
    for key, gs := range self.states {
        delete(self.states, key)
        return gs, nil
    }
    return gameState{}, errors.New("No game state to pop")
}


func (self *gameManager) Get() (gameState, error) {
    for _, gs := range self.states {
        return gs, nil
    }
    return gameState{}, errors.New("No game state to get")
}


func (self *gameManager) Add(gs gameState) {
    self.states[gs.Hash()] = gs
}


func (self *gameManager) Update(other gameManager) {
    for hash, state := range other.states {
        self.states[hash] = state
    }
}


func (self *gameManager) Size() int {
    return len(self.states)
}


func (self *gameManager) Draw(n int) gameManager {
    ret := gameManager()
    for _, state := range self.states {
        ret.Update(state.Draw(n))
    }
    return ret
}
