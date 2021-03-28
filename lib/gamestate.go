package lib


import (
    "log"
    "strings"
    "strconv"
)


// The gameState is an immutable object which describes a snapshot in time
// during a game. Any change in game state, like drawing a card or casting a
// spell, is enacted by creating a new state.
type gameState struct {
    battlefield cardMap
    done bool
    hand cardMap
    hash string
    landPlays int
    library cardArray
    jsonCache string
    jsonLog string
    manaDebt mana
    manaPool mana
    onThePlay bool
    turn int
}


func (self *gameState) NextStates() []gameState {
    ret := self.passTurn()

    // If we have 6 mana, no need to accumulate more

    // If we have plenty of mana but no Titan, don't bother with Grazer, etc

    for c, _ := range self.hand.Items() {
        if c.IsLand() {
            for _, state := range self.play(c) {
                ret = append(ret, state)
            }
        } else {
            for _, state := range self.cast(c) {
                ret = append(ret, state)
            }
        }
    }
    for c, _ := range self.battlefield.Items() {
        if c.HasAbility() {
            for _, state := range self.activate(c) {
                ret = append(ret, state)
            }
        }
    }
    return ret
}






func (clone gameState) passTurn() []gameState {
    clone.turn += 1
    clone.logBreak()
    clone.logText("turn " + strconv.Itoa(clone.turn))
    // Empty mana pool then tap out
    clone.manaPool = mana{}
    for c, n := range clone.battlefield.Items() {
        m := c.TapsFor()
        clone.manaPool = clone.manaPool.Plus(m.Times(n))
    }
    clone.logManaPool()
    // Pay for Pact
    if clone.manaDebt.Total > 0 {
        m, err := clone.manaPool.Minus(clone.manaDebt)
        if err != nil {
            return []gameState{}
        }
        clone.manaPool = m
        clone.manaDebt = Mana("")
        clone.logText(", pay for pact")
        clone.logManaPool()
    }
    // TODO: pay for Pact
    // Reset land drops. Check for Dryad, Scout, Azusa
    clone.landPlays = 1 +
        clone.battlefield.Count(Card("Dryad of the Ilysian Grove")) +
        clone.battlefield.Count(Card("Sakura-Tribe Scout")) +
        2*clone.battlefield.Count(Card("Azusa, Lost but Seeking"))
    if clone.turn > 1 || !clone.onThePlay {
        return clone.draw(1)
    } else {
        return []gameState{clone}
    }
}


func (clone gameState) activate(c card) []gameState {
    // Is this card on the battlefield?
    if clone.battlefield.Count(c) == 0 {
        return []gameState{}
    }
    // Do we have enough mana to activate it?
    cost := c.ActivationCost()
    m, err := clone.manaPool.Minus(cost)
    if err != nil {
        return []gameState{}
    }
    clone.manaPool = m
    clone.logBreak()
    clone.logText("activate ")
    clone.logCard(c)
    clone.logManaPool()
    // Now figure out what it does
    switch c.name {
        case "Castle Garenbrig":
            return clone.activateCastleGarenbrig()
    }
    log.Fatal("not sure how to activate: " + c.name)
    return []gameState{}
}


func (clone gameState) cast(c card) []gameState {
    // Is this spell in our hand?
    if clone.hand.Count(c) == 0 {
        return []gameState{}
    }
    // Do we have enough mana to cast it?
    cost := c.CastingCost()
    m, err := clone.manaPool.Minus(cost)
    if err != nil {
        return []gameState{}
    }
    clone.manaPool = m
    clone.logBreak()
    clone.logText("cast ")
    clone.logCard(c)
    clone.hand = clone.hand.Minus(c)
    clone.logManaPool()
    // Now figure out what it does
    switch c.name {
        case "Adventurous Impulse":
            return clone.castAdventurousImpulse()
        case "Amulet of Vigor":
            return clone.castAmuletOfVigor()
        case "Ancient Stirrings":
            return clone.castAncientStirrings()
        case "Arboreal Grazer":
            return clone.castArborealGrazer()
        case "Azusa, Lost but Seeking":
            return clone.castAzusaLostButSeeking()
        case "Dryad of the Ilysian Grove":
            return clone.castDryadOfTheIlysianGrove()
        case "Explore":
            return clone.castExplore()
        case "Primeval Titan":
            return clone.castPrimevalTitan()
        case "Summoner's Pact":
            return clone.castSummonersPact()
    }
    log.Fatal("not sure how to cast: " + c.name)
    return []gameState{}
}


