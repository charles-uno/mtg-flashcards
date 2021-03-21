package main

import (
    "fmt"
    "math/rand"
    "time"

    "github.com/charles-uno/mtgserver/lib"
)


// Note: we want to be able to run multiple models for the same opening hand.
// That'll be easier if we send the hand and the library into the game state
// constructor separately. Rather than, say, passing in a list of 60 cards and
// having the constructor shuffle and draw.

func main() {

    deck := lib.LoadDeck()
    hand, library := deck[:7], deck[7:]
    onThePlay := flip()

    gs := lib.GameState(hand, library, onThePlay)

    fmt.Println(gs)

}


func flip() bool {
    // Random generator should be seeded from shuffling, but let's be sure. We
    // only call this once per game anyway.
    rand.Seed(time.Now().UTC().UnixNano())
    return rand.Intn(2) == 0
}
