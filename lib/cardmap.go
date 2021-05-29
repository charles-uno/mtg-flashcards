package lib


import (
    "log"
    "sort"
    "strconv"
    "strings"
)


type cardMap struct {
    // A non-ordered container of cards, such as a hand or battlefield
    counts map[card]int
}


func CardMap(cards []card) cardMap {
    counts := make(map[card]int)
    for _, c := range cards {
        counts[c] += 1
    }
    return cardMap{counts: counts}
}


func (self *cardMap) Pretty() string {
    chunks := []string{}
    for c, n := range self.counts {
        if n == 0 {
            continue
        }
        chunk := c.Pretty()
        if n > 1 {
            chunk += "*" + strconv.Itoa(n)
        }
        chunks = append(chunks, chunk)
    }
    sort.Strings(chunks)
    return strings.Join(chunks, " ")
}


func (self *cardMap) Items() map[card]int {
    // TODO: iterate over keys in alphabetical order
    return self.counts
}


func (self *cardMap) Count(c card) int {
    return self.counts[c]
}


func (self *cardMap) Replace(c0 card, c1 card) cardMap {
    counts := make(map[card]int)
    for c, n := range self.counts {
        if c == c0 {
            counts[c1] += n
        } else {
            counts[c] += n
        }
    }
    return cardMap{counts: counts}
}


func (self *cardMap) Plus(cards ...card) cardMap {
    counts := make(map[card]int)
    for c, n := range self.counts {
        counts[c] = n
    }
    for _, c := range cards {
        counts[c] += 1
    }
    return cardMap{counts: counts}
}


func (self *cardMap) Minus(cards ...card) cardMap {
    counts := make(map[card]int)
    for c, n := range self.counts {
        counts[c] = n
    }
    for _, c := range cards {
        if counts[c] > 0 {
            counts[c] -= 1
        } else {
            log.Fatal("can't pop ", c, " from ", counts)
        }
        // Don't leave keys sitting around with zero value
        if counts[c] == 0 {
            delete(counts, c)
        }
    }
    return cardMap{counts: counts}
}
