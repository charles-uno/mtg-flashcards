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
    name string
}


func Card(name string) card {
    return card{name: name}
}


func (self *card) Pretty() string {
    return GetCardData(self.name).Pretty
}


type cardArray struct {
    // An ordered sequence of card names, such as a library
    arr []card
}


func CardArray(names []string) cardArray {
    arr := []card{}
    for _, name := range names {
        arr = append(arr, Card(name))
    }
    return cardArray{arr: arr}
}


func (self *cardArray) Pretty() string {
    chunks := []string{}
    for _, c := range self.arr {
        chunks = append(chunks, c.Pretty())
    }
    return strings.Join(chunks, " ")
}


func (ca cardArray) SplitAfter(n int) ([]card, cardArray) {
    popped := ca.arr[:n]
    ca.arr = ca.arr[n:]
    return popped, ca
}






type cards struct {
    // A non-ordered container of cards, such as a hand or battlefield
    cardsCounts map[card]int
}


func Cards(names []string) cards {
    cardsCounts := make(map[card]int)
    for _, name := range names {
        cardsCounts[Card(name)] += 1
    }
    return cards{cardsCounts: cardsCounts}
}


func (self *cards) Pretty() string {
    chunks := []string{}
    for c, n := range self.cardsCounts {
        chunk := c.Pretty()
        if n > 1 {
            chunk += "*" + strconv.Itoa(n)
        }
        chunks = append(chunks, chunk)
    }
    sort.Strings(chunks)
    return strings.Join(chunks, " ")
}


func (self *cards) Plus(elts ...card) cards {
    cardsCounts := make(map[card]int)
    // Deep copy the original
    for c, n := range self.cardsCounts {
        cardsCounts[c] = n
    }
    // Append the new elements
    for _, c := range elts {
        cardsCounts[c] += 1
    }
    return cards{cardsCounts: cardsCounts}
}


type cardData struct {
    // No need to duplicate card metadata over and over. Cache it by card name
    // and look it up as needed.
    Name string
    CastingCost mana
    EntersTapped bool
    Pretty string
    Type string
    TapsFor mana
}


// Cache card data by name so we don't re-read the file repeatedly
var card_cache = make(map[string]cardData)


func InitCardDataCache() {
    card_data_raw := []cardData{}
    text_bytes, err := ioutil.ReadFile("carddata.yaml")
    if err != nil { log.Fatal(err) }
    err = yaml.Unmarshal(text_bytes, &card_data_raw)
    if err != nil { log.Fatal(err) }
    log.Println("loaded carddata.yaml")
    for _, cd := range card_data_raw {
        if cd.Pretty == "" {
            cd.Pretty = slug(cd.Name)
        }
        card_cache[cd.Name] = cd
    }
}


func GetCardData(card_name string) cardData {
    if len(card_cache) == 0 {
        InitCardDataCache()
    }
    // If data for a card is missing, we need to stop and add it immediately
    cd, ok := card_cache[card_name]
    if !ok { log.Fatal("no data for: " + card_name) }
    return cd
}


func slug(s string) string {
    for _, c := range []string{" ", "-", "'", ","} {
        s = strings.ReplaceAll(s, c, "")
    }
    return s
}
