package lib

type gameState struct {
    battlefield []card
    hand []card
    library []card
    Log string
    manaPool mana
    landPlays int
    turn int
}


func GameState() {}
