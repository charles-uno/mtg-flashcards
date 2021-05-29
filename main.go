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
    Verbose     bool        `json:"verbose"`
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
        Verbose: false,
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
    maxTurns := 4

    log.Println(oh)

    game, err := lib.NewGame(
        lib.Shuffled(oh.Library),
        oh.Hand,
        oh.OnThePlay,
        oh.Verbose,
        maxTurns,
    )
    if err != nil {
        reply := map[string]string{"error": err.Error()}
        b, _ := json.Marshal(reply)
        http.Error(w, string(b), http.StatusInternalServerError)
        log.Println("failed to start game at /api/play")
        return
    }
    // Iterate through the turns
    for !game.IsDone() {
        game = game.NextTurn()
    }
    fmt.Fprintf(w, game.ToJSON())
    log.Println("done with calculation at /api/play")
    fmt.Println(game.Pretty())
}


func handleEndToEnd(w http.ResponseWriter, r *http.Request) {
    deck, err := lib.LoadDeck()
    if err != nil {
        reply := map[string]string{"error": err.Error()}
        b, _ := json.Marshal(reply)
        http.Error(w, string(b), http.StatusInternalServerError)
        log.Println("failed to load deck at /api/e2e")
        return
    }
    oh := openingHand{
        Hand: deck[:7],
        Library: deck[7:],
        OnThePlay: flip(),
        Verbose: false,
    }
    maxTurns := 4
    game, err := lib.NewGame(
        lib.Shuffled(oh.Library),
        oh.Hand,
        oh.OnThePlay,
        oh.Verbose,
        maxTurns,
    )
    if err != nil {
        reply := map[string]string{"error": err.Error()}
        b, _ := json.Marshal(reply)
        http.Error(w, string(b), http.StatusInternalServerError)
        log.Println("failed to start game at /api/e2e")
        return
    }
    // Iterate through the turns
    for !game.IsDone() {
        game = game.NextTurn()
    }
    fmt.Fprintf(w, game.ToMiniJSON())
    log.Println("done with calculation at /api/e2e")
    fmt.Println(game.ToMiniJSON())
    fmt.Println(game.Pretty())
}


func main() {
    log.Println("launching service")
    mux := http.NewServeMux()
    mux.HandleFunc("/api/hand", handleOpeningHand)
    mux.HandleFunc("/api/play", handleSequencing)
    mux.HandleFunc("/api/e2e", handleEndToEnd)
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
