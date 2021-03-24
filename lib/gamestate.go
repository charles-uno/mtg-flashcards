package lib


import (
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
    manaPool mana
    onThePlay bool
    turn int
}


func (self *gameState) Pretty() string {
    lines := []string{}
    lines = append(lines, "hand: " + self.hand.Pretty())
    return strings.Join(lines, "\n")
}


func (clone gameState) Draw(n int) gameManager {
    popped, library := clone.library.SplitAfter(n)
    clone.library = library
    clone.hand = clone.hand.Plus(popped...)
    return GameManager(clone)
}


func (clone gameState) PassTurn() gameManager {
    // Empty mana pool then tap out
    clone.manaPool = mana{}
    for c, n := range clone.battlefield.Items() {
        m := c.TapsFor()
        clone.manaPool = clone.manaPool.Plus(m.Times(n))
    }
    // TODO: pay for Pact
    // Reset land drops. Check for Dryad, Scout, Azusa
    clone.landPlays = 1 +
        clone.battlefield.Count(Card("Dryad of the Ilysian Grove")) +
        clone.battlefield.Count(Card("Sakura-Tribe Scout")) +
        2*clone.battlefield.Count(Card("Azusa, Lost but Seeking"))
    return clone.Draw(1)
}


func (gs *gameState) GetLog() string {
    return gs.log
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
