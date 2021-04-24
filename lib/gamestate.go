package lib


import (
    "encoding/json"
    "log"
    "strings"
    "strconv"
)


// The gameState is an immutable object which describes a snapshot in time
// during a game. Any change in game state, like drawing a card or casting a
// spell, is enacted by creating a new state.
type gameState struct {
    battlefield cardMap
    deadEnd bool
    hand cardMap
    hash string
    landPlays int
    library cardArray
    jsonCache string
    jsonLog string
    manaDebt mana
    manaPool mana
    maxTurns int
    onThePlay bool
    success bool
    timestamp int64
    turn int
}


func NewGameState(library []card, hand []card, otp bool, maxTurns int) gameState {
    state := gameState{
        hand: CardMap(hand),
        // Empty string is fine for the initial game state
        hash: "",
        landPlays: 0,
        library: CardArray(library),
        maxTurns: maxTurns,
        onThePlay: otp,
        timestamp: timestamp(),
        turn: 0,
    }
    if otp {
        state.logText("on the play")
    } else {
        state.logText("on the draw")
    }
    state.logText(", opening hand: ")
    state.logCardMap(state.hand)
    return state
}


func (self *gameState) NextStates() []gameState {
    ret := []gameState{}
    // If we're out of time, see about wrapping up gracefully. Note that
    // timestamp is measured in nanoseconds
    if timestamp() - self.timestamp > 4e9 {
        self.logBreak()
        self.logText("timeout")
        self.MarkDeadEnd()
    }
    // If we've flagged this state as a dead end, just wait out the clock
    if self.deadEnd {
        return self.passTurn()
    }
    // Try to identify doomed lines early rather than playing them out
    ret = append(ret, self.checkForFailure()...)
    if len(ret) > 0 {
        return ret
    }

    // TODO: Try to identify lines of play that are correct 99% of the time.
    // Like if we have a bounce land, turn 1 Amulet is best.

    // TODO: Make sure we allow cantrips to whiff. Otherwise we'll get stuck
    // with the "always play" setting.

    // Don't skip land drops, and don't skip some spells
    if !self.skippedLandDrop() && !self.skippedSpell() {
        ret = append(ret, self.passTurn()...)
    }
    for c, _ := range self.hand.Items() {
        if c.IsLand() {
            // No need to ever go above 6 mana
            if self.manaPool.Total >= 6 {
                continue
            }
            ret = append(ret, self.play(c)...)
        } else {
            // TODO: If we have 6 mana, only cast Titan or Pact?
            ret = append(ret, self.cast(c)...)
        }
    }
    for c, _ := range self.battlefield.Items() {
        if c.HasAbility() {
            ret = append(ret, self.activate(c)...)
        }
    }
    return ret
}


func (self *gameState) checkForFailure() []gameState {
    if self.turn < self.maxTurns {
        return []gameState{}
    }
    // If we don't have Primeval Titan or a way to find it, bail
    noTitan := true
    for c, _ := range self.hand.Items() {
        if c.CanBeTitan() {
            noTitan = false
        }
    }
    if noTitan {
        clone := self.clone()
        clone.logBreak()
        clone.logText("failed to find ")
        clone.logCard(Card("Primeval Titan"))
        clone.MarkDeadEnd()
        return clone.passTurn()
    }
    return []gameState{}
}


func (self *gameState) skippedLandDrop() bool {
    // Return true if all of the following are true:
    // 1. We have a non-bounce land in hand
    // 2. We have a land play remaining
    // 3. We have less than 6 mana available
    // Note: this has a very small chance to miss lines! For example, if we
    // have 4 mana it might make sense to keep Bojuka Bog in hand in case we
    // draw into a second Amulet.
    if self.landPlays > 0 && self.manaPool.Total < 6 {
        for c, _ := range self.hand.Items() {
            if c.IsLand() && !c.IsBounceLand() {
                return true
            }
        }
    }
    return false
}


func (self *gameState) skippedSpell() bool {
    // The computer will sometimes find it advantageou to cast a cantrip later
    // because it "knows" the order of the deck. Do what we can to suppress
    // that non-human play pattern. We might want to hold onto Explore for ramp
    // purposes with multiple copies of Amulet, but a human player is never
    // going to pass the turn rather than cast Ancient Stirrings.
    for c, _ := range self.hand.Items() {
        if c.AlwaysCast() && self.manaPool.CanPay(c.CastingCost()) {
            return true
        }
    }
    return false
}


func (clone gameState) clone() gameState {
    return clone
}


