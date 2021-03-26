package lib


import (
    "log"
    "strings"
    "strconv"
)


// The gameState is an immutable object which describes a snapshot in time
// during a game. Any change in game state, like drawing a card or casting a
// spell, is enacted by creating a new state.
type gameState struct {
    battlefield cardMap
    done bool
    hand cardMap
    hash string
    landPlays int
    library cardArray
    log string
    manaDebt mana
    manaPool mana
    onThePlay bool
    turn int
}


func (self *gameState) NextStates() []gameState {
    ret := self.passTurn()
    for c, _ := range self.hand.Items() {
        if c.IsLand() {
            for _, state := range self.play(c) {
                ret = append(ret, state)
            }
        } else {
            for _, state := range self.cast(c) {
                ret = append(ret, state)
            }
        }
    }
    return ret
}


func (clone gameState) passTurn() []gameState {
    clone.turn += 1
    clone.note("\n--- turn " + strconv.Itoa(clone.turn))
    // Empty mana pool then tap out
    clone.manaPool = mana{}
    for c, n := range clone.battlefield.Items() {
        m := c.TapsFor()
        clone.manaPool = clone.manaPool.Plus(m.Times(n))
    }
    clone.noteManaPool()
    // TODO: pay for Pact
    // Reset land drops. Check for Dryad, Scout, Azusa
    clone.landPlays = 1 +
        clone.battlefield.Count(Card("Dryad of the Ilysian Grove")) +
        clone.battlefield.Count(Card("Sakura-Tribe Scout")) +
        2*clone.battlefield.Count(Card("Azusa, Lost but Seeking"))
    if clone.turn > 1 || !clone.onThePlay {
        return clone.draw(1)
    } else {
        return []gameState{clone}
    }
}


func (clone gameState) cast(c card) []gameState {
    // Is this spell in our hand?
    if clone.hand.Count(c) == 0 {
        return []gameState{}
    }
    // Do we have enough mana to cast it?
    cost := c.CastingCost()
    m, err := clone.manaPool.Minus(cost)
    if err != nil {
        return []gameState{}
    }
    clone.manaPool = m
    clone.note("\ncast " + c.Pretty())
    clone.noteManaPool()
    // Now figure out what it does
    switch c.name {
        case "Explore":
            return clone.castExplore()
        case "Primeval Titan":
            return clone.castPrimevalTitan()
    }
    log.Fatal("not sure how to cast: " + c.name)
    return []gameState{}
}


func (clone gameState) play(c card) []gameState {
    // Is this land in our hand?
    if clone.hand.Count(c) == 0 {
        return []gameState{}
    }
    // Do we have at least one land play remaining?
    if clone.landPlays <= 0 {
        return []gameState{}
    }
    // Tap out immediately
    if c.EntersTapped() {
        nAmulets := clone.battlefield.Count(Card("Amulet of Vigor"))
        m := c.TapsFor()
        clone.manaPool = clone.manaPool.Plus(m.Times(nAmulets))
    } else {
        clone.manaPool = clone.manaPool.Plus(c.TapsFor())
    }
    clone.note("\nplay " + c.Pretty())
    clone.noteManaPool()
    clone.landPlays -= 1
    clone.battlefield = clone.battlefield.Plus(c)
    // Watch out for additional effects, if any
    switch c.name {
        case "Forest":
            return clone.playForest()
    }
    log.Fatal("not sure how to play: " + c.name)
    return []gameState{}
}


func (clone gameState) castPrimevalTitan() []gameState {
    clone.done = true
    return []gameState{clone}
}


func (clone gameState) castExplore() []gameState {
    clone.landDrops += 1
    return clone.draw(1)
}


func (clone gameState) playForest() []gameState {
    return []gameState{clone}
}


func (gs *gameState) Pretty() string {
    return gs.log
}


func (clone gameState) draw(n int) []gameState {
    popped, library := clone.library.SplitAfter(n)
    clone.library = library
    clone.hand = clone.hand.Plus(popped...)

    popped_map := CardMap(popped)
    clone.note(", draw " + popped_map.Pretty())

    return []gameState{clone}
}


func (self *gameState) note(s string) {
    self.log += s
}


func (self *gameState) noteManaPool() {
    self.log += ", " + self.manaPool.Pretty() + " in pool"
}


func (state *gameState) Hash() string {
    // We don't care about order for battlefield or hand, but we do care about
    // the order of the library
    return strings.Join(
        []string{
            state.hand.Pretty(),
            state.battlefield.Pretty(),
            state.manaPool.Pretty(),
            strconv.FormatBool(state.done),
            strconv.Itoa(state.landPlays),
            state.library.Pretty(),
        },
        ";",
    )
}
