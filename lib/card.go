package lib


import (
    "errors"
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "log"
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


func (self *card) ToJSON() string {
    var t tag
    if self.IsLand() {
        t = Tag("land", self.name)
    } else {
        t = Tag("spell", self.name)
    }
    return t.ToJSON()
}


func (self *card) TapsFor() mana {
    return GetCardData(self.name).TapsFor
}


func (self *card) CastingCost() mana {
    return GetCardData(self.name).CastingCost
}


func (self *card) CanBeTitan() bool {
    // Most common fail case is that we can't find Primeval Titan. Let's try to
    // identify those situations sooner.
    return GetCardData(self.name).CanBeTitan
}


func (self *card) AlwaysCast() bool {
    return GetCardData(self.name).AlwaysCast
}


func (self *card) IsLand() bool {
    return GetCardData(self.name).Type == "land"
}


func (self *card) IsBounceLand() bool {
    return self.name == "Simic Growth Chamber"
}


func (self *card) IsCreature() bool {
    return GetCardData(self.name).Type == "creature"
}


func (self *card) IsColorless() bool {
    return GetCardData(self.name).Type == "land" || self.name == "Amulet of Vigor"
}


func (self *card) HasAbility() bool {
    return GetCardData(self.name).ActivationCost.Total != 0
}


func (self *card) ActivationCost() mana {
    return GetCardData(self.name).ActivationCost
}


func (self *card) EntersTapped() bool {
    return GetCardData(self.name).EntersTapped
}


type cardData struct {
    // No need to duplicate card metadata over and over. Cache it by card name
    // and look it up as needed.
    Name string         `yaml:"name"`
    ActivationCost mana `yaml:"activation_cost"`
    CastingCost mana    `yaml:"casting_cost"`
    EntersTapped bool   `yaml:"enters_tapped"`
    Pretty string       `yaml:"pretty"`
    Type string         `yaml:"type"`
    TapsFor mana        `yaml:"taps_for"`
    CanBeTitan bool     `yaml:"can_be_titan"`
    AlwaysCast bool     `yaml:"always_cast"`
}


// Cache card data by name so we don't re-read the file repeatedly
var cardCache = make(map[string]cardData)


func InitCardDataCache() {
    cardDataRaw := []cardData{}
    textBytes, err := ioutil.ReadFile("carddata.yaml")
    if err != nil {
        log.Fatal(err)
    }
    err = yaml.Unmarshal(textBytes, &cardDataRaw)
    if err != nil {
        log.Fatal(err)
    }
    log.Println("loading carddata.yaml")
    for _, cd := range cardDataRaw {
        if cd.Pretty == "" {
            cd.Pretty = slug(cd.Name)
        }
        cardCache[cd.Name] = cd
    }
}


func GetCardData(cardName string) cardData {
    if len(cardCache) == 0 {
        InitCardDataCache()
    }
    return cardCache[cardName]
}


func EnsureCardData(cardNames []string) error {
    if len(cardCache) == 0 {
        InitCardDataCache()
    }
    for _, cardName := range cardNames {
        _, ok := cardCache[cardName]
        if !ok {
            return errors.New("no data for: " + cardName)
        }
    }
    return nil
}
