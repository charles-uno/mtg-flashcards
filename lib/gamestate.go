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
        library: CardArray(library),
        onThePlay: otp,
    }

    // TODO: compute the hash here rather than dynamically

    return gs
}







func (gs *gameState) getLog() string {
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
        "&",
    )
}
