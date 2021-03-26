package lib

import (
    "errors"
    "log"
    "strconv"
)


type mana struct {
    Green int   `yaml:"green"`
    Total int   `yaml:"total"`
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


func (self *mana) Minus(other mana) (mana, error) {
    if self.Green >= other.Green && self.Total >= other.Total {
        total := self.Total - other.Total
        green := self.Green - other.Green
        if green > total {
            green = total
        }
        return mana{Green: green, Total: total}, nil
    } else {
        text := "can't subtract " + self.Pretty() + " - " + other.Pretty()
        return mana{}, errors.New(text)
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


func Mana(s string) mana {
    green := 0
    total := 0
    for _, c := range s {
        if c == 'G' {
            green += 1
            total += 1
        } else if '0' <= c && c <= '9' {
            total += int(c - '0')
        } else {
            log.Fatal("failed to parse mana cost: " + s)
        }
    }
    return mana{Green: green, Total: total}
}
