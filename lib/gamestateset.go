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
        gss.Add(gs)
    }
    return gss
}


func (self *gameStateSet) Pretty() string {
    gs, err := self.Get()
    if err != nil {
        log.Fatal(err)
    }
    return gs.Pretty()
}


func (self *gameStateSet) Pop() (gameState, error) {
    for key, gs := range self.states {
        delete(self.states, key)
        return gs, nil
    }
    return gameState{}, errors.New("No game state to pop")
}


func (self *gameStateSet) Get() (gameState, error) {
    for _, gs := range self.states {
        return gs, nil
    }
    return gameState{}, errors.New("No game state to get")
}


func (self *gameStateSet) Add(gs gameState) {
    self.states[gs.Hash()] = gs
}


func (self *gameStateSet) Update(other gameStateSet) {
    for hash, state := range other.states {
        self.states[hash] = state
    }
}


func (self *gameStateSet) Size() int {
    return len(self.states)
}


func (self *gameStateSet) Draw(n int) gameStateSet {
    ret := GameStateSet()
    for _, state := range self.states {
        ret.Update(state.Draw(n))
    }
    return ret
}