func (clone gameState) passTurn() []gameState {
    clone.turn += 1
    if clone.turn > clone.maxTurns {
        // Nice to have output here in terms of traceability when debugging,
        // but it doesn't read nicely.
//        clone.logBreak()
//        clone.logText("no more turns")
        clone.MarkDeadEnd()
        return []gameState{clone}
    }
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
    if c.name != "Primeval Titan" && c.name != "Summoner's Pact" {
        clone.logManaPool()
    }
    // Now figure out what it does
    switch c.name {
        case "Abundant Harvest":
            return clone.castAbundantHarvest()
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
        case "Elvish Spirit Guide":
            return clone.castElvishSpiritGuide()
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
        case "Wastes":
            return clone.playWastes()
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
    self.logText(", mill ")
    self.logCardMap(milled)
    self.library = remaining
    for c, _ := range milled.Items() {
        if c.IsLand() || c.IsCreature() {
            clone := self.clone()
            clone.logText(", grab ")
            clone.logCard(c)
            clone.hand = clone.hand.Plus(c)
            ret = append(ret, clone)
        }
    }
    if len(ret) == 0 {
        clone := self.clone()
        clone.logText(", whiff")
        ret = append(ret, clone)
    }
    return ret
}

func (self *gameState) castAbundantHarvest() []gameState {
    ret := []gameState{}
    for _, chooseLand := range []bool{true, false} {
        clone := self.clone()
        if chooseLand {
            clone.logText(", choose land")
        } else {
            clone.logText(", choose nonland")
        }
        i := 0
        for {
            nextCard := self.library.Get(i)
            if (chooseLand && nextCard.IsLand()) || (!chooseLand && !nextCard.IsLand()) {
                break
            }
            i += 1
        }
        revealed, library := self.library.SplitAfter(i+1)
        keep := revealed[i]
        clone.library = library
        clone.hand = clone.hand.Plus(keep)
        clone.logText(", reveal")
        for _, c := range revealed {
            clone.logText(" ")
            clone.logCard(c)
        }
        clone.logText(", grab ")
        clone.logCard(keep)
        ret = append(ret, clone)
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
    self.library = remaining
    self.logText(", mill ")
    self.logCardMap(milled)
    for c, _ := range milled.Items() {
        if c.IsColorless() {
            clone := self.clone()
            clone.logText(", grab ")
            clone.logCard(c)
            clone.hand = clone.hand.Plus(c)
            ret = append(ret, clone)
        }
    }
    if len(ret) == 0 {
        clone := self.clone()
        clone.logText(", whiff")
        ret = append(ret, clone)
    }
    return ret
}


func (clone gameState) castAzusaLostButSeeking() []gameState {
    clone.battlefield = clone.battlefield.Plus(Card("Azusa, Lost but Seeking"))
    return []gameState{clone}
}


func (clone gameState) castDryadOfTheIlysianGrove() []gameState {
    clone.battlefield = clone.battlefield.Plus(Card("Dryad of the Ilysian Grove"))
    clone.landPlays += 1
    return []gameState{clone}
}


func (clone gameState) castElvishSpiritGuide() []gameState {
    if clone.turn <= 0 {
        return []gameState{}
    }
    clone.manaPool = clone.manaPool.Plus(Mana("G"))
    return []gameState{clone}
}


func (clone gameState) castExplore() []gameState {
    clone.landPlays += 1
    return clone.draw(1)
}


func (clone gameState) castPrimevalTitan() []gameState {
    clone.success = true
    return []gameState{clone}
}


func (self *gameState) castSummonersPact() []gameState {
    ret := []gameState{}
    for c, _ := range self.library.Items() {
        if !c.IsCreature() {
            continue
        }
        // Don't Pact for a card we already have in hand
        if self.hand.Count(c) > 0 {
            continue
        }
        clone := self.clone()
        clone.hand = clone.hand.Plus(c)
        clone.logText(", grab ")
        clone.logCard(c)
        clone.manaDebt = clone.manaDebt.Plus(Mana("2GG"))
        // To cut down on near-identical sequences, cast immediately
        for _, state := range clone.cast(c) {
            ret = append(ret, state)
        }
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


func (clone gameState) playWastes() []gameState {
    return []gameState{clone}
}


func (self *gameState) playSimicGrowthChamber() []gameState {
    ret := []gameState{}
    nAmulets := self.battlefield.Count(Card("Amulet of Vigor"))
    for c, _ := range self.battlefield.Items() {
        if !c.IsLand() {
            continue
        }
        // Without an Amulet, never pick itself back up
        if nAmulets == 0 && c.IsBounceLand() {
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
    // Might want this back later, but for now let's skip it
    return

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
        if n > 1 {
            // TODO: Use the unicode multiplication symbol instead
            self.logText(strconv.Itoa(n) + "*")
        }
        self.logCard(c)
        self.logText(" ")
    }
    // Get rid of the trailing whitespace
    self.jsonCache = ""
}


func (self *gameState) MarkDeadEnd() {
    if !self.deadEnd {
        self.logBreak()
        self.logText("giving up!")
        self.deadEnd = true
    }
}


func (self *gameState) ToJSON() string {
    self.resolveCache()
    turn := "-1"
    if self.success {
        turn = strconv.Itoa(self.turn)
    }
    // Pull off the last trailing comma so we have a valid JSON list of objects
    return "{\"turn\": " + turn + ", " +
        "\"plays\": [" + self.jsonLog[:len(self.jsonLog)-1] + "]}\n"
}


type miniGame struct {
    Turn        int     `json:"turn"`
    OnThePlay   bool    `json:"onThePlay"`
}


func (self *gameState) ToMiniJSON() string {
    mini := miniGame{
        Turn: self.turn,
        OnThePlay: self.onThePlay,
    }
    b, _ := json.Marshal(mini)
    return string(b) + "\n"
}


func (self *gameState) LogSize() int {
    return len(self.jsonLog)
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
            strconv.FormatBool(state.success),
            strconv.FormatBool(state.deadEnd),
            strconv.Itoa(state.landPlays),
            state.library.Pretty(),
        },
        ";",
    )
}