func (clone gameState) play(c card) []gameState {
    // Is this land in our hand?
    if clone.hand.Count(c) == 0 {
        return []gameState{}
    }
    // Do we have at least one land play remaining?
    if clone.landPlays <= 0 {
        return []gameState{}
    }
    clone.landPlays -= 1
    clone.logBreak()
    clone.logText("play ")
    clone.logCard(c)
    if c.name == "Castle Garenbrig" {
        if clone.battlefield.Count(Card("Forest")) > 0 {
            return clone.playUntapped(c)
        } else {
            return clone.playTapped(c)
        }
    }
    if c.EntersTapped() {
        return clone.playTapped(c)
    } else {
        return clone.playUntapped(c)
    }
}


func (clone gameState) playTapped(c card) []gameState {
    nAmulets := clone.battlefield.Count(Card("Amulet of Vigor"))
    m := c.TapsFor()
    for i := 0; i < nAmulets; i++ {
        clone.manaPool = clone.manaPool.Plus(m)
        clone.logManaPool()
    }
    return clone.playHelper(c)
}


func (clone gameState) playUntapped(c card) []gameState {
    clone.manaPool = clone.manaPool.Plus(c.TapsFor())
    clone.logManaPool()
    return clone.playHelper(c)
}


func (clone gameState) playHelper(c card) []gameState {
    clone.hand = clone.hand.Minus(c)
    clone.battlefield = clone.battlefield.Plus(c)
    // Watch out for additional effects, if any
    switch c.name {
        case "Bojuka Bog":
            return clone.playBojukaBog()
        case "Castle Garenbrig":
            return clone.playCastleGarenbrig()
        case "Forest":
            return clone.playForest()
        case "Simic Growth Chamber":
            return clone.playSimicGrowthChamber()
    }
    log.Fatal("not sure how to play: " + c.name)
    return []gameState{}
}


func (clone gameState) activateCastleGarenbrig() []gameState {
    clone.manaPool = clone.manaPool.Plus(Mana("GGGGGG"))
    clone.logManaPool()
    // Only activate immediately before casting Titan
    return clone.cast(Card("Primeval Titan"))
}


func (clone gameState) castAmuletOfVigor() []gameState {
    clone.battlefield = clone.battlefield.Plus(Card("Amulet of Vigor"))
    return []gameState{clone}
}


func (self *gameState) castAdventurousImpulse() []gameState {
    ret := []gameState{}
    milled_raw, remaining := self.library.SplitAfter(3)
    milled := CardMap(milled_raw)
    for c, _ := range milled.Items() {
        if c.IsLand() || c.IsCreature() {
            clone := self.clone()
            clone.logText(", mill ")
            clone.logCardMap(milled)
            clone.logText(", grab ")
            clone.logCard(c)
            clone.hand = clone.hand.Plus(c)
            clone.library = remaining
            ret = append(ret, clone)
        }
    }
    return ret
}


func (self *gameState) castArborealGrazer() []gameState {
    ret := []gameState{}
    for c, _ := range self.hand.Items() {
        if !c.IsLand() {
            continue
        }
        clone := self.clone()
        clone.logText(", play ")
        clone.logCard(c)
        ret = append(ret, clone.playTapped(c)...)
    }
    return ret
}


func (self *gameState) castAncientStirrings() []gameState {
    ret := []gameState{}
    milled_raw, remaining := self.library.SplitAfter(5)
    milled := CardMap(milled_raw)
    for c, _ := range milled.Items() {
        if c.IsColorless() {
            clone := self.clone()
            clone.logText(", mill ")
            clone.logCardMap(milled)
            clone.logText(", grab ")
            clone.logCard(c)
            clone.hand = clone.hand.Plus(c)
            clone.library = remaining
            ret = append(ret, clone)
        }
    }
    return ret
}


