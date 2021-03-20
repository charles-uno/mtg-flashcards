package lib


import (
    "strings"
    "strconv"
)


// The gameState is an immutable object which describes a snapshot in time
// during a game. Any change in game state, like drawing a card or casting a
// spell, is enacted by creating a new state.
type gameState struct {
    battlefield []card
    done bool
    hand []card
    hash string
    landPlays int
    library []card
    log string
    manaPool mana
    onThePlay bool
    turn int
}


func GameState(hand []card, library []card, otp bool) gameState {
    gs := gameState{
        hand: hand,
        library: library,
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
            PrettyCards(gs.hand),
            PrettyCards(gs.battlefield),
            gs.manaPool.Pretty(),
            strconv.FormatBool(gs.done),
            PrettyCardsOrdered(gs.library),
        },
        "&",
    )
}
