package lib

import (
    "errors"
    "fmt"
    "gopkg.in/yaml.v2"
    "io/ioutil"
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
        err = yaml.Unmarshal(text_bytes, &card_data)
        if err != nil {
            panic(err)
        }
        fmt.Println(card_data)
        for _, c := range card_data {
            card_cache[c.Name] = c
        }
    }
    // If data for a card is missing, we need to stop and add it immediately
    c, ok := card_cache[card_name]
    if !ok {
        panic(errors.New("no data for: " + card_name))
    }
    return c
}


func LoadDeck() ([]card, error) {
    card_names, err := loadCardNames()
    if err != nil {
        return []card{}, err
    }
    deck := []card{}
    for _, card_name := range card_names {
        deck = append(deck, Card(card_name))
    }
    return shuffled(deck), nil
}


func loadCardNames() ([]string, error) {
    lines, err := readLines("decklist.txt")
    if err != nil {
        return []string{}, err
    }
    list := []string{}
    for _, line := range lines {
        n_card := strings.SplitN(line, " ", 2)
        n, err := strconv.Atoi(n_card[0])
        if err != nil {
            return []string{}, err
        }
        for i := 0; i<n; i++ {
            list = append(list, n_card[1])
        }
    }
    return list, nil
}


func shuffled(deck []card) []card {
    rand.Seed(time.Now().UTC().UnixNano())
    ret := make([]card, len(deck))
    for i, j := range rand.Perm(len(deck)) {
        ret[i] = deck[j]
    }
    return ret
}


func readLines(filename string) ([]string, error) {
    lines := []string{}
    text_bytes, err := ioutil.ReadFile(filename)
    if err != nil { return lines, err }
    for _, line := range strings.Split(string(text_bytes), "\n") {
        // Skip empty lines and comments
        if len(line) == 0 { continue }
        if line[:1] == "#" { continue }
        lines = append(lines, line)
    }
    return lines, nil
}
