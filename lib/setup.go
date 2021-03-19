package lib

import (
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "log"
    "math/rand"
    "strconv"
    "strings"
    "time"
)


type mana struct {
    Green int
    Total int
}


type card struct {
    Name string
    CastingCost mana
    EntersTapped bool
    Type string
    TapsFor mana
}


// Cache card data by name so we don't re-read the file repeatedly
var card_cache = make(map[string]card)


func Card(card_name string) card {
    if len(card_cache) == 0 {
        card_data := []card{}
        text_bytes, err := ioutil.ReadFile("carddata.yaml")
        if err != nil { log.Fatal(err) }
        err = yaml.Unmarshal(text_bytes, &card_data)
        if err != nil { log.Fatal(err) }
        log.Println("loaded carddata.yaml")
        for _, c := range card_data {
            card_cache[c.Name] = c
        }
    }
    // If data for a card is missing, we need to stop and add it immediately
    c, ok := card_cache[card_name]
    if !ok { log.Fatal("no data for: " + card_name) }
    return c
}


func LoadDeck() []card {
    deck := []card{}
    for _, card_name := range loadCardNames() {
        deck = append(deck, Card(card_name))
    }
    return shuffled(deck)
}


func loadCardNames() []string {
    list := []string{}
    for _, line := range readLines("decklist.txt") {
        n_card := strings.SplitN(line, " ", 2)
        n, err := strconv.Atoi(n_card[0])
        if err != nil { log.Fatal(err) }
        for i := 0; i<n; i++ {
            list = append(list, n_card[1])
        }
    }
    return list
}


func shuffled(deck []card) []card {
    rand.Seed(time.Now().UTC().UnixNano())
    ret := make([]card, len(deck))
    for i, j := range rand.Perm(len(deck)) {
        ret[i] = deck[j]
    }
    return ret
}


func readLines(filename string) []string {
    lines := []string{}
    text_bytes, err := ioutil.ReadFile(filename)
    if err != nil { log.Fatal(err) }
    for _, line := range strings.Split(string(text_bytes), "\n") {
        // Skip empty lines and comments
        if len(line) == 0 { continue }
        if line[:1] == "#" { continue }
        lines = append(lines, line)
    }
    return lines
}