func (clone gameState) castAzusaLostButSeeking() []gameState {
    clone.battlefield = clone.battlefield.Plus(Card("Azusa, Lost but Seeking"))
    return []gameState{clone}
}


func (clone gameState) castDryadOfTheIlysianGrove() []gameState {
    clone.battlefield = clone.battlefield.Plus(Card("Dryad of the Ilysian Grove"))
    return []gameState{clone}
}


func (clone gameState) castExplore() []gameState {
    clone.landPlays += 1
    return clone.draw(1)
}


func (clone gameState) castPrimevalTitan() []gameState {
    clone.done = true
    return []gameState{clone}
}


func (self *gameState) castSummonersPact() []gameState {
    ret := []gameState{}
    for c, _ := range self.library.Items() {
        if !c.IsCreature() {
            continue
        }
        clone := self.clone()
        clone.hand = clone.hand.Plus(c)
        clone.logText(", grab ")
        clone.logCard(c)
        clone.manaDebt = clone.manaDebt.Plus(Mana("2GG"))
        ret = append(ret, clone)
    }
    return ret
}


func (clone gameState) playBojukaBog() []gameState {
    return []gameState{clone}
}


func (clone gameState) playCastleGarenbrig() []gameState {
    return []gameState{clone}
}


func (clone gameState) playForest() []gameState {
    return []gameState{clone}
}


func (self *gameState) playSimicGrowthChamber() []gameState {
    ret := []gameState{}
    for c, _ := range self.battlefield.Items() {
        if !c.IsLand() {
            continue
        }
        clone := self.clone()
        clone.battlefield = clone.battlefield.Minus(c)
        clone.hand = clone.hand.Plus(c)
        clone.logText(", bounce ")
        clone.logCard(c)
        ret = append(ret, clone)
    }
    return ret
}


func (clone gameState) clone() gameState {
    return clone
}


func (clone gameState) draw(n int) []gameState {
    popped, library := clone.library.SplitAfter(n)
    clone.library = library
    clone.hand = clone.hand.Plus(popped...)
    // ing a card map already throws an extra space in there
    clone.logText(", draw ")
    clone.logCardMap(CardMap(popped))
    return []gameState{clone}
}


func (self *gameState) logManaPool() {
    if self.manaPool.Total > 0 {
        self.logText(", ")
        self.logMana(self.manaPool)
        self.logText(" in pool")
    }
}


func (self *gameState) logBreak() {
    self.resolveCache()
    t := Tag("break", "")
    self.jsonLog += t.ToJSON() + ","
}


func (self *gameState) logText(s string) {
    self.jsonCache += s
}


func (self *gameState) resolveCache() {
    if self.jsonCache != "" {
        t := Tag("text", self.jsonCache)
        self.jsonLog += t.ToJSON() + ","
        self.jsonCache = ""
    }
}


func (self *gameState) logMana(m mana) {
    self.resolveCache()
    self.jsonLog += m.ToJSON() + ","
}


func (self *gameState) logCard(c card) {
    self.resolveCache()
    self.jsonLog += c.ToJSON() + ","
}


func (self *gameState) logCardMap(cm cardMap) {
    for c, n := range cm.Items() {
        self.logCard(c)
        if n > 1 {
            self.logText("*" + strconv.Itoa(n))
        }
        self.logText(" ")
    }
}


func (self *gameState) ToJSON() string {
    // Pull off the last trailing comma so we have a valid JSON list of objects
    return "[" + self.jsonLog[:len(self.jsonLog)-1] + "]"
}


func (self *gameState) Pretty() string {
    return PrettyJSON(self.ToJSON())
}



func (state *gameState) Hash() string {
    // We don't care about order for battlefield or hand, but we do care about
    // the order of the library
    return strings.Join(
        []string{
            state.hand.Pretty(),
            state.battlefield.Pretty(),
            state.manaPool.Pretty(),
            strconv.FormatBool(state.done),
            strconv.Itoa(state.landPlays),
            state.library.Pretty(),
        },
        ";",
    )
}
