package lib


import (
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "log"
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


func (self *card) TapsFor() mana {
    return GetCardData(self.name).TapsFor
}


func (self *card) CastingCost() mana {
    return GetCardData(self.name).CastingCost
}


func (self *card) IsLand() bool {
    return GetCardData(self.name).Type == "land"
}


func (self *card) EntersTapped() bool {
    return GetCardData(self.name).EntersTapped
}


type cardData struct {
    // No need to duplicate card metadata over and over. Cache it by card name
    // and look it up as needed.
    Name string         `yaml:"name"`
    CastingCost mana    `yaml:"casting_cost"`
    EntersTapped bool   `yaml:"enters_tapped"`
    Pretty string       `yaml:"pretty"`
    Type string         `yaml:"type"`
    TapsFor mana        `yaml:"taps_for"`
}


// Cache card data by name so we don't re-read the file repeatedly
var card_cache = make(map[string]cardData)


func InitCardDataCache() {
    card_data_raw := []cardData{}
    text_bytes, err := ioutil.ReadFile("carddata.yaml")
    if err != nil {
        log.Fatal(err)
    }
    err = yaml.Unmarshal(text_bytes, &card_data_raw)
    if err != nil {
        log.Fatal(err)
    }
    log.Println("loaded carddata.yaml")
    for _, cd := range card_data_raw {
        if cd.Pretty == "" {
            cd.Pretty = slug(cd.Name)
        }
        card_cache[cd.Name] = cd
    }

    log.Println(card_cache["Primeval Titan"])

}


func GetCardData(card_name string) cardData {
    if len(card_cache) == 0 {
        InitCardDataCache()
    }
    // If data for a card is missing, we need to stop and add it immediately
    cd, ok := card_cache[card_name]
    if !ok {
        log.Fatal("no data for: " + card_name)
    }
    return cd
}


func slug(s string) string {
    for _, c := range []string{" ", "-", "'", ","} {
        s = strings.ReplaceAll(s, c, "")
    }
    return s
}
