package lib


import (
    "log"
    "strings"
)


type gameManager struct {
    maxTurns int
    // Use a map to imitate a Python-style set of game states
    states map[string]gameState
    success bool
    turn int
}


func NewGame(libraryRaw []string, handRaw []string, otp bool, verbose bool, maxTurns int) (gameManager, error) {
    allCardNames := []string{}
    handCards := []card{}
    for _, cardName := range handRaw {
        handCards = append(handCards, Card(cardName))
        allCardNames = append(allCardNames, cardName)
    }
    libraryCards := []card{}
    for _, cardName := range libraryRaw {
        libraryCards = append(libraryCards, Card(cardName))
        allCardNames = append(allCardNames, cardName)
    }
    err := EnsureCardData(allCardNames)
    if err != nil {
        return gameManager{}, err
    }
    state := NewGameState(libraryCards, handCards, otp, verbose, maxTurns)
    return GameManager(state), nil
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


func (self *gameManager) NextTurn() gameManager {
    if self.Size() == 0 {
        log.Fatal("called NextTurn on empty gameManager")
    }
    // Once we find a line, we're done iterating
    if self.success {
        return *self
    }
    if self.turn > 0 {
        log.Println("starting turn", self.turn, "with", self.Size(), "states")
    }
    ret := GameManager()
    for self.Size() > 0 {
        stateOld := self.Pop()
        for _, stateNew := range stateOld.NextStates() {
            // If we find a state that gets there, we're done
            if stateNew.success {
                return GameManager(stateNew)
            }
            if stateNew.turn == self.turn {
                self.Add(stateNew)
            } else {
                ret.Add(stateNew)
            }
        }
    }
    // After turn four or so, further work is expensive but not interesting.
    // Pop off the longest log we can find to show we tried.
    if ret.turn > self.maxTurns {
        log.Println("giving up on turn", ret.turn, "with", ret.Size(), "states")
        bestState := ret.Pop()
        for ret.Size() > 0 {
            state := ret.Pop()
            if state.LogSize() > bestState.LogSize() {
                bestState = state
            }
        }
        bestState.MarkDeadEnd()
        return GameManager(bestState)
    }
    return ret
}


func (self *gameManager) Pretty() string {
    lines := []string{}
    for _, state := range self.states {
        lines = append(lines, state.Pretty())
    }
    return strings.Join(lines, "\n~~~\n")
}


func (self *gameManager) ToJSON() string {
    lines := []string{}
    for _, state := range self.states {
        lines = append(lines, state.ToJSON())
    }
    return strings.Join(lines, "\n~~~\n")
}


func (self *gameManager) ToMiniJSON() string {
    lines := []string{}
    for _, state := range self.states {
        lines = append(lines, state.ToMiniJSON())
    }
    return strings.Join(lines, "\n~~~\n")
}


func (self *gameManager) Add(state gameState) {
    self.states[state.Hash()] = state
    self.maxTurns = state.maxTurns
    // By construction, in-progress states and completed states never mix
    self.success = state.success
    // Turn is uniform for all states within a gameManager
    self.turn = state.turn
}


func (self *gameManager) Pop() gameState {
    for hash, state := range self.states {
        delete(self.states, hash)
        return state
    }
    log.Fatal("pop from empty gameManager")
    return gameState{}
}


func (self *gameManager) Size() int {
    return len(self.states)
}


func (self *gameManager) IsDone() bool {
    return self.turn > self.maxTurns || self.success
}


func (self *gameManager) Update(other gameManager) {
    for hash, state := range other.states {
        self.states[hash] = state
    }
}
