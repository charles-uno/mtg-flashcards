package lib


import (
    "errors"
    "log"
)


type gameStateSet struct {
    // Use a map to imitate a Python-style set of game states
    states map[string]gameState
}


func GameStateSet(states ...gameState) gameStateSet {
    gss := gameStateSet{
        states: make(map[string]gameState),
    }
    for _, gs := range states {
        gss.add(gs)
    }
    return gss
}


func (self *gameStateSet) Pretty() string {
    gs, err := self.get()
    if err != nil {
        log.Fatal(err)
    }
    return gs.Pretty()
}


func (self *gameStateSet) pop() (gameState, error) {
    for key, gs := range self.states {
        delete(self.states, key)
        return gs, nil
    }
    return gameState{}, errors.New("No game state to pop")
}


func (self *gameStateSet) get() (gameState, error) {
    for _, gs := range self.states {
        return gs, nil
    }
    return gameState{}, errors.New("No game state to get")
}


func (self *gameStateSet) add(gs gameState) {
    self.states[gs.Hash()] = gs
}


func (self *gameStateSet) size() int {
    return len(self.states)
}
