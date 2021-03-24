package lib

import (
    "strconv"
)


type mana struct {
    Green int
    Total int
}


func (self *mana) Times(n int) mana {
    return mana{
        Green: self.Green*n,
        Total: self.Total*n,
    }
}


func (self *mana) Plus(other mana) mana {
    return mana{
        Green: self.Green + other.Green,
        Total: self.Total + other.Total,
    }
}


func (m *mana) Pretty() string {
    s := ""
    if m.Total > m.Green || m.Green == 0 {
        s += strconv.Itoa(m.Total - m.Green)
    }
    for i := 0; i < m.Green; i++ {
        s += "G"
    }
    return s
}
