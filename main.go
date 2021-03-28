package main

import (
    "fmt"
    "log"
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

    game := lib.NewGame(hand, library, onThePlay)

    for game.IsNotDone() {
        if game.Turn > 0 {
            log.Println("starting turn", game.Turn, "with", game.Size(), "states")
        }
        game = game.NextTurn()

    }

    fmt.Println("\n\n")
    fmt.Println(game.ToJSON())
    fmt.Println("\n\n")
    fmt.Println(game.Pretty())

}


func flip() bool {
    // Random generator should be seeded from shuffling, but let's be sure. We
    // only call this once per game anyway.
    rand.Seed(time.Now().UTC().UnixNano())
    return rand.Intn(2) == 0
}
