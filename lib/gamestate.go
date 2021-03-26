package lib


import (
    "errors"
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


func (self *gameState) NextSteps() gameManager {
    gm := self.passTurn()
    for c, _ := range self.hand.Items() {
        if c.IsLand() {
            g, err := self.play(c)
            if err != nil {
                log.Fatal(err)
            }
            gm.Update(g)
        } else {
            g, err := self.cast(c)
            if err != nil {
                log.Fatal(err)
            }
            gm.Update(g)
        }
    }
    return gm
}


func (clone gameState) passTurn() gameManager {
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
        return clone.Draw(1)
    } else {
        return GameManager(clone)
    }
}


func (clone gameState) Draw(n int) gameManager {
    popped, library := clone.library.SplitAfter(n)
    clone.library = library
    clone.hand = clone.hand.Plus(popped...)
    return GameManager(clone)
}


func (clone gameState) cast(c card) (gameManager, error) {
    // Is this spell in our hand?
    if clone.hand.Count(c) == 0 {
        return GameManager(), nil
    }
    // Do we have enough mana to cast it?
    m, err := clone.manaPool.Minus(c.CastingCost())
    if err != nil {
        return GameManager(), nil
    }
    clone.manaPool = m
    clone.note("\ncast " + c.Pretty())
    clone.noteManaPool()
    // Now figure out what it does
    switch c.name {
        case "Primeval Titan":
            return clone.castPrimevalTitan(), nil
    }
    return GameManager(), errors.New("not sure how to cast: " + c.name)
}


func (clone gameState) play(c card) (gameManager, error) {
    // Is this land in our hand?
    if clone.hand.Count(c) == 0 {
        return GameManager(), nil
    }
    // Do we have at least one land play remaining?
    if clone.landPlays <= 0 {
        return GameManager(), nil
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
    // Watch out for additional effects, if any
    switch c.name {
        case "Forest":
            return clone.playForest(), nil
    }
    return GameManager(), errors.New("not sure how to play: " + c.name)
}


func (clone gameState) castPrimevalTitan() gameManager {
    clone.done = true
    return GameManager(clone)
}


func (clone gameState) playForest() gameManager {
    return GameManager(clone)
}


func (gs *gameState) Pretty() string {
    return gs.log
}


func (self *gameState) note(s string) {
    self.log += s
}


func (self *gameState) noteManaPool() {
    self.log += ", " + self.manaPool.Pretty() + " in pool"
}



func (gs *gameState) Hash() string {
    // We don't care about order for battlefield or hand, but we do care about
    // the order of the library
    return strings.Join(
        []string{
            gs.hand.Pretty(),
            gs.battlefield.Pretty(),
            gs.manaPool.Pretty(),
            strconv.FormatBool(gs.done),
            gs.library.Pretty(),
        },
        ";",
    )
}
