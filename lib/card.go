package lib


import (
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "log"
    "sort"
    "strconv"
    "strings"
)


type card struct {
    Name string
    CastingCost mana
    EntersTapped bool
    Pretty string
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
            if c.Pretty == "" {
                c.Pretty = prettyName(c.Name)
            }
            card_cache[c.Name] = c
        }
    }
    // If data for a card is missing, we need to stop and add it immediately
    c, ok := card_cache[card_name]
    if !ok { log.Fatal("no data for: " + card_name) }
    return c
}


func prettyName(s string) string {
    for _, c := range []string{" ", "-", "'", ","} {
        s = strings.ReplaceAll(s, c, "")
    }
    return s
}


func PrettyCards(cards []card) string {
    counts := make(map[string]int)
    for _, c := range cards {
        counts[c.Pretty] += 1
    }
    name_count := []string{}
    for name, count := range counts {
        nc := name
        if count > 1 {
            nc += "*" + strconv.Itoa(count)
        }
        name_count = append(name_count, nc)
    }
    sort.Strings(name_count)
    return strings.Join(name_count, " ")
}


func PrettyCardsOrdered(cards []card) string {
    s := ""
    for _, c := range cards {
        s += " " + c.Pretty
    }
    return s
}




type cardContainerOrdered struct {
    cards []card
}

type cardContainer struct {


}
