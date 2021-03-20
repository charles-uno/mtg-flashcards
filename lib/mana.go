package lib

import (
    "strconv"
)


type mana struct {
    Green int
    Total int
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
