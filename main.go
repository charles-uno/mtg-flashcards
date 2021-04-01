package main

import (
    "encoding/json"
    "fmt"
    "log"
    "math/rand"
    "net/http"
    "time"

    "github.com/charles-uno/mtgserver/lib"
    "github.com/rs/cors"
)


type openingHand struct {
    Hand        []string    `json:"hand"`
    Library     []string    `json:"library"`
    OnThePlay   bool        `json:"onThePlay"`
}


func handleOpeningHand(w http.ResponseWriter, r *http.Request) {
    deck, err := lib.LoadDeck()
    if err != nil {
        reply := map[string]string{"error": err.Error()}
        b, _ := json.Marshal(reply)
        http.Error(w, string(b), http.StatusInternalServerError)
        log.Println("failed to load deck at /api/hand")
        return
    }
    oh := openingHand{
        Hand: deck[:7],
        Library: deck[7:],
        OnThePlay: flip(),
    }
    log.Println("endpoint hit: /api/hand")
    json.NewEncoder(w).Encode(oh)
}


func handleSequencing(w http.ResponseWriter, r *http.Request) {
    oh := openingHand{}
    err := json.NewDecoder(r.Body).Decode(&oh)
    if err != nil {
        reply := map[string]string{"error": err.Error()}
        b, _ := json.Marshal(reply)
        http.Error(w, string(b), http.StatusBadRequest)
        log.Println("bad payload at /api/play")
        return
    }
    game, err := lib.NewGame(oh.Hand, lib.Shuffled(oh.Library), oh.OnThePlay)
    if err != nil {
        reply := map[string]string{"error": err.Error()}
        b, _ := json.Marshal(reply)
        http.Error(w, string(b), http.StatusInternalServerError)
        log.Println("failed to start game at /api/play")
        return
    }
    // Iterate through the turns
    for game.IsNotDone() {
        log.Println("starting turn", game.Turn, "with", game.Size(), "states")
        game = game.NextTurn()
    }
    fmt.Fprintf(w, game.ToJSON())
    log.Println("done with calculation at /api/play")
    fmt.Println(game.Pretty())
}


func main() {
    log.Println("launching service")
    mux := http.NewServeMux()
    mux.HandleFunc("/api/hand", handleOpeningHand)
    mux.HandleFunc("/api/play", handleSequencing)
    // Default CORS handler allows GET and POST from anywhere. To go back to
    // default settings, lose the handler and use nil instead
    handler := cors.Default().Handler(mux)
    log.Fatal(http.ListenAndServe(":5001", handler))
}


func flip() bool {
    // Random generator should be seeded from shuffling, but let's be sure. We
    // only call this once per game anyway.
    rand.Seed(time.Now().UTC().UnixNano())
    return rand.Intn(2) == 0
}
