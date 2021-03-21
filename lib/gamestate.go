package lib


import (
    "strings"
    "strconv"
)


// The gameState is an immutable object which describes a snapshot in time
// during a game. Any change in game state, like drawing a card or casting a
// spell, is enacted by creating a new state.
type gameState struct {
    battlefield cards
    done bool
    hand cards
    hash string
    landPlays int
    library cardArray
    log string
    manaPool mana
    onThePlay bool
    turn int
}


func GameState(hand []string, library []string, otp bool) gameState {
    gs := gameState{
        hand: Cards(hand),
        // Empty string is fine for the initial game state
        hash: "",
        library: CardArray(library),
        onThePlay: otp,
    }
    return gs
}


func (self *gameState) Pretty() string {
    lines := []string{}
    lines = append(lines, "hand: " + self.hand.Pretty())
    return strings.Join(lines, "\n")
}


func (clone gameState) Draw(n int) gameStateSet {
    popped, library := clone.library.SplitAfter(n)
    clone.library = library
    clone.hand = clone.hand.Plus(popped...)
    return GameStateSet(clone)
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
