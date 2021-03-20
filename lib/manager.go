package lib


import (
    "errors"
)


// Use a map to imitate a Python-style set of game states
type uniqueGameStates struct {
    states map[string]gameState
}


func UniqueGameStates(states ...gameState) uniqueGameStates {
    ugs := uniqueGameStates{
        states: make(map[string]gameState),
    }
    for _, gs := range states {
        ugs.add(gs)
    }
    return ugs
}


func (ugs *uniqueGameStates) pop() (gameState, error) {
    for key, gs := range ugs.states {
        delete(ugs.states, key)
        return gs, nil
    }
    return gameState{}, errors.New("No game state to pop")
}


func (ugs *uniqueGameStates) get() (gameState, error) {
    for _, gs := range ugs.states {
        return gs, nil
    }
    return gameState{}, errors.New("No game state to get")
}


func (ugs *uniqueGameStates) add(gs gameState) {
    ugs.states[gs.Hash()] = gs
}


func (ugs *uniqueGameStates) size() int {
    return len(ugs.states)
}
