package lib


import (
    "strings"
)


type cardArray struct {
    // An ordered sequence of card names, such as a library
    arr []card
}


func CardArray(arr []card) cardArray {
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
