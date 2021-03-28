package lib


import (
    "encoding/json"
    "log"
    "strings"
)


type tag struct {
    Type string `json:"type"`
    Text string `json:"text"`
}


func Tag(tagType string, tagText string) tag {
    return tag{Type: tagType, Text: tagText}
}


func (self *tag) ToJSON() string {
    b, err := json.Marshal(self)
    if err != nil {
        log.Fatal("failed to marshal:", self)
    }
    return string(b)
//    return "{\"type\": \"" + self.Type + "\", \"text\": \"" + self.Text + "\"}"
}


func PrettyJSON(s string) string {
    tags := []tag{}
    err := json.Unmarshal([]byte(s), &tags)
    if err != nil {
        log.Fatal("failed to unmarshal:", s)
    }
    ret := ""
    for _, t := range tags {
        if t.Type == "text" {
            ret += t.Text
        } else if t.Type == "break" {
            ret += "\n"
        } else if t.Type == "mana" {
            ret += "\u001b[35m" + t.Text + "\u001b[0m"
        } else if t.Type == "land" {
            ret += "\u001b[33m" + slug(t.Text) + "\u001b[0m"
        } else if t.Type == "spell" {
            ret += "\u001b[32m" + slug(t.Text) + "\u001b[0m"
        } else {
            log.Fatal("not sure how to export type", t.Type)
        }
    }
    return ret
}


func slug(s string) string {
    for _, c := range []string{" ", "-", "'", ","} {
        s = strings.ReplaceAll(s, c, "")
    }
    return s
}
