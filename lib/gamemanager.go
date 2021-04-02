package lib


import (
    "log"
    "strings"
)


type gameManager struct {
    // Use a map to imitate a Python-style set of game states
    states map[string]gameState
    Turn int
    done bool
}


func NewGame(handRaw []string, libraryRaw []string, otp bool) (gameManager, error) {
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
    hand := CardMap(handCards)
    var playOrder string
    if otp {
        playOrder = "on the play"
    } else {
        playOrder = "on the draw"
    }
    state := gameState{
        hand: hand,
        // Empty string is fine for the initial game state
        hash: "",
        landPlays: 0,
        library: CardArray(libraryCards),
        onThePlay: otp,
        turn: 0,
    }
    state.logText(playOrder + ", opening hand: ")
    state.logCardMap(hand)
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


func (self *gameManager) NextTurn(maxTurns int) gameManager {
    if self.Size() == 0 {
        log.Fatal("called NextTurn on empty gameManager")
    }
    // Once we find a line, we're done iterating
    if self.done {
        return *self
    }
    if self.Turn > 0 {
        log.Println("starting turn", self.Turn, "with", self.Size(), "states")
    }
    ret := GameManager()
    for self.Size() > 0 {
        state_old := self.Pop()
        for _, state_new := range state_old.NextStates(maxTurns) {
            // If we find a state that gets there, we're done
            if state_new.done {
                return GameManager(state_new)
            }
            if state_new.turn == self.Turn {
                self.Add(state_new)
            } else {
                ret.Add(state_new)
            }
        }
    }
    // After turn four or so, further work is expensive but not interesting.
    // Pop off the longest log we can find to show we tried.
    if ret.Turn > maxTurns {
        log.Println("giving up on turn", ret.Turn, "with", ret.Size(), "states")
        bestState := ret.Pop()
        for ret.Size() > 0 {
            state := ret.Pop()
            if state.LogSize() > bestState.LogSize() {
                bestState = state
            }
        }
        bestState.GiveUp()
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


func (self *gameManager) Add(state gameState) {
    self.states[state.Hash()] = state
    self.done = state.done
    // Turn is uniform for all states within a gameManager
    self.Turn = state.turn
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


func (self *gameManager) IsNotDone() bool {
    return !self.done
}


func (self *gameManager) Update(other gameManager) {
    for hash, state := range other.states {
        self.states[hash] = state
    }
}
