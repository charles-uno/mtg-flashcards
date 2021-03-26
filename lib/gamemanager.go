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


func NewGame(hand_raw []card, library_raw []card, otp bool) gameManager {
    hand := CardMap(hand_raw)
    var play_order string
    if otp {
        play_order = "on the play"
    } else {
        play_order = "on the draw"
    }
    state := gameState{
        hand: hand,
        // Empty string is fine for the initial game state
        hash: "",
        landPlays: 0,
        library: CardArray(library_raw),
        log: play_order + ", draw " + hand.Pretty(),
        onThePlay: otp,
        turn: 0,
    }
    state.exportBreak()
    state.exportText(play_order + ", draw ")
    state.exportCardMap(hand)
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


func (self *gameManager) NextTurn() gameManager {
    if self.Size() == 0 {
        log.Fatal("called NextTurn on empty gameManager")
    }
    // Once we find a line, we're done iterating
    if self.done {
        return *self
    }
    ret := GameManager()
    for self.Size() > 0 {
        state_old := self.Pop()
        for _, state_new := range state_old.NextStates() {
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
    return ret
}


func (self *gameManager) Pretty() string {
    lines := []string{}
    for _, state := range self.states {
        lines = append(lines, state.Pretty())
    }
    return strings.Join(lines, "\n~~~\n")
}


func (self *gameManager) Export() string {
    lines := []string{}
    for _, state := range self.states {
        lines = append(lines, state.Export())
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
