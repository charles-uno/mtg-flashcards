package lib


import (
    "io/ioutil"
    "log"
    "math/rand"
    "strconv"
    "strings"
    "time"
)


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
