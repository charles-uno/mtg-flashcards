# MTG Server

Mulligan flashcards for the MTG Modern deck [Amulet Titan][amulet_titan]. The service generates opening hands, and can also evaluate opening hands by playing them out.

[amulet_titan]: https://www.mtggoldfish.com/archetype/amulet-titan#paper

For more information about the model, see my article at [charles.uno][amulet_model]

[amulet_model]: https://charles.uno/amulet-simulation/


## Usage

Launch the service via:

```
go run main.go
```

The service supports two endpoints on port 5001:

- `/api/hand` returns an opening game position:
  - `hand`, a list of seven card names corresponding to the opening hand
  - `library`, a list of the remaining fifty-three cards in the deck
  - `on_the_play`, a boolean indicating whether we are playing first or drawing first
- `/api/play` accepts the same data format returned above. It then shuffles the fifty-three card deck and plays it out. It returns:
  - `success`, indicating whether it was able to cast Primeval Titan by turn four
  - `plays`, a list of maps which describe the computer's sequence of plays over the first few turns of the game. The intention is that these maps can be turned into HTML, complete with formatting for card and mana elements

## Limitations of the Model

The model present here is pretty stripped-down in the interest of performance. For example, it only handles green mana. Adding blue mana into the mix is computationally demanding, and Tolaria West just doesn't matter that often in the first few turns of the game.
